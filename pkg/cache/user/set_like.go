package usercache

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func (c *UserCache) SetLikeVideos(ctx context.Context, userID int64, videoIDs []int64) error {
	args := make([]any, len(videoIDs))
	for i, v := range videoIDs {
		args[i] = v
	}

	pipe := c.c.Pipeline()
	pipe.SAdd(ctx, getLikedVideosKey(userID), args...)
	pipe.Expire(ctx, getLikedVideosKey(userID), constants.LikeCacheExpiration)
	_, err := pipe.Exec(ctx)

	if err != nil {
		return errors.Wrapf(err, "SetLikeVideos failed, userID=%d", userID)
	}

	return nil
}

func (c *UserCache) SetLikeVideo(ctx context.Context, userID int64, videoID int64) error {
	pipe := c.c.Pipeline()
	pipe.SAdd(ctx, getLikedVideosKey(userID), videoID)
	pipe.Expire(ctx, getLikedVideosKey(userID), constants.LikeCacheExpiration)
	_, err := pipe.Exec(ctx)

	if err != nil {
		return errors.Wrapf(err, "SetLikeVideo failed, userID=%d", userID)
	}

	return nil
}

func (c *UserCache) SetUnlikeVideo(ctx context.Context, userID int64, videoID int64) error {
	if err := c.c.SRem(ctx, getLikedVideosKey(userID), videoID).Err(); err != nil {
		return errors.Wrapf(err, "SetUnlikeVideo failed, userID=%d", userID)
	}
	return nil
}

func (c *UserCache) ClearLikedVideos(ctx context.Context, userID int64) error {
	if err := c.c.Del(ctx, getLikedVideosKey(userID)).Err(); err != nil {
		return errors.Wrapf(err, "ClearLikedVideos failed, userID=%d", userID)
	}
	return nil
}
