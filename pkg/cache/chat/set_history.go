package chatcache

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (c *ChatCache) SetChatHistory(ctx context.Context, userID int64, otherUserID int64, pageSize int, pageNum int, messages []*model.ChatMessage) error {
	data, err := json.Marshal(messages)
	if err != nil {
		return errors.Wrapf(err, "SetChatHistory failed, marshal message failed, userID=%d, otherUserID=%d", userID, otherUserID)
	}

	return c.c.Set(ctx, getHistoryKey(userID, otherUserID, pageSize, pageNum), data, constants.ChatHistoryCacheExpiration).Err()
}

func (c *ChatCache) ClearChatHistory(ctx context.Context, userID int64, otherUserID int64) error {
	return c.deleteByPattern(ctx, getHistoryKeyPattern(userID, otherUserID))
}
