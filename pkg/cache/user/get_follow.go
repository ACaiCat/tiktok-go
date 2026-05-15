package usercache

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func (c *UserCache) GetFollowing(ctx context.Context, userID int64) ([]int64, error) {
	videoIDsStr, err := c.c.SMembers(ctx, getFollowingKey(userID)).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "GetFollowing failed, userID=%d", userID)
	}

	videoIDs := make([]int64, len(videoIDsStr))
	for i, v := range videoIDsStr {
		var id int64
		_, err := fmt.Sscanf(v, "%d", &id)
		if err != nil {
			return nil, errors.Wrapf(err, "GetFollowing failed, userID=%d", userID)
		}
		videoIDs[i] = id
	}

	return videoIDs, nil
}

func (c *UserCache) IsFollowed(ctx context.Context, userID int64, followingID int64) (bool, error) {
	exist, err := c.c.Exists(ctx, getFollowingKey(userID)).Result()
	if err != nil {
		return false, errors.Wrapf(err, "IsFollowed failed, userID=%d", userID)
	}
	if exist == 0 {
		return false, redis.Nil
	}

	result, err := c.c.SIsMember(ctx, getFollowingKey(userID), followingID).Result()
	if err != nil {
		return false, errors.Wrapf(err, "IsFollowed failed, userID=%d", userID)
	}

	return result, nil
}
