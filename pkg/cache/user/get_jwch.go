package usercache

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

func (c *UserCache) GetJwchSession(ctx context.Context, userID int64) (string, string, error) {
	data, err := c.c.Get(ctx, getJwchSessionKey(userID)).Result()
	if err != nil {
		return "", "", errors.Wrapf(err, "GetJwchSession failed, userID=%d", userID)
	}

	var session jwchSession
	err = json.Unmarshal([]byte(data), &session)
	if err != nil {
		return "", "", errors.Wrapf(err, "GetJwchSession failed, userID=%d", userID)
	}

	return session.ID, session.Cookie, nil
}
