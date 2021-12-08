package misc

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strings"
)

var RedisClient *redis.Client

func initRedis() {
	var err error
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     strings.Join([]string{Conf.RedisCfg.Address, ":", Conf.RedisCfg.Port}, ""),
		Password: Conf.RedisCfg.Password,
		DB:       Conf.RedisCfg.Db,
	})
	Logger.Info("init redis client", zap.Any("redis", RedisClient))
	if _, err = RedisClient.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}
}
