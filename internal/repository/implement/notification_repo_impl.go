package implement

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"gorm.io/gorm"
)

type notificationRepoImpl struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) repository.Notification {
	return &notificationRepoImpl{db}
}

func (r *notificationRepoImpl) CreateNotification(ctx context.Context, notification *model.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}