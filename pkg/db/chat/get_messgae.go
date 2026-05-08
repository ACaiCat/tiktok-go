package chatdao

import (
	"context"
	"time"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
	"github.com/pkg/errors"
)

func (c *ChatDao) GetUnreadMessages(ctx context.Context, userID int64, senderID int64) ([]*model.ChatMessage, error) {
	messages, err := c.q.ChatMessage.WithContext(ctx).
		Where(
			c.q.ChatMessage.SenderID.Eq(senderID),
			c.q.ChatMessage.ReceiverID.Eq(userID),
			c.q.ChatMessage.ReadAt.IsNull(),
		).
		Find()

	if err != nil {
		return nil, errors.Wrapf(err, "GetUnreadMessages failed, senderID: %d, receiverID: %d", userID, senderID)
	}

	return messages, nil
}

func (c *ChatDao) MarkMessagesAsRead(ctx context.Context, userID int64, senderID int64) error {
	_, err := c.q.ChatMessage.WithContext(ctx).
		Where(
			c.q.ChatMessage.SenderID.Eq(senderID),
			c.q.ChatMessage.ReceiverID.Eq(userID),
			c.q.ChatMessage.ReadAt.IsNull(),
		).
		Update(c.q.ChatMessage.ReadAt, time.Now())

	if err != nil {
		return errors.Wrapf(err, "MarkMessagesAsRead failed, senderID: %d, receiverID: %d", userID, senderID)
	}
	return nil
}

func (c *ChatDao) GetChatHistory(ctx context.Context, userID int64, otherUserID int64, pageSize int, pageNum int) ([]*model.ChatMessage, error) {
	messages, err := c.q.ChatMessage.WithContext(ctx).
		Where(
			c.q.ChatMessage.SenderID.In(userID, otherUserID),
			c.q.ChatMessage.ReceiverID.In(userID, otherUserID),
		).
		Where(c.q.ChatMessage.Or(c.q.ChatMessage.IsAi.Is(false)).Or(c.q.ChatMessage.SenderID.Eq(otherUserID))).
		Order(c.q.ChatMessage.CreatedAt.Desc()).
		Offset(pageSize * pageNum).
		Limit(pageSize).
		Find()

	if err != nil {
		return nil, errors.Wrapf(err, "GetChatHistory failed, senderID: %d, receiverID: %d", userID, otherUserID)
	}

	return messages, nil
}
