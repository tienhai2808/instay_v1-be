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

	mqWorker := worker.NewMQWorker(cfg, ctn.MQProvider, ctn.SMTPProvider, s3.Client, logger, ctn.SSEHub)
	mqWorker.Start()

	listenWorker := worker.NewListenWorker(cfg, ctn.BookingRepo, ctn.SfGen, logger)
	listenWorker.Start()

	go ctn.SSEHub.Run()
	go ctn.WSHub.Run()

	r := gin.Default()
	_ = r.SetTrustedProxies(nil)

	corsConfig := cors.Config{
		AllowOrigins:     cfg.Server.AllowOrigins,
		AllowMethods:     cfg.Server.AllowMethods,
		AllowHeaders:     cfg.Server.AllowHeaders,
		ExposeHeaders:    cfg.Server.ExposeHeaders,
		AllowCredentials: cfg.Server.AllowCredentials,
		MaxAge:           cfg.Server.MaxAge,
	}

	r.Use(cors.New(corsConfig))
	r.Use(ctn.ReqMid.Recovery())
	r.Use(ctn.ReqMid.ErrorHandler())

	api := r.Group(cfg.Server.APIPrefix)

	router.FileRouter(api, ctn.FileCtn.Hdl)
	router.UserRouter(api, ctn.UserCtn.Hdl, ctn.AuthMid)
	router.AuthRouter(api, ctn.AuthCtn.Hdl, ctn.AuthMid)
	router.DepartmentRouter(api, ctn.DepartmentCtn.Hdl, ctn.AuthMid)
	router.ServiceRouter(api, ctn.ServiceCtn.Hdl, ctn.AuthMid)
	router.RequestRouter(api, ctn.RequestCtn.Hdl, ctn.AuthMid)
	router.RoomRouter(api, ctn.RoomCtn.Hdl, ctn.AuthMid)
	router.BookingRouter(api, ctn.BookingCtn.Hdl, ctn.AuthMid)
	router.OrderRouter(api, ctn.OrderCtn.Hdl, ctn.AuthMid)
	router.NotificationRouter(api, ctn.NotificationCtn.Hdl, ctn.AuthMid)
	router.ChatRouter(api, ctn.ChatCtn.Hdl, ctn.AuthMid)
	router.ReviewRouter(api, ctn.ReviewCtn.Hdl, ctn.AuthMid)
	router.DashboardRouter(api, ctn.DashboardCtn.Hdl, ctn.AuthMid)
	router.SSERouter(api, ctn.SSECtn.Hdl, ctn.AuthMid)
	router.WSRouter(api, ctn.WSCtn.Hdl, ctn.AuthMid)

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	http := &http.Server{
		Addr:           addr,
		Handler:        r,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes * 1024 * 1024,
		IdleTimeout:    cfg.Server.IdleTimeout,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
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
