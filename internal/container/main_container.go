package container

import (
	"github.com/InstaySystem/is-be/internal/config"
	"github.com/InstaySystem/is-be/internal/initialization"
	"github.com/InstaySystem/is-be/internal/middleware"
	"github.com/InstaySystem/is-be/internal/provider/jwt"
	"github.com/InstaySystem/is-be/pkg/bcrypt"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"github.com/sony/sonyflake/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	FileCtn *FileContainer
	AuthCtn *AuthContainer
	UserCtn *UserContainer
	AuthMid *middleware.AuthMiddleware
}

func NewContainer(
	cfg *config.Config,
	db *gorm.DB,
	s3 *initialization.S3,
	sf *sonyflake.Sonyflake,
	logger *zap.Logger,
) *Container {
	sfGen := snowflake.NewGenerator(sf)
	bHash := bcrypt.NewHasher(10)
	jwtProvider := jwt.NewJWTProvider(cfg.JWT.SecretKey)

	fileCtn := NewFileContainer(cfg, s3, logger)
	authCtn := NewAuthContainer(cfg, db, logger, bHash, jwtProvider)
	userCtn := NewUserContainer(db, sfGen, logger, bHash)

	authMid := middleware.NewAuthMiddleware(cfg.JWT.AccessName, cfg.JWT.RefreshName, userCtn.Repo, jwtProvider, logger)

	return &Container{
		fileCtn,
		authCtn,
		userCtn,
		authMid,
	}
}
