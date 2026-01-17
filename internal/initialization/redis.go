package initialization

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/InstaySystem/is_v1-be/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

func InitRedis(cfg *config.Config) (*redis.Client, error) {
	rAddr := cfg.Redis.Host + fmt.Sprintf(":%d", cfg.Redis.Port)

	options := &redis.Options{
		Addr:     rAddr,
		Password: cfg.Redis.Password,
		DB:       0,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	}
	if cfg.Redis.UseSSL {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	rdb := redis.NewClient(options)
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("cache - %w", err)
	}

	return rdb, nil
}
