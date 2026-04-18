package service

import (
	"github.com/ACaiCat/tiktok-go/api/model/ws"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func MessageDaoToDto(message *model.ChatMessage) *ws.ChatMessage {
	return &ws.ChatMessage{
		SenderID:   message.SenderID,
		ReceiverID: message.ReceiverID,
		Content:    message.Content,
		Timestamp:  message.CreatedAt.Unix(),
	}
}

func MessagesDaoToDto(messages []*model.ChatMessage) []*ws.ChatMessage {
	chatMessages := make([]*ws.ChatMessage, len(messages))
	for i, message := range messages {
		chatMessages[i] = MessageDaoToDto(message)
	}
	return chatMessages
}
