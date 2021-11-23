package config

import (
	"github.com/go-redis/redis/v8"
	"net/http"
)

type AppConfig struct {
	RedisDB *redis.Client
	Client  *http.Client
}

var Config *AppConfig
