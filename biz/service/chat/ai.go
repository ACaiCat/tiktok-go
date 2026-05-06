package service

import (
	"log"
	"slices"
	"time"

	"github.com/ACaiCat/tiktok-go/biz/model/ws"
)

func (s *ChatService) replyWithAI(userID int64, receiverID int64) {
	messages, err := s.chatDao.GetChatHistory(s.ctx, userID, receiverID, 10, 0)
	if err != nil {
		log.Println("failed to get message history:", err)
		return
	}

	// 反转使最新的消息在上下文底部
	slices.Reverse(messages)
	history := BuildAIHistory(messages, userID, receiverID)

	reply, content, err := s.ChatAI(s.ctx, history)
	if err != nil {
		log.Println("failed to AI message:", err)
		return
	}
	if !reply {
		return
	}

	now := time.Now().UnixMilli()
	receiverMessage := &ws.ChatMessage{
		SenderID:   userID,
		ReceiverID: receiverID,
		IsAI:       true,
		Content:    content,
		Timestamp:  now,
	}
	receiverOnline := s.sendMessageToUser(receiverID, ws.MessageTypeChat, receiverMessage, "failed to forward message to receiver")
	if !s.saveChatMessage(userID, receiverID, content, receiverOnline, true) {
		return
	}

	senderMessage := &ws.ChatMessage{
		SenderID:   receiverID,
		ReceiverID: userID,
		IsAI:       true,
		Content:    content,
		Timestamp:  now,
	}
	senderOnline := s.sendMessageToUser(userID, ws.MessageTypeChat, senderMessage, "failed to forward message to sender")
	_ = s.saveChatMessage(receiverID, userID, content, senderOnline, true)
}
