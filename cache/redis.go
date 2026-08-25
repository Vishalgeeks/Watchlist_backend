package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()
	pong, err := Rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}
	log.Println("Redis connected:", pong)
}
