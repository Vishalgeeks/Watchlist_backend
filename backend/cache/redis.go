package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

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
		log.Println("Redis connection failed (continuing without cache):", err)
		return
	}
	log.Println("Redis connected:", pong)
}

// GetJSON implements the cache-aside read: on a cache hit it decodes the
// stored JSON into dest and returns (true, nil). On a miss (or if Redis is
// unavailable) it returns (false, nil) so the caller falls back to PostgreSQL.
func GetJSON(ctx context.Context, key string, dest interface{}) (bool, error) {
	if Rdb == nil {
		return false, nil
	}
	val, err := Rdb.Get(ctx, key).Result()
	if err != nil {
		// Cache miss or Redis down — graceful fallback is the caller's job.
		return false, nil
	}
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false, nil
	}
	return true, nil
}

// SetJSON stores val serialized as JSON under key with the given TTL.
// Errors are non-fatal: the request still succeeds without caching.
func SetJSON(ctx context.Context, key string, val interface{}, ttl time.Duration) {
	if Rdb == nil {
		return
	}
	data, err := json.Marshal(val)
	if err != nil {
		return
	}
	Rdb.Set(ctx, key, data, ttl)
}

// Del removes one or more keys. Used to invalidate caches after mutations.
// Errors are non-fatal: the keys simply expire naturally.
func Del(ctx context.Context, keys ...string) {
	if Rdb == nil || len(keys) == 0 {
		return
	}
	Rdb.Del(ctx, keys...)
}
