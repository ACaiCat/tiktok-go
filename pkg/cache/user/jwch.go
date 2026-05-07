package usercache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func getJwchSessionKey(userID int64) string {
	return fmt.Sprintf("user:%d:jwch", userID)
}

type jwchSession struct {
	ID     string
	Cookie string
}

func (c *UserCache) SetJwchSession(ctx context.Context, userID int64, jwchID string, cookie string) error {
	session := &jwchSession{
		ID:     jwchID,
		Cookie: cookie,
	}

	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	}

	pipe := c.c.Pipeline()
	pipe.Set(ctx, getJwchSessionKey(userID), sessionJSON, constants.JwchSessionCacheExpiration)
	_, err = pipe.Exec(ctx)

	return err
}

func (c *UserCache) GetJwchSession(ctx context.Context, userID int64) (string, string, error) {
	data, err := c.c.Get(ctx, getJwchSessionKey(userID)).Result()
	if err != nil {
		return "", "", err
	}

	var session jwchSession
	err = json.Unmarshal([]byte(data), &session)
	if err != nil {
		return "", "", err
	}

	return session.ID, session.Cookie, nil
}

func (c *UserCache) CleanJwchSession(ctx context.Context, userID int64) error {
	err := c.c.Del(ctx, getJwchSessionKey(userID)).Err()
	if err != nil {
		return err
	}

	return nil
}
