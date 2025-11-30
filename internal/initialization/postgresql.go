package initialization

import (
	"database/sql"
	"fmt"

	"github.com/InstaySystem/is-be/internal/config"
	"github.com/InstaySystem/is-be/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var allModels = []any{
	&model.User{},
	&model.Department{},
	&model.ServiceType{},
	&model.Service{},
	&model.ServiceImage{},
	&model.RequestType{},
	&model.Request{},
	&model.RoomType{},
	&model.Floor{},
	&model.Room{},
	&model.Booking{},
	&model.OrderRoom{},
	&model.OrderService{},
	&model.Notification{},
	&model.NotificationStaff{},
	&model.Chat{},
	&model.Message{},
	&model.MessageStaff{},
	&model.Review{},
}

type DB struct {
	Gorm *gorm.DB
	sql  *sql.DB
}

func InitPostgreSQL(cfg *config.Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.DBName,
		cfg.PostgreSQL.User,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.SSLMode,
	)
	gDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("database - %w", err)
	}

	if err := runAutoMigrations(gDB); err != nil {
		return nil, fmt.Errorf("database - %w", err)
	}

	sqlDB, err := gDB.DB()
	if err != nil {
		return nil, fmt.Errorf("database - %w", err)
	}

	return &DB{
		gDB,
		sqlDB,
	}, nil
}

func (d *DB) Close() {
	_ = d.sql.Close()
}

func runAutoMigrations(db *gorm.DB) error {
	return db.AutoMigrate(allModels...)
}
