package initialization

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/InstaySystem/is-be/internal/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg *config.Config) (*redis.Client, error) {
	rAddr := cfg.Redis.Host + fmt.Sprintf(":%d", cfg.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:      rAddr,
		Password:  cfg.Redis.Password,
		TLSConfig: &tls.Config{},
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
