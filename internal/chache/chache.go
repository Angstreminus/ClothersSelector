package chache

import (
	"strconv"

	"github.com/Angstreminus/ClothersSelector/config"
	"github.com/redis/go-redis/v9"
)

type Chache struct {
	RedisChahe *redis.Client
	Config     *config.Config
}

func NewChache(cfg *config.Config) (*Chache, error) {
	dbNum, err := strconv.Atoi(cfg.RedisDatabase)
	if err != nil {
		return nil, err
	}
	return &Chache{
		RedisChahe: redis.NewClient(
			&redis.Options{
				Addr:     cfg.RedisAddr,
				Password: cfg.RedisPassword,
				DB:       dbNum,
			},
		),
		Config: cfg,
	}, nil
}
