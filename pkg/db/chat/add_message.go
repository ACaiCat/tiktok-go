package chatdao

import (
	"log"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (c *ChatDao) AddMessage(senderID int64, receiverID int64, content string, read bool) error {
	var err error

	message := model.ChatMessage{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		Read:       read,
	}

	err = c.q.ChatMessage.Create(&message)
	if err != nil {
		log.Println("failed to add message from userID", senderID, "to userID", receiverID, ":", err)
		return err
	}

	return nil
}
