package chatcache

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (c *ChatCache) GetChatHistory(ctx context.Context, userID int64, otherUserID int64, pageSize int, pageNum int) ([]*model.ChatMessage, error) {
	data, err := c.c.Get(ctx, getHistoryKey(userID, otherUserID, pageSize, pageNum)).Bytes()
	if err != nil {
		return nil, errors.Wrapf(err, "GetChatHistory failed, userID=%d, otherUserID=%d", userID, otherUserID)
	}

	var messages []*model.ChatMessage
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
