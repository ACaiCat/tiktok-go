package videoCache

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/redis/go-redis/v9"
)

func likeCountKey(videoID int64) string {
	return "video:like_count:" + strconv.FormatInt(videoID, 10)
}

func (v *VideoCache) SetLikeCount(videoID int64, count int64) error {
	ctx := context.Background()
	err := v.c.Set(ctx, likeCountKey(videoID), strconv.FormatInt(count, 10), constants.LikeCountCacheExpiration).Err()
	if err != nil {
		log.Println("failed to set like count in cache:", err)
		return err
	}
	return nil
}

func (v *VideoCache) GetLikeCount(videoID int64) (int64, error) {
	ctx := context.Background()
	result, err := v.c.Get(ctx, likeCountKey(videoID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}

		log.Println("failed to get like count from cache:", err)
		return 0, err
	}
	count, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		log.Println("failed to parse like count:", err)
		return 0, err
	}
	return count, nil
}

func (v *VideoCache) IncrLikeCount(videoID int64) error {
	ctx := context.Background()
	err := v.c.Incr(ctx, likeCountKey(videoID)).Err()
	if err != nil {
		log.Println("failed to increment like count in cache:", err)
		return err
	}
	return nil
}

func (v *VideoCache) DecrLikeCount(videoID int64) error {
	ctx := context.Background()
	err := v.c.Decr(ctx, likeCountKey(videoID)).Err()
	if err != nil {
		log.Println("failed to decrement like count in cache:", err)
		return err
	}
	return nil
}

func (v *VideoCache) DeleteLikeCount(videoID int64) error {
	ctx := context.Background()
	err := v.c.Del(ctx, likeCountKey(videoID)).Err()
	if err != nil {
		log.Println("failed to delete like count from cache:", err)
		return err
	}
	return nil
}
