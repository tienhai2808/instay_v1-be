package container

import (
	"github.com/InstaySystem/is-be/internal/handler"
	"github.com/InstaySystem/is-be/internal/repository"
	svcImpl "github.com/InstaySystem/is-be/internal/service/implement"
	"github.com/InstaySystem/is-be/pkg/snowflake"
	"go.uber.org/zap"
)

type DepartmentContainer struct {
	Hdl *handler.DepartmentHandler
}

func NewDepartmentContainer(
	departmentRepo repository.DepartmentRepository,
	sfGen snowflake.Generator,
	logger *zap.Logger,
) *DepartmentContainer {
	svc := svcImpl.NewDepartmentService(departmentRepo, sfGen, logger)
	hdl := handler.NewDepartmentHandler(svc)

	return &DepartmentContainer{hdl}
}
