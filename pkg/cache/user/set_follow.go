package usercache

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

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

func (c *UserCache) ClearFollowing(ctx context.Context, userID int64) error {
	if err := c.c.Del(ctx, getFollowingKey(userID)).Err(); err != nil {
		return errors.Wrapf(err, "ClearFollowing failed, userID=%d", userID)
	}
	return nil
}
