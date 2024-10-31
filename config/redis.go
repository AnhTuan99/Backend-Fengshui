package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	redisHost   = os.Getenv("REDIS_HOST")
	redisPrefix = os.Getenv("REDIS_PREFIX")
)

var RedisClient *redis.Client
var Context context.Context

func ConnectToRedis() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisHost,
		DB:   0,
	})

	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Đã kết nối đến Redis:", pong)
	// Use context.Background() as the default context
	ctx := context.Background()

	RedisClient = redisClient
	Context = ctx
}

func getRedisValue(prefixedKey string) (error, string, error, bool) {
	// Thử lấy giá trị từ Redis
	val, err := RedisClient.Get(Context, prefixedKey).Result()
	if err == nil {
		// Nếu có giá trị trong Redis, trả về luôn
		return nil, val, nil, true
	} else if !errors.Is(err, redis.Nil) {
		// Xảy ra lỗi khác (nếu không phải lỗi "không tìm thấy" từ Redis), trả về lỗi
		return nil, "", err, true
	}
	return err, "", nil, false
}

func SetRedisValue(key, value string) error {
	prefixedKey := redisPrefix + key
	err := RedisClient.Set(Context, prefixedKey, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
