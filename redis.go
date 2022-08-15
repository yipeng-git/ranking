package main

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

)

var redisCli *redis.Client

func initRedis() {
	redisCli = redis.NewClient(&redis.Options{
		Addr:         "0.0.0.0:6379",
		DialTimeout:  time.Duration(8) * time.Second,
		ReadTimeout:  time.Duration(8) * time.Second,
		WriteTimeout: time.Duration(8) * time.Second,
		DB:           0,
	})
	_, _ = redisCli.Ping(context.Background()).Result()
	return
}
