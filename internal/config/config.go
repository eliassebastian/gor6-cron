package config

import "github.com/go-redis/redis/v8"

type AppConfig struct {
	RedisDB *redis.Client
}

var Config *AppConfig
