package redis

import (
	"context"
	"fmt"
	"time"

	"boilerplate/internal/config"

	goRedis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var redisClient *goRedis.Client

func Init() {
	redisClient = goRedis.NewClient(&goRedis.Options{
		Addr:            fmt.Sprintf("%v:%v", config.Redis().Host, config.Redis().Port),
		Password:        config.Redis().Pass,
		DB:              0,
		PoolSize:        10,
		MaxActiveConns:  20,
		MaxIdleConns:    5,
		ConnMaxIdleTime: time.Hour,
		ConnMaxLifetime: 5 * time.Minute,
	})

	if err := redisClient.Ping(context.TODO()).Err(); err != nil {
		panic(fmt.Errorf("failed to connect to redis, error: %v", err))
	}

	logrus.Info("successfully connected to redis")
}

func Client() *goRedis.Client {
	return redisClient
}

func Close() {
	if err := redisClient.Close(); err != nil {
		logrus.WithField("message", "failed to close redis connection").Error(err.Error())
	}
	logrus.Info("redis connection to closed")
}
