package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client redis.Cmdable // use cmdable for easier usage of a mock redis server while testing
	// maybe check with performance later on
	ctx context.Context
}

var instance *RedisClient

// var once sync.Once

func Initialize(client redis.Cmdable) *RedisClient {
	instance = &RedisClient{
		client: client,
		ctx:    context.Background(),
	}
	return instance
}

func Connect() *redis.Client {
	fmt.Print("Connecting to redis...")
	address := os.Getenv("REDIS_ADDRESS")
	password := os.Getenv("REDIS_PASSWORD")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("error while parsing db value from .env file: %s", err)
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	// check redis connection

	return client
}

// set key-value pair
func (r *RedisClient) Set(key, value string) error {
	return r.client.Set(r.ctx, key, value, 0).Err()
}

// set key-value pair with ttl
func (r *RedisClient) SetWithTTL(key, value string, ttl time.Duration) error {
	return r.client.Set(r.ctx, key, value, ttl).Err()
}

// get value by key
func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

// delete by key
func (r *RedisClient) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// close redis client connection
func (r *RedisClient) Close() error {
	if client, ok := r.client.(*redis.Client); ok {
		return client.Close()
	}
	return nil
}
