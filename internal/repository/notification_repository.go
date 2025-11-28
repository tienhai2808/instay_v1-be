package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/types"
	"gorm.io/gorm"
)

type Notification interface {
	CreateNotificationTx(ctx context.Context, tx *gorm.DB, notification *model.Notification) error

	FindAllUnreadNotificationsByContentIDAndType(ctx context.Context, staffID, contentID int64, contentType string) ([]*model.Notification, error)

	CreateNotificationStaffs(ctx context.Context, notificationStaffs []*model.NotificationStaff) error

	FindAllUnreadNotificationsByDepartmentID(ctx context.Context, staffID, departmentID int64) ([]*model.Notification, error)

	FindAllUnreadNotificationsByOrderRoomID(ctx context.Context, orderRoomID int64) ([]*model.Notification, error)

	FindAllNotificationsByOrderRoomID(ctx context.Context, orderRoomID int64) ([]*model.Notification, error)

	CountUnreadNotificationsByDepartmentID(ctx context.Context, userID, departmentID int64) (int64, error)

	CountUnreadNotificationsByOrderRoomID(ctx context.Context, orderRoomID int64) (int64, error)

	UpdateNotifications(ctx context.Context, notificationIDs []int64, updateData map[string]any) error

	UpdateNotificationsByOrderRoomIDAndType(ctx context.Context, orderRoomID int64, contentType string, updateData map[string]any) error

	FindAllNotificationsByDepartmentIDWithStaffsReadPaginated(ctx context.Context, query types.NotificationPaginationQuery, staffID, departmentID int64) ([]*model.Notification, int64, error)
}
