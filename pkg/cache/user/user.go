package userCache

import (
	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	c *redis.Client
}

func NewUserCache(cache *redis.Client) *UserCache {
	return &UserCache{c: cache}
}
