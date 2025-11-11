package container

import (
	"time"

	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/provider/cache"
	"github.com/InstaySystem/is-be/internal/repository"
	repoImpl "github.com/InstaySystem/is-be/internal/repository/implement"
	svcImpl "github.com/InstaySystem/is-be/internal/service/implement"
	"github.com/InstaySystem/is-be/pkg/bcrypt"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserContainer struct {
	Hdl  *handler.UserHandler
	Repo repository.UserRepository
}

func NewUserContainer(
	db *gorm.DB,
	sfGen snowflake.Generator,
	logger *zap.Logger,
	bHash bcrypt.Hasher,
	refreshExpiresIn time.Duration,
	cacheProvider cache.CacheProvider,
) *UserContainer {
	userRepo := repoImpl.NewUserRepository(db)
	departmentRepo := repoImpl.NewDepartmentRepository(db)
	svc := svcImpl.NewUserService(userRepo, departmentRepo, sfGen, logger, bHash, refreshExpiresIn, cacheProvider)
	hdl := handler.NewUserHandler(svc)

	return &UserContainer{
		hdl,
		userRepo,
	}
}
