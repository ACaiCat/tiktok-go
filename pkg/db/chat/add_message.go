package chatdao

import (
	"log"
	"time"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (c *ChatDao) AddMessage(senderID int64, receiverID int64, content string, isRead bool) error {
	var err error

	var readAt *time.Time

	if isRead {
		readAt = new(time.Now())
	}

	message := model.ChatMessage{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		ReadAt:     readAt,
	}

	err = c.q.ChatMessage.Create(&message)
	if err != nil {
		log.Println("failed to add message from userID", senderID, "to userID", receiverID, ":", err)
		return err
	}

	return nil
}
