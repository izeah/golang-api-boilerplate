package config

import (
	"os"
	"sync"

	"boilerplate/pkg/util/priority"
)

type RedisConfig struct {
	Host string
	Port string
	Pass string
}

var (
	redisConfig *RedisConfig
	redisOnce   sync.Once
)

func Redis() *RedisConfig {
	redisOnce.Do(func() {
		redisConfig = &RedisConfig{
			Host: priority.PriorityString(os.Getenv("REDIS_HOST"), "localhost"),
			Port: priority.PriorityString(os.Getenv("REDIS_PORT"), "6379"),
			Pass: priority.PriorityString(os.Getenv("REDIS_PASSWORD"), ""),
		}
	})
	return redisConfig
}
