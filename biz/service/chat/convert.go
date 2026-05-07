package service

import (
	"github.com/ACaiCat/tiktok-go/biz/model/ws"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (s *ChatService) MessageDaoToDto(message *model.ChatMessage) *ws.ChatMessage {
	content := message.Content

	return &ws.ChatMessage{
		SenderID:   message.SenderID,
		ReceiverID: message.ReceiverID,
		Content:    content,
		Timestamp:  message.CreatedAt.Unix(),
	}
}

func (s *ChatService) MessagesDaoToDto(messages []*model.ChatMessage) []*ws.ChatMessage {
	chatMessages := make([]*ws.ChatMessage, len(messages))
	for i, message := range messages {
		chatMessages[i] = s.MessageDaoToDto(message)
	}
	return chatMessages
}
