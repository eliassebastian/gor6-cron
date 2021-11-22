package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type baseCache struct {
	cache *redis.Client
	ctx   context.Context
}

var Cache *baseCache

func InitCache(ctx context.Context) error {
	//TODO TLS Connection
	conn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	res := conn.Ping(ctx)

	if err := res.Err(); err != nil {
		return err
	}

	Cache = &baseCache{
		cache: conn,
		ctx:   ctx,
	}

	return nil
}
