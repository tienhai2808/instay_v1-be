package container

import (
	"github.com/InstaySystem/is-be/internal/config"
	"github.com/InstaySystem/is-be/internal/initialization"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/InstaySystem/is-be/internal/provider/cache"
	"github.com/InstaySystem/is-be/internal/provider/jwt"
	"github.com/InstaySystem/is-be/internal/provider/mq"
	"github.com/InstaySystem/is-be/internal/provider/smtp"
	"github.com/InstaySystem/is-be/pkg/bcrypt"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/sony/sonyflake/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	FileCtn      *FileContainer
	AuthCtn      *AuthContainer
	UserCtn      *UserContainer
	AuthMid      *middleware.AuthMiddleware
	SMTPProvider smtp.SMTPProvider
	MQProvider   mq.MessageQueueProvider
}

func NewContainer(
	cfg *config.Config,
	db *gorm.DB,
	rdb *redis.Client,
	s3 *initialization.S3,
	sf *sonyflake.Sonyflake,
	logger *zap.Logger,
	mqConn *amqp091.Connection,
	mqChan *amqp091.Channel,
) *Container {
	sfGen := snowflake.NewGenerator(sf)
	bHash := bcrypt.NewHasher(10)
	jwtProvider := jwt.NewJWTProvider(cfg.JWT.SecretKey)
	smtpProvider := smtp.NewSMTPProvider(cfg)
	mqProvider := mq.NewMessageQueueProvider(mqConn, mqChan, logger)
	cacheProvider := cache.NewCacheProvider(rdb)

	fileCtn := NewFileContainer(cfg, s3, logger)
	authCtn := NewAuthContainer(cfg, db, logger, bHash, jwtProvider, cacheProvider, mqProvider)
	userCtn := NewUserContainer(db, sfGen, logger, bHash)

	authMid := middleware.NewAuthMiddleware(cfg.JWT.AccessName, cfg.JWT.RefreshName, userCtn.Repo, jwtProvider, logger)

	return &Container{
		fileCtn,
		authCtn,
		userCtn,
		authMid,
		smtpProvider,
		mqProvider,
	}
}
