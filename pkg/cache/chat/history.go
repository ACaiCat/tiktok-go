package chatcache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func normalizeConversationUserIDs(userID int64, otherUserID int64) (int64, int64) {
	if userID > otherUserID {
		return otherUserID, userID
	}
	return userID, otherUserID
}

func getHistoryKey(userID int64, otherUserID int64, pageSize int, pageNum int) string {
	left, right := normalizeConversationUserIDs(userID, otherUserID)
	return fmt.Sprintf("chat:history:%d:%d:%d:%d", left, right, pageSize, pageNum)
}

func getHistoryKeyPattern(userID int64, otherUserID int64) string {
	left, right := normalizeConversationUserIDs(userID, otherUserID)
	return fmt.Sprintf("chat:history:%d:%d:*", left, right)
}

func (c *ChatCache) SetChatHistory(ctx context.Context, userID int64, otherUserID int64, pageSize int, pageNum int, messages []*model.ChatMessage) error {
	data, err := json.Marshal(messages)
	if err != nil {
		return errors.Wrapf(err, "SetChatHistory failed, marshal message failed, userID=%d, otherUserID=%d", userID, otherUserID)
	}

	return c.c.Set(ctx, getHistoryKey(userID, otherUserID, pageSize, pageNum), data, constants.ChatHistoryCacheExpiration).Err()
}

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

func (c *ChatCache) ClearChatHistory(ctx context.Context, userID int64, otherUserID int64) error {
	return c.deleteByPattern(ctx, getHistoryKeyPattern(userID, otherUserID))
}
