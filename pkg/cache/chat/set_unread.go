package chatcache

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (c *ChatCache) SetUnreadMessages(ctx context.Context, userID int64, senderID int64, messages []*model.ChatMessage) error {
	data, err := json.Marshal(messages)
	if err != nil {
		return errors.Wrapf(err, "SetUnreadMessages failed, userID=%d, otherUserID=%d", userID, senderID)
	}

	return c.c.Set(ctx, getUnreadKey(userID, senderID), data, constants.ChatUnreadCacheExpiration).Err()
}

func (c *ChatCache) ClearUnreadMessages(ctx context.Context, userID int64, senderID int64) error {
	if err := c.c.Del(ctx, getUnreadKey(userID, senderID)).Err(); err != nil {
		return errors.Wrapf(err, "ClearUnreadMessages failed, userID=%d, otherUserID=%d", userID, senderID)
	}
	return nil
}
