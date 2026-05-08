package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ACaiCat/tiktok-go/biz/model/ws"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (s *ChatService) handleUnreadMessage(userID int64, unreadRequest *ws.UnreadRequest) error {
	unreadMessages, err := s.getUnreadMessages(userID, unreadRequest.Sender)
	if err != nil {
		return errors.Wrapf(err, "getUnreadMessages failed, userID=%d, sender=%d", userID, unreadRequest.Sender)
	}

	online, err := s.sendMessageToUser(userID, ws.MessageTypeUnread, &ws.UnreadMessage{
		Messages: s.MessagesDaoToDto(unreadMessages),
	})
	if err != nil {
		return errors.Wrapf(err, "sendMessageToUser failed, userID=%d", userID)
	}
	if !online {
		return nil
	}

	if err := s.chatDao.MarkMessagesAsRead(s.ctx, userID, unreadRequest.Sender); err != nil {
		return errors.Wrapf(err, "MarkMessagesAsRead failed, userID=%d, sender=%d", userID, unreadRequest.Sender)
	}

	go s.clearUnreadMessagesCache(userID, unreadRequest.Sender)
	return nil
}

func (s *ChatService) getUnreadMessages(userID int64, senderID int64) ([]*model.ChatMessage, error) {
	messages, err := s.cache.GetUnreadMessages(s.ctx, userID, senderID)
	if err == nil {
		return messages, nil
	}
	if !errors.Is(err, redis.Nil) {
		hlog.CtxErrorf(s.ctx, "get unread messages cache err user=%d sender=%d: %v", userID, senderID, err)
	}

	messages, err = s.chatDao.GetUnreadMessages(s.ctx, userID, senderID)
	if err != nil {
		return nil, err
	}

	if err := s.cache.SetUnreadMessages(s.ctx, userID, senderID, messages); err != nil {
		hlog.CtxErrorf(s.ctx, "set unread messages cache err user=%d sender=%d: %v", userID, senderID, err)
	}

	return messages, nil
}

func (s *ChatService) clearUnreadMessagesCache(userID int64, senderID int64) {
	if err := s.cache.ClearUnreadMessages(context.Background(), userID, senderID); err != nil {
		hlog.Errorf("ClearUnreadMessages failed, user=%d sender=%d: %v", userID, senderID, err)
	}
}
