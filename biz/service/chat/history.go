package service

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ACaiCat/tiktok-go/biz/model/ws"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (s *ChatService) HandleHistoryMessage(userID int64, historyRequest *ws.HistoryRequest) error {
	historyMessages, err := s.getChatHistory(userID, historyRequest.Sender, historyRequest.PageSize, historyRequest.Page)
	if err != nil {
		return errors.Wrapf(err, "getChatHistory failed, userID=%d, sender=%d", userID, historyRequest.Sender)
	}

	if _, err := s.sendMessage(userID, ws.MessageTypeHistory, &ws.HistoryMessage{
		Messages: s.MessagesDaoToDto(historyMessages),
	}); err != nil {
		return errors.Wrapf(err, "sendMessageToUser failed, userID=%d", userID)
	}
	return nil
}

func (s *ChatService) getChatHistory(userID int64, otherUserID int64, pageSize int, pageNum int) ([]*model.ChatMessage, error) {
	messages, err := s.cache.GetChatHistory(s.ctx, userID, otherUserID, pageSize, pageNum)
	if err == nil {
		return messages, nil
	}
	if !errors.Is(err, redis.Nil) {
		hlog.CtxErrorf(s.ctx, "get chat history cache err user=%d other=%d: %v", userID, otherUserID, err)
	}

	messages, err = s.chatDao.GetChatHistory(s.ctx, userID, otherUserID, pageSize, pageNum)
	if err != nil {
		return nil, err
	}

	if err := s.cache.SetChatHistory(s.ctx, userID, otherUserID, pageSize, pageNum, messages); err != nil {
		hlog.CtxErrorf(s.ctx, "set chat history cache err user=%d other=%d: %v", userID, otherUserID, err)
	}

	return messages, nil
}
