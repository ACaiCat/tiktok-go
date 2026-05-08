package usercache

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func getLikedVideosKey(userID int64) string {
	return fmt.Sprintf("user:%d:liked_videos", userID)
}

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
		return errors.Wrapf(err, "SetLikeVideosFailed, userID=%d", userID)
	}

	return nil
}

func (c *UserCache) GetLikedVideos(ctx context.Context, userID int64) ([]int64, error) {
	videoIDsStr, err := c.c.SMembers(ctx, getLikedVideosKey(userID)).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "get liked videos failed, userID=%d", userID)
	}

	videoIDs := make([]int64, len(videoIDsStr))
	for i, v := range videoIDsStr {
		var id int64
		_, err := fmt.Sscanf(v, "%d", &id)
		if err != nil {
			return nil, errors.Wrapf(err, "GetLikedVideosById failed, userID=%d", userID)
		}
		videoIDs[i] = id
	}

	return videoIDs, nil
}

func (c *UserCache) SetLikeVideo(ctx context.Context, userID int64, videoID int64) error {
	pipe := c.c.Pipeline()
	pipe.SAdd(ctx, getLikedVideosKey(userID), videoID)
	pipe.Expire(ctx, getLikedVideosKey(userID), constants.LikeCacheExpiration)
	_, err := pipe.Exec(ctx)

	if err != nil {
		return errors.Wrapf(err, "SetLikedVideosById failed, userID=%d", userID)
	}

	return nil
}

func (c *UserCache) SetUnlikeVideo(ctx context.Context, userID int64, videoID int64) error {
	if err := c.c.Del(ctx, getLikedVideosKey(userID)).Err(); err != nil {
		return errors.Wrapf(err, "CleanLikedVideosById failed, userID=%d", userID)
	}
	return nil
}

func (c *UserCache) IsVideoLiked(ctx context.Context, userID int64, videoID int64) (bool, error) {
	result, err := c.c.SIsMember(ctx, getLikedVideosKey(userID), videoID).Result()
	if err != nil {
		return false, errors.Wrapf(err, "IsLikedVideosById failed, userID=%d", userID)
	}

	return result, nil
}

func (c *UserCache) ClearLikedVideos(ctx context.Context, userID int64) error {
	if err := c.c.Del(ctx, getLikedVideosKey(userID)).Err(); err != nil {
		return errors.Wrapf(err, "CleanLikedVideosById failed, userID=%d", userID)
	}
	return nil
}
