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
	"github.com/InstaySystem/is-be/internal/initialization"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	cfg  *config.Config
	http *http.Server
	db   *initialization.DB
	rdb  *redis.Client
	mq   *initialization.MQ
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

	_ = r.Group(cfg.Server.APIPrefix)

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
	}, nil
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) {
	if s.db != nil {
		s.db.Close()
	}

	if s.rdb != nil {
		s.rdb.Close()
	}

	if s.mq != nil {
		s.mq.Close()
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