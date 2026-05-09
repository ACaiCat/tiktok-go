package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/biz/model/ws"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func (s *ChatService) HandleChatMessage(userID int64, chatMessage *ws.ChatMessage) error {
	if !s.ensureFriend(userID, chatMessage.ReceiverID) {
		return nil
	}

	if _, err := s.sendMessage(userID, ws.MessageTypeChat, chatMessage); err != nil {
		return errors.Wrapf(err, "sendMessageToUser failed, userID=%d", userID)
	}

	receiverOnline, err := s.sendMessage(chatMessage.ReceiverID, ws.MessageTypeChat, chatMessage)
	if err != nil {
		return errors.Wrapf(err, "sendMessageToUser failed, userID=%d", chatMessage.ReceiverID)
	}
	if !s.saveChatMessage(userID, chatMessage.ReceiverID, chatMessage.Content, receiverOnline, chatMessage.IsAI) {
		return nil
	}

	go s.replyWithAI(userID, chatMessage.ReceiverID)
	return nil
}

func (s *ChatService) ensureFriend(userID int64, receiverID int64) bool {
	isFriend, err := s.followerDao.IsExistFriend(s.ctx, userID, receiverID)
	if err != nil {
		s.SendErr(userID, errno.ServiceErr)
		return false
	}
	if !isFriend {
		s.SendErr(userID, errno.ChatNotFriendErr)
		return false
	}

	return true
}

func (s *ChatService) saveChatMessage(senderID int64, receiverID int64, content string, isRead bool, isAI bool) bool {
	if err := s.chatDao.AddMessage(s.ctx, senderID, receiverID, content, isRead, isAI); err != nil {
		s.SendErr(senderID, errno.ServiceErr)
		return false
	}

	go s.clearConversationHistoryCache(senderID, receiverID)
	go s.clearUnreadMessagesCache(receiverID, senderID)

	return true
}

func (s *ChatService) clearConversationHistoryCache(senderID int64, receiverID int64) {
	if err := s.cache.ClearChatHistory(context.Background(), senderID, receiverID); err != nil {
		hlog.Errorf("ClearChatHistory failed, sender=%d receiver=%d: %v", senderID, receiverID, err)
	}
}
