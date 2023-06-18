package redis_provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/huytq/tech-core/logic"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}
type IRedisProvider interface {
	SetKey(ctx context.Context, request SetRedisValue) error
	GetKey(ctx context.Context, key string, data interface{}) error
}

type redisProvider struct {
	redisClient *redis.Client
}

func NewRedisProvider(config RedisConfig) (IRedisProvider, func()) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password, // no password set
		DB:       config.DB,       // use default DB
	})
	return &redisProvider{redisClient: rdb}, func() {
		rdb.Close()
	}
}

type SetRedisValue struct {
	Key     string
	Value   interface{}
	Expired time.Duration
}

func (redis *redisProvider) SetKey(ctx context.Context, request SetRedisValue) error {

	valueStr, err := json.Marshal(request.Value)
	if err != nil {
		return err
	}
	err = redis.redisClient.Set(ctx, request.Key, valueStr, request.Expired).Err()
	if err != nil {
		return err
	}

	return nil
}

func (redis *redisProvider) GetKey(ctx context.Context, key string, data interface{}) error {

	value := redis.redisClient.Get(ctx, key)
	if value == nil {
		return errors.New("Error Get Key Redis")
	}
	if value.Err() != nil {
		return value.Err()
	}

	if result, err := value.Result(); err != nil {
		return err
	} else {
		if err = logic.ParseStringToStruct(result, &data); err != nil {
			return err
		}
	}

	return nil
}
