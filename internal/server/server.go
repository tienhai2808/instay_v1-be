package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/InstaySystem/is-be/internal/config"
	"github.com/InstaySystem/is-be/internal/container"
	"github.com/InstaySystem/is-be/internal/initialization"
	"github.com/InstaySystem/is-be/internal/router"
	"github.com/InstaySystem/is-be/internal/worker"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type Server struct {
	cfg          *config.Config
	http         *http.Server
	db           *initialization.DB
	rdb          *redis.Client
	mq           *initialization.MQ
	listenWorker *worker.ListenWorker
	logger       *zap.Logger
}

func NewServer(cfg *config.Config) (*Server, error) {
	db, err := initialization.InitPostgreSQL(cfg)
	if err != nil {
		return nil, err
	}

	rdb, err := initialization.InitRedis(cfg)
	if err != nil {
		return nil, err
	}

	mq, err := initialization.InitRabbitMQ(cfg)
	if err != nil {
		return nil, err
	}

	s3, err := initialization.InitS3(cfg)
	if err != nil {
		return nil, err
	}

	sf, err := initialization.InitSnowFlake()
	if err != nil {
		return nil, err
	}

	logger, err := initialization.InitLogger()
	if err != nil {
		return nil, err
	}

	ctn := container.NewContainer(cfg, db.Gorm, rdb, s3, sf, logger, mq.Conn, mq.Chan)

	mqWorker := worker.NewMQWorker(cfg, ctn.MQProvider, ctn.SMTPProvider, s3.Client, logger)
	mqWorker.Start()

	listenWorker := worker.NewListenWorker(cfg, ctn.BookingCtn.Repo, ctn.SfGen, logger)
	listenWorker.Start()

	r := gin.Default()
	if err = r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		return nil, fmt.Errorf("setup Proxy failed: %w", err)
	}

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(corsConfig))
	r.Use(ctn.ReqMid.Recovery())

	api := r.Group(cfg.Server.APIPrefix)

	router.FileRouter(api, ctn.FileCtn.Hdl)
	router.UserRouter(api, ctn.UserCtn.Hdl, ctn.AuthMid)
	router.AuthRouter(api, ctn.AuthCtn.Hdl, ctn.AuthMid)
	router.DepartmentRouter(api, ctn.DepartmentCtn.Hdl, ctn.AuthMid)
	router.ServiceRouter(api, ctn.ServiceCtn.Hdl, ctn.AuthMid)
	router.RequestRouter(api, ctn.RequestCtn.Hdl, ctn.AuthMid)
	router.RoomRouter(api, ctn.RoomCtn.Hdl, ctn.AuthMid)
	router.BookingRouter(api, ctn.BookingCtn.Hdl, ctn.AuthMid)

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	http := &http.Server{
		Addr:           addr,
		Handler:        r,
		MaxHeaderBytes: 5 * 1024 * 1024,
	}

	return &Server{
		cfg,
		http,
		db,
		rdb,
		mq,
		listenWorker,
		logger,
	}, nil
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) {
	if s.listenWorker != nil {
		s.listenWorker.Stop()
	}
	
	if s.db != nil {
		s.db.Close()
	}

	if s.rdb != nil {
		s.rdb.Close()
	}

	if s.mq != nil {
		s.mq.Close()
	}

	if s.logger != nil {
		s.logger.Sync()
	}

	if s.http != nil {
		if err := s.http.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown failed: %v", err)
			return
		}
	}

	log.Println("Server stopped successfully")
}

func (s *Server) GracefulShutdown(ch <-chan error) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-ch:
		log.Printf("Server run failed: %v", err)
	case <-ctx.Done():
		log.Println("Server stop signal")
	}

	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.Shutdown(shutdownCtx)
}
