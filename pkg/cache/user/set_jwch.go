package usercache

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func (c *UserCache) SetJwchSession(ctx context.Context, userID int64, jwchID string, cookie string) error {
	session := &jwchSession{
		ID:     jwchID,
		Cookie: cookie,
	}

	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return errors.Wrapf(err, "SetJwchSession failed, userID=%d, cookie=%s", userID, session.Cookie)
	}

	pipe := c.c.Pipeline()
	pipe.Set(ctx, getJwchSessionKey(userID), sessionJSON, constants.JwchSessionCacheExpiration)
	_, err = pipe.Exec(ctx)

	return err
}

func (c *UserCache) CleanJwchSession(ctx context.Context, userID int64) error {
	err := c.c.Del(ctx, getJwchSessionKey(userID)).Err()
	if err != nil {
		return errors.Wrapf(err, "CleanJwchSession failed, userID=%d", userID)
	}

	return nil
}
