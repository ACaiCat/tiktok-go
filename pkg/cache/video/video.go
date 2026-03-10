package videoCache

import "github.com/redis/go-redis/v9"

type VideoCache struct {
	c *redis.Client
}

func NewVideoCache(cache *redis.Client) *VideoCache {
	return &VideoCache{c: cache}
}
