package cache

import (
	"context"
	"strconv"

	"github.com/ACaiCat/tiktok-go/config"
	"github.com/redis/go-redis/v9"
)

var Cache *redis.Client

func InitRedis() {
	var err error
	Cache = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.Redis.Host + ":" + strconv.Itoa(config.AppConfig.Redis.Port),
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	})

	err = Cache.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

}
