package container

import (
	"github.com/InstaySystem/is-be/internal/repository"
	repoImpl "github.com/InstaySystem/is-be/internal/repository/implement"
	"gorm.io/gorm"
)

type BookingContainer struct {
	Repo repository.BookingRepository
}

func NewBookingContainer(db *gorm.DB) *BookingContainer {
	repo := repoImpl.NewBookingRepository(db)
	
	return &BookingContainer{repo}
}