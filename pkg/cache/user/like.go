package usercache

import (
	"context"
	"fmt"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func getLikedVideosKey(userID int64) string {
	return fmt.Sprintf("user:%d:liked_videos", userID)
}

func (c *UserCache) SetLikeVideos(userID int64, videoIDs []int64) error {
	ctx := context.Background()

	args := make([]interface{}, len(videoIDs))
	for i, v := range videoIDs {
		args[i] = v
	}

	pipe := c.c.Pipeline()
	pipe.SAdd(ctx, getLikedVideosKey(userID), args...)
	pipe.Expire(ctx, getLikedVideosKey(userID), constants.LikeCacheExpiration)
	_, err := pipe.Exec(ctx)

	return err
}

func (c *UserCache) GetLikedVideos(userID int64) ([]int64, error) {
	ctx := context.Background()

	videoIDsStr, err := c.c.SMembers(ctx, getLikedVideosKey(userID)).Result()
	if err != nil {
		return nil, err
	}

	videoIDs := make([]int64, len(videoIDsStr))
	for i, v := range videoIDsStr {
		var id int64
		_, err := fmt.Sscanf(v, "%d", &id)
		if err != nil {
			return nil, err
		}
		videoIDs[i] = id
	}

	return videoIDs, nil
}

func (c *UserCache) SetLikeVideo(userID int64, videoID int64) error {
	ctx := context.Background()
	pipe := c.c.Pipeline()
	pipe.SAdd(ctx, getLikedVideosKey(userID), videoID)
	pipe.Expire(ctx, getLikedVideosKey(userID), constants.LikeCacheExpiration)
	_, err := pipe.Exec(ctx)
	return err
}

func (c *UserCache) SetUnlikeVideo(userID int64, videoID int64) error {
	ctx := context.Background()

	return c.c.SRem(ctx, getLikedVideosKey(userID), videoID).Err()
}

func (c *UserCache) IsVideoLiked(userID int64, videoID int64) (bool, error) {
	ctx := context.Background()

	return c.c.SIsMember(ctx, getLikedVideosKey(userID), videoID).Result()
}

func (c *UserCache) ClearLikedVideos(userID int64) error {
	ctx := context.Background()

	return c.c.Del(ctx, getLikedVideosKey(userID)).Err()
}
