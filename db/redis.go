package db

import (
	"context"
	"log"
	"strings"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis(address, port, password string, db int) {
	var err error
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     strings.Join([]string{address, ":", port}, ""),
		Password: password,
		DB:       db,
	})
	log.Println("db redis client", RedisClient)
	if _, err = RedisClient.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}
}
