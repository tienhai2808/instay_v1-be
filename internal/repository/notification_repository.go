package repository

import (
	"context"

	"github.com/InstaySystem/is-be/internal/model"
)

type Notification interface {
	CreateNotification(ctx context.Context, notification *model.Notification) error
}