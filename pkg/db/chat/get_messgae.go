package chatdao

import (
	"context"
	"log"
	"time"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
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
		log.Println("failed to get unread messages for userID", userID, "from senderID", senderID, ":", err)
		return nil, err
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
		log.Println("failed to mark messages as read for userID", userID, "from senderID", senderID, ":", err)
		return err
	}
	return nil
}

func (c *ChatDao) GetChatHistory(ctx context.Context, userID int64, otherUserID int64, pageSize int, pageNum int) ([]*model.ChatMessage, error) {
	messages, err := c.q.ChatMessage.WithContext(ctx).
		Where(
			c.q.ChatMessage.SenderID.In(userID, otherUserID),
			c.q.ChatMessage.ReceiverID.In(userID, otherUserID),
		).
		Order(c.q.ChatMessage.CreatedAt.Desc()).
		Offset(pageSize * pageNum).
		Limit(pageSize).
		Find()

	if err != nil {
		log.Println("failed to get chat history between userID", userID, "and otherUserID", otherUserID, ":", err)
		return nil, err
	}

	return messages, nil
}
