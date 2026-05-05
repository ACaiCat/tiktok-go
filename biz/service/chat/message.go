package service

import (
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
		if !s.unmarshalBody(userID, message.Body, &chatMessage, "聊天消息格式错误：") {
			return
		}
		s.handleChatMessage(userID, &chatMessage)

	case ws.MessageTypeUnread:
		var unreadRequest ws.UnreadRequest
		if !s.unmarshalBody(userID, message.Body, &unreadRequest, "未读消息请求格式错误：") {
			return
		}
		s.handleUnreadMessage(userID, &unreadRequest)

	case ws.MessageTypeHistory:
		var historyRequest ws.HistoryRequest
		if !s.unmarshalBody(userID, message.Body, &historyRequest, "历史消息请求格式错误：") {
			return
		}
		s.handleHistoryMessage(userID, &historyRequest)

	default:
		s.SendErr(userID, errno.ChatMsgTypeErr)
	}
}

func (s *ChatService) unmarshalBody(userID int64, body []byte, target any, errPrefix string) bool {
	if err := json.Unmarshal(body, target); err != nil {
		s.SendErr(userID, errno.ChatMsgParseErr.WithMessage(errPrefix+err.Error()))
		return false
	}

	return true
}
