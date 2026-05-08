package service

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/json"

	"github.com/ACaiCat/tiktok-go/biz/model/ws"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func (s *ChatService) SendErr(userID int64, err errno.ErrNo) {
	if u, online := s.manager.GetOnlineUser(userID); online {
		u.SendError(int(err.ErrCode), err.ErrMsg)
	}
}

func (s *ChatService) HandleMessage(userID int64, messageText string) {
	var message ws.Message
	if err := json.Unmarshal([]byte(messageText), &message); err != nil {
		s.SendErr(userID, errno.ChatMsgParseErr.WithMessage("消息格式错误："+err.Error()))
		return
	}

	switch message.Type {
	case ws.MessageTypeChat:
		var chatMessage ws.ChatMessage
		if err := json.Unmarshal(message.Body, &chatMessage); err != nil {
			s.SendErr(userID, errno.ChatMsgParseErr)
			return
		}
		if err := s.handleChatMessage(userID, &chatMessage); err != nil {
			s.SendErr(userID, errno.ServiceErr)
			hlog.CtxErrorf(s.ctx, "handleChatMessage: %v", err)
		}

	case ws.MessageTypeUnread:
		var unreadRequest ws.UnreadRequest
		if err := json.Unmarshal(message.Body, &unreadRequest); err != nil {
			s.SendErr(userID, errno.ChatMsgParseErr)
			return
		}
		if err := s.handleUnreadMessage(userID, &unreadRequest); err != nil {
			s.SendErr(userID, errno.ServiceErr)
			hlog.CtxErrorf(s.ctx, "handleUnreadMessage: %v", err)
		}

	case ws.MessageTypeHistory:
		var historyRequest ws.HistoryRequest
		if err := json.Unmarshal(message.Body, &historyRequest); err != nil {
			s.SendErr(userID, errno.ChatMsgParseErr)
			return
		}
		if err := s.handleHistoryMessage(userID, &historyRequest); err != nil {
			s.SendErr(userID, errno.ServiceErr)
			hlog.CtxErrorf(s.ctx, "handleHistoryMessage: %v", err)
		}

	default:
		s.SendErr(userID, errno.ChatMsgTypeErr)
	}
}
