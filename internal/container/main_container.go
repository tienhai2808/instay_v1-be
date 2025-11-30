package container

import (
	"github.com/InstaySystem/is-be/internal/config"
	"github.com/InstaySystem/is-be/internal/hub"
	"github.com/InstaySystem/is-be/internal/initialization"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/InstaySystem/is-be/internal/provider/cache"
	"github.com/InstaySystem/is-be/internal/provider/jwt"
	"github.com/InstaySystem/is-be/internal/provider/mq"
	"github.com/InstaySystem/is-be/internal/provider/smtp"
	"github.com/InstaySystem/is-be/internal/repository"
	repoImpl "github.com/InstaySystem/is-be/internal/repository/implement"
	"github.com/InstaySystem/is-be/pkg/bcrypt"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/sony/sonyflake/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	FileCtn         *FileContainer
	AuthCtn         *AuthContainer
	UserCtn         *UserContainer
	DepartmentCtn   *DepartmentContainer
	ServiceCtn      *ServiceContainer
	RequestCtn      *RequestContainer
	RoomCtn         *RoomContainer
	BookingCtn      *BookingContainer
	OrderCtn        *OrderContainer
	NotificationCtn *NotificationContainer
	ChatCtn         *ChatContainer
	ReviewCtn       *ReviewContainer
	SSECtn          *SSEContainer
	WSCtn           *WSContainer
	AuthMid         *middleware.AuthMiddleware
	ReqMid          *middleware.RequestMiddleware
	SMTPProvider    smtp.SMTPProvider
	MQProvider      mq.MessageQueueProvider
	SfGen           snowflake.Generator
	BookingRepo     repository.BookingRepository
	SSEHub          *hub.SSEHub
	WSHub           *hub.WSHub
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
	sseHub := hub.NewSSEHub()

	userRepo := repoImpl.NewUserRepository(db)
	serviceRepo := repoImpl.NewServiceRepository(db)
	roomRepo := repoImpl.NewRoomRepository(db)
	requestRepo := repoImpl.NewRequestRepository(db)
	departmentRepo := repoImpl.NewDepartmentRepository(db)
	bookingRepo := repoImpl.NewBookingRepository(db)
	orderRepo := repoImpl.NewOrderRepository(db)
	notificationRepo := repoImpl.NewNotificationRepository(db)
	chatRepo := repoImpl.NewChatRepository(db)
	reviewRepo := repoImpl.NewReviewRepository(db)

	fileCtn := NewFileContainer(cfg, s3, logger)
	authCtn := NewAuthContainer(cfg, db, logger, bHash, jwtProvider, cacheProvider, mqProvider)
	userCtn := NewUserContainer(userRepo, sfGen, logger, bHash, cfg.JWT.RefreshExpiresIn, cacheProvider)
	departmentCtn := NewDepartmentContainer(departmentRepo, sfGen, logger)
	serviceCtn := NewServiceContainer(db, serviceRepo, sfGen, logger, mqProvider)
	requestCtn := NewRequestContainer(db, requestRepo, orderRepo, notificationRepo, sfGen, logger, mqProvider)
	roomCtn := NewRoomContainer(roomRepo, sfGen, logger)
	bookingCtn := NewBookingContainer(bookingRepo, logger)
	orderCtn := NewOrderContainer(db, orderRepo, bookingRepo, serviceRepo, notificationRepo, sfGen, logger, cacheProvider, jwtProvider, mqProvider, cfg.JWT.GuestName)
	notificationCtn := NewNotificationContainer(db, notificationRepo, logger, sfGen)
	chatCtn := NewChatContainer(db, chatRepo, orderRepo, sfGen, logger)
	reviewCtn := NewReviewContainer(reviewRepo, sfGen, logger)
	wsHub := hub.NewWSHub(chatCtn.Svc)
	sseCtn := NewSSEContainer(sseHub)
	wsCtn := NewWSContainer(wsHub)

	authMid := middleware.NewAuthMiddleware(cfg.JWT.AccessName, cfg.JWT.RefreshName, cfg.JWT.GuestName, userRepo, jwtProvider, logger, cacheProvider)
	reqMid := middleware.NewRequestMiddleware(logger)

	return &Container{
		fileCtn,
		authCtn,
		userCtn,
		departmentCtn,
		serviceCtn,
		requestCtn,
		roomCtn,
		bookingCtn,
		orderCtn,
		notificationCtn,
		chatCtn,
		reviewCtn,
		sseCtn,
		wsCtn,
		authMid,
		reqMid,
		smtpProvider,
		mqProvider,
		sfGen,
		bookingRepo,
		sseHub,
		wsHub,
	}
}
