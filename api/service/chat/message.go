package service

import (
	"log"

	"github.com/cloudwego/hertz/pkg/common/json"

	"github.com/ACaiCat/tiktok-go/api/model/ws"
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
		chatMessage := ws.ChatMessage{}
		if err := json.Unmarshal(message.Body, &chatMessage); err != nil {
			s.SendErr(userID, errno.ChatMsgParseErr.WithMessage("聊天消息格式错误："+err.Error()))
			return
		}

		isFriend, err := s.followerDao.IsExistFriend(userID, chatMessage.ReceiverID)
		if err != nil {
			s.SendErr(userID, errno.ServiceErr.WithMessage("查询好友关系失败"))
			return
		}
		if !isFriend {
			s.SendErr(userID, errno.ChatNotFriendErr)
			return
		}

		if sender, online := s.manager.GetOnlineUser(userID); online {
			if err := sender.SendMessage(ws.MessageTypeChat, &chatMessage); err != nil {
				log.Println("failed to echo message to sender:", err)
			}
		}

		receiver, receiverOnline := s.manager.GetOnlineUser(chatMessage.ReceiverID)
		if receiverOnline {
			if err := receiver.SendMessage(ws.MessageTypeChat, &chatMessage); err != nil {
				log.Println("failed to forward message to receiver:", err)
			}
			if err := s.chatDao.AddMessage(userID, chatMessage.ReceiverID, chatMessage.Content, true); err != nil {
				s.SendErr(userID, errno.ServiceErr.WithMessage("消息保存失败"))
				return
			}
		} else {
			if err := s.chatDao.AddMessage(userID, chatMessage.ReceiverID, chatMessage.Content, false); err != nil {
				s.SendErr(userID, errno.ServiceErr.WithMessage("消息保存失败"))
				return
			}
		}

	case ws.MessageTypeUnread:
		unreadRequest := ws.UnreadRequest{}
		if err := json.Unmarshal(message.Body, &unreadRequest); err != nil {
			s.SendErr(userID, errno.ChatMsgParseErr.WithMessage("未读消息请求格式错误："+err.Error()))
			return
		}

		unreadMessages, err := s.chatDao.GetUnreadMessages(userID, unreadRequest.Sender)
		if err != nil {
			s.SendErr(userID, errno.ServiceErr.WithMessage("获取未读消息失败："+err.Error()))
			return
		}

		user, online := s.manager.GetOnlineUser(userID)
		if online {
			if err := user.SendMessage(ws.MessageTypeUnread, &ws.HistoryMessage{
				Messages: MessagesDaoToDto(unreadMessages),
			}); err != nil {
				log.Println("failed to send unread messages to user:", err)
				return
			}
			if err := s.chatDao.MarkMessagesAsRead(userID, unreadRequest.Sender); err != nil {
				s.SendErr(userID, errno.ServiceErr.WithMessage("标记消息已读失败"))
				return
			}
		}

	case ws.MessageTypeHistory:
		historyRequest := ws.HistoryRequest{}
		if err := json.Unmarshal(message.Body, &historyRequest); err != nil {
			s.SendErr(userID, errno.ChatMsgParseErr.WithMessage("历史消息请求格式错误："+err.Error()))
			return
		}

		historyMessages, err := s.chatDao.GetChatHistory(userID, historyRequest.Sender, historyRequest.PageSize, historyRequest.Page)
		if err != nil {
			s.SendErr(userID, errno.ServiceErr.WithMessage("获取历史消息失败"))
			return
		}

		user, online := s.manager.GetOnlineUser(userID)
		if online {
			if err := user.SendMessage(ws.MessageTypeHistory, &ws.HistoryMessage{
				Messages: MessagesDaoToDto(historyMessages),
			}); err != nil {
				log.Println("failed to send chat history to user:", err)
			}
		}

	default:
		s.SendErr(userID, errno.ChatMsgTypeErr)
	}
}
