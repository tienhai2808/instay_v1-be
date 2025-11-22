package container

import (
	"time"

	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/provider/cache"
	"github.com/InstaySystem/is-be/internal/repository"
	svcImpl "github.com/InstaySystem/is-be/internal/service/implement"
	"github.com/InstaySystem/is-be/pkg/bcrypt"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
)

type UserContainer struct {
	Hdl  *handler.UserHandler
}

func NewUserContainer(
	userRepo repository.UserRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
	bHash bcrypt.Hasher,
	refreshExpiresIn time.Duration,
	cacheProvider cache.CacheProvider,
) *UserContainer {
	svc := svcImpl.NewUserService(userRepo, sfGen, logger, bHash, refreshExpiresIn, cacheProvider)
	hdl := handler.NewUserHandler(svc)

	return &UserContainer{hdl}
}
