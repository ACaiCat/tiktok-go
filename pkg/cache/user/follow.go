package usercache

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func getFollowingKey(userID int64) string {
	return fmt.Sprintf("user:%d:following", userID)
}

func (c *UserCache) SetFollowings(ctx context.Context, userID int64, followingIDs []int64) error {
	args := make([]any, len(followingIDs))
	for i, v := range followingIDs {
		args[i] = v
	}

	pipe := c.c.Pipeline()
	pipe.SAdd(ctx, getFollowingKey(userID), args...)
	pipe.Expire(ctx, getFollowingKey(userID), constants.LikeCacheExpiration)
	_, err := pipe.Exec(ctx)

	if err != nil {
		return errors.Wrapf(err, "SetFollowings failed, userID=%d", userID)
	}

	return nil
}

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

func (c *UserCache) SetFollow(ctx context.Context, userID int64, followingID int64) error {
	pipe := c.c.Pipeline()
	pipe.SAdd(ctx, getFollowingKey(userID), followingID)
	pipe.Expire(ctx, getFollowingKey(userID), constants.LikeCacheExpiration)
	_, err := pipe.Exec(ctx)

	if err != nil {
		return errors.Wrapf(err, "SetFollow failed, userID=%d", userID)
	}

	return nil
}

func (c *UserCache) SetUnfollow(ctx context.Context, userID int64, followingID int64) error {
	if err := c.c.SRem(ctx, getFollowingKey(userID), followingID).Err(); err != nil {
		return errors.Wrapf(err, "SetUnfollow failed, userID=%d", userID)
	}
	return nil
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

func (c *UserCache) ClearFollowing(ctx context.Context, userID int64) error {
	if err := c.c.Del(ctx, getFollowingKey(userID)).Err(); err != nil {
		return errors.Wrapf(err, "ClearFollowing failed, userID=%d", userID)
	}
	return nil
}
