package chatdao

import (
	"log"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (c *ChatDao) GetUnreadMessages(userID int64, senderID int64) ([]*model.ChatMessage, error) {
	messages, err := c.q.ChatMessage.
		Where(
			c.q.ChatMessage.SenderID.Eq(senderID),
			c.q.ChatMessage.ReceiverID.Eq(userID),
			c.q.ChatMessage.Read.Is(false),
		).
		Find()

	if err != nil {
		log.Println("failed to get unread messages for userID", userID, "from senderID", senderID, ":", err)
		return nil, err
	}

	return messages, nil
}

func (c *ChatDao) MarkMessagesAsRead(userID int64, senderID int64) error {
	_, err := c.q.ChatMessage.
		Where(
			c.q.ChatMessage.SenderID.Eq(senderID),
			c.q.ChatMessage.ReceiverID.Eq(userID),
			c.q.ChatMessage.Read.Is(false),
		).
		Update(c.q.ChatMessage.Read, true)

	if err != nil {
		log.Println("failed to mark messages as read for userID", userID, "from senderID", senderID, ":", err)
		return err
	}
	return nil
}

func (c *ChatDao) GetChatHistory(userID int64, otherUserID int64, pageSize int, pageNum int) ([]*model.ChatMessage, error) {
	messages, err := c.q.ChatMessage.
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
