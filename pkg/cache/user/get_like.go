package usercache

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func (c *UserCache) GetLikedVideos(ctx context.Context, userID int64) ([]int64, error) {
	videoIDsStr, err := c.c.SMembers(ctx, getLikedVideosKey(userID)).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "GetLikedVideos failed, userID=%d", userID)
	}

	videoIDs := make([]int64, len(videoIDsStr))
	for i, v := range videoIDsStr {
		var id int64
		_, err := fmt.Sscanf(v, "%d", &id)
		if err != nil {
			return nil, errors.Wrapf(err, "GetLikedVideos failed, userID=%d", userID)
		}
		videoIDs[i] = id
	}

	return videoIDs, nil
}

func (c *UserCache) IsVideoLiked(ctx context.Context, userID int64, videoID int64) (bool, error) {
	exist, err := c.c.Exists(ctx, getFollowingKey(userID)).Result()
	if err != nil {
		return false, errors.Wrapf(err, "IsVideoLiked failed, userID=%d", userID)
	}
	if exist == 0 {
		return false, redis.Nil
	}

	result, err := c.c.SIsMember(ctx, getLikedVideosKey(userID), videoID).Result()
	if err != nil {
		return false, errors.Wrapf(err, "IsVideoLiked failed, userID=%d", userID)
	}

	return result, nil
}
