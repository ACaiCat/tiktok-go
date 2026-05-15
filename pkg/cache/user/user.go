package usercache

import (
	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	c *redis.Client
}

type jwchSession struct {
	ID     string
	Cookie string
}

func NewUserCache(cache *redis.Client) *UserCache {
	return &UserCache{c: cache}
}
