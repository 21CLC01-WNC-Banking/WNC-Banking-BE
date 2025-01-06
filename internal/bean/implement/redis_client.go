package beanimplement

import (
	"context"
	"fmt"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/env"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/bean"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/constants"
	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService() bean.RedisClient {
	redisHost, err := env.GetEnv("REDIS_HOST")
	redisPort, err := env.GetEnv("REDIS_PORT")
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "",
		DB:       0,
	})

	return &RedisService{client: client}
}

func (r *RedisService) Set(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, key, value, constants.RESET_PASSWORD_EXP_TIME).Err()
}

func (r *RedisService) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (r *RedisService) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
