package redis

import (
	"auth-service/pkg/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func ConnectDB() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return rdb
}

type RedisStorage struct {
	rdb *redis.Client
	log *slog.Logger
}

func NewRedisStorage(rdb *redis.Client, logger *slog.Logger) *RedisStorage {
	return &RedisStorage{rdb: rdb, log: logger}
}

func (r *RedisStorage) SetCode(ctx context.Context, email, code string) error {

	err := r.rdb.Set(ctx, email, code, 5*time.Minute).Err()
	if err != nil {
		return errors.Wrap(err, "failed to set code in Redis")
	}

	return nil
}

func (r *RedisStorage) GetCodes(ctx context.Context, email string) (string, error) {
	code, err := r.rdb.Get(ctx, email).Result()
	if err != nil {
		if err == redis.Nil {
			r.log.Error("no code found for email", "email", email)
			return "", fmt.Errorf("no code found for email: %s", email)
		}
		return "", errors.Wrap(err, "failed to get code from Redis")
	}
	return code, nil
}

func (r *RedisStorage) SetRegister(ctx context.Context, in models.RegisterRequest1) error {

	info, err := json.Marshal(in)
	if err != nil {
		r.log.Error("Error in Marshal", "err", err)
		return err
	}

	err = r.rdb.Set(ctx, in.Email, info, 5*time.Minute).Err()
	if err != nil {
		r.log.Error("Error in Set", "err", err)
		return errors.Wrap(err, "failed to set code in Redis")
	}

	return nil
}

func (r *RedisStorage) GetRegister(ctx context.Context, email string) (models.RegisterRequest1, error) {

	res, err := r.rdb.Get(ctx, email).Result()
	if err != nil {
		r.log.Error("Error in GetRegister", "err", err)
		return models.RegisterRequest1{}, err
	}

	var response models.RegisterRequest1

	err = json.Unmarshal([]byte(res), &response)
	if err != nil {
		r.log.Error("Error in Unmarshal", "err", err)
		return models.RegisterRequest1{}, err
	}

	return response, nil
}
