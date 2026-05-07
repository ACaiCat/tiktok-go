package chatcache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type ChatCache struct {
	c *redis.Client
}

func NewChatCache(cache *redis.Client) *ChatCache {
	return &ChatCache{c: cache}
}

func (c *ChatCache) deleteByPattern(ctx context.Context, pattern string) error {
	var (
		cursor uint64
		keys   []string
	)

	var patternScanBatch int64 = 100

	for {
		batch, nextCursor, err := c.c.Scan(ctx, cursor, pattern, patternScanBatch).Result()
		if err != nil {
			return err
		}
		keys = append(keys, batch...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	if len(keys) == 0 {
		return nil
	}

	return c.c.Del(ctx, keys...).Err()
}
