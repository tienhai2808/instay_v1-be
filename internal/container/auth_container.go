package container

import (
	"github.com/InstaySystem/is-be/internal/config"
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/provider/cache"
	"github.com/InstaySystem/is-be/internal/provider/jwt"
	"github.com/InstaySystem/is-be/internal/provider/mq"
	repoImpl "github.com/InstaySystem/is-be/internal/repository/implement"
	svcImpl "github.com/InstaySystem/is-be/internal/service/implement"
	"github.com/InstaySystem/is-be/pkg/bcrypt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthContainer struct {
	Hdl *handler.AuthHandler
}

func NewAuthContainer(
	cfg *config.Config,
	db *gorm.DB,
	logger *zap.Logger,
	bHash bcrypt.Hasher,
	jwtProvider jwt.JWTProvider,
	cacheProvider cache.CacheProvider,
	mqProvider mq.MessageQueueProvider,
) *AuthContainer {
	userRepo := repoImpl.NewUserRepository(db)
	svc := svcImpl.NewAuthService(userRepo, logger, bHash, jwtProvider, cfg, cacheProvider, mqProvider)
	hdl := handler.NewAuthHandler(svc, cfg)

	return &AuthContainer{hdl}
}
