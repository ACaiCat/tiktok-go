package chatdao

import (
	"context"
	"time"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
	"github.com/pkg/errors"
)

func (c *ChatDao) AddMessage(ctx context.Context, senderID int64, receiverID int64, content string, isRead bool, isAi bool) error {
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
		IsAi:       isAi,
	}

	err = c.q.ChatMessage.WithContext(ctx).Create(&message)
	if err != nil {
		return errors.Wrapf(err, "AddMessage failed, senderID: %d, receiverID: %d", message.SenderID, message.ReceiverID)
	}

	return nil
}
