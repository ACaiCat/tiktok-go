package chatcache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func getUnreadKey(userID int64, senderID int64) string {
	return fmt.Sprintf("chat:unread:%d:%d", userID, senderID)
}

func (c *ChatCache) SetUnreadMessages(ctx context.Context, userID int64, senderID int64, messages []*model.ChatMessage) error {
	data, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	return c.c.Set(ctx, getUnreadKey(userID, senderID), data, constants.ChatUnreadCacheExpiration).Err()
}

func (c *ChatCache) GetUnreadMessages(ctx context.Context, userID int64, senderID int64) ([]*model.ChatMessage, error) {
	data, err := c.c.Get(ctx, getUnreadKey(userID, senderID)).Bytes()
	if err != nil {
		return nil, err
	}

	var messages []*model.ChatMessage
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (c *ChatCache) ClearUnreadMessages(ctx context.Context, userID int64, senderID int64) error {
	return c.c.Del(ctx, getUnreadKey(userID, senderID)).Err()
}
