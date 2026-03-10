package videoCache

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/redis/go-redis/v9"
)

func commentCountKey(videoID int64) string {
	return "video:comment_count:" + strconv.FormatInt(videoID, 10)
}
func (v *VideoCache) SetCommentCount(videoID int64, count int64) error {
	ctx := context.Background()
	_, err := v.c.Set(ctx, commentCountKey(videoID), count, constants.CommentCountCacheExpiration).Result()
	if err != nil {
		log.Println("failed to set comment count in cache:", err)
		return err
	}
	return nil
}

func (v *VideoCache) GetCommentCount(videoID int64) (int64, error) {
	ctx := context.Background()
	result, err := v.c.Get(ctx, commentCountKey(videoID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		log.Println("failed to get comment count from cache:", err)
		return 0, err
	}

	count, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		log.Println("failed to parse comment count from cache:", err)
		return 0, err
	}

	return count, nil
}

func (v *VideoCache) IncrCommentCount(videoID int64) error {
	ctx := context.Background()
	_, err := v.c.Incr(ctx, commentCountKey(videoID)).Result()
	if err != nil {
		log.Println("failed to increment comment count in cache:", err)
		return err
	}
	return nil
}

func (v *VideoCache) DecrCommentCount(videoID int64) error {
	ctx := context.Background()
	_, err := v.c.Decr(ctx, commentCountKey(videoID)).Result()
	if err != nil {
		log.Println("failed to decrement comment count in cache:", err)
		return err
	}
	return nil
}

func (v *VideoCache) DeleteCommentCount(videoID int64) error {
	ctx := context.Background()
	_, err := v.c.Del(ctx, commentCountKey(videoID)).Result()
	if err != nil {
		log.Println("failed to delete comment count from cache:", err)
		return err
	}
	return nil
}
