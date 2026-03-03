package container

import (
	"cloud.google.com/go/storage"
	"github.com/InstaySystem/is_v1-be/internal/config"
	"github.com/InstaySystem/is_v1-be/internal/hub"
	"github.com/InstaySystem/is_v1-be/internal/middleware"
	"github.com/InstaySystem/is_v1-be/internal/provider/cache"
	"github.com/InstaySystem/is_v1-be/internal/provider/jwt"
	"github.com/InstaySystem/is_v1-be/internal/provider/mq"
	"github.com/InstaySystem/is_v1-be/internal/provider/smtp"
	"github.com/InstaySystem/is_v1-be/internal/repository"
	repoImpl "github.com/InstaySystem/is_v1-be/internal/repository/implement"
	"github.com/InstaySystem/is_v1-be/pkg/bcrypt"
	"github.com/InstaySystem/is_v1-be/pkg/snowflake"
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
	DashboardCtn    *DashboardContainer
	SSECtn          *SSEContainer
	WSCtn           *WSContainer
	AuthMid         *middleware.AuthMiddleware
	ReqMid          *middleware.RequestMiddleware
	SMTPProvider    smtp.SMTPProvider
	MQProvider      mq.MessageQueueProvider
	SfGen           snowflake.Generator
	BHash           bcrypt.Hasher
	BookingRepo     repository.BookingRepository
	UserRepo        repository.UserRepository
	SSEHub          *hub.SSEHub
	WSHub           *hub.WSHub
}

func NewContainer(
	cfg *config.Config,
	db *gorm.DB,
	rdb *redis.Client,
	gcs *storage.Client,
	sf *sonyflake.Sonyflake,
	logger *zap.Logger,
	rmq *amqp091.Connection,
) *Container {
	sfGen := snowflake.NewGenerator(sf)
	bHash := bcrypt.NewHasher(10)
	jwtProvider := jwt.NewJWTProvider(cfg.JWT.SecretKey)
	smtpProvider := smtp.NewSMTPProvider(cfg)
	mqProvider := mq.NewMessageQueueProvider(rmq, logger)
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

	fileCtn := NewFileContainer(cfg, gcs, logger)
	authCtn := NewAuthContainer(cfg, db, userRepo, logger, bHash, jwtProvider, cacheProvider, mqProvider)
	userCtn := NewUserContainer(userRepo, sfGen, logger, bHash, cfg.JWT.RefreshExpiresIn, cacheProvider)
	departmentCtn := NewDepartmentContainer(departmentRepo, sfGen, logger)
	serviceCtn := NewServiceContainer(db, serviceRepo, sfGen, logger, mqProvider)
	requestCtn := NewRequestContainer(db, requestRepo, orderRepo, notificationRepo, sfGen, logger, mqProvider)
	roomCtn := NewRoomContainer(roomRepo, sfGen, logger)
	bookingCtn := NewBookingContainer(bookingRepo, logger)
	orderCtn := NewOrderContainer(db, orderRepo, bookingRepo, roomRepo, serviceRepo, notificationRepo, chatRepo, sfGen, logger, cacheProvider, jwtProvider, mqProvider, cfg.JWT.GuestName)
	notificationCtn := NewNotificationContainer(db, notificationRepo, logger, sfGen)
	chatCtn := NewChatContainer(db, chatRepo, orderRepo, userRepo, sfGen, logger)
	reviewCtn := NewReviewContainer(reviewRepo, sfGen, logger)
	dashboardCtn := NewDashboardContainer(userRepo, roomRepo, serviceRepo, bookingRepo, orderRepo, requestRepo, reviewRepo, logger)
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
		dashboardCtn,
		sseCtn,
		wsCtn,
		authMid,
		reqMid,
		smtpProvider,
		mqProvider,
		sfGen,
		bHash,
		bookingRepo,
		userRepo,
		sseHub,
		wsHub,
	}
}
