package container

import (
	"github.com/InstaySystem/is-be/internal/handler"
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
) *UserContainer {
	repo := repoImpl.NewUserRepository(db)
	svc := svcImpl.NewUserService(repo, sfGen, logger, bHash)
	hdl := handler.NewUserHandler(svc)

	return &UserContainer{
		hdl,
		repo,
	}
}
