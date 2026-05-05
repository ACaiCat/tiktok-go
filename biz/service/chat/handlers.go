package service

import (
	"log"

	"github.com/ACaiCat/tiktok-go/biz/model/ws"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func (s *ChatService) handleChatMessage(userID int64, chatMessage *ws.ChatMessage) {
	if !s.ensureFriend(userID, chatMessage.ReceiverID) {
		return
	}

	s.sendMessageToUser(userID, ws.MessageTypeChat, chatMessage, "failed to echo message to sender")

	receiverOnline := s.sendMessageToUser(chatMessage.ReceiverID, ws.MessageTypeChat, chatMessage, "failed to forward message to receiver")
	if !s.saveChatMessage(userID, chatMessage.ReceiverID, chatMessage.Content, receiverOnline, chatMessage.IsAI) {
		return
	}

	go s.replyWithAI(userID, chatMessage.ReceiverID)
}

func (s *ChatService) handleUnreadMessage(userID int64, unreadRequest *ws.UnreadRequest) {
	unreadMessages, err := s.chatDao.GetUnreadMessages(s.ctx, userID, unreadRequest.Sender)
	if err != nil {
		s.SendErr(userID, errno.ServiceErr.WithMessage("获取未读消息失败："+err.Error()))
		return
	}

	if !s.sendMessageToUser(userID, ws.MessageTypeUnread, &ws.HistoryMessage{
		Messages: s.MessagesDaoToDto(unreadMessages),
	}, "failed to send unread messages to user") {
		return
	}

	if err := s.chatDao.MarkMessagesAsRead(s.ctx, userID, unreadRequest.Sender); err != nil {
		s.SendErr(userID, errno.ServiceErr.WithMessage("标记消息已读失败"))
	}
}

func (s *ChatService) handleHistoryMessage(userID int64, historyRequest *ws.HistoryRequest) {
	historyMessages, err := s.chatDao.GetChatHistory(s.ctx, userID, historyRequest.Sender, historyRequest.PageSize, historyRequest.Page)
	if err != nil {
		s.SendErr(userID, errno.ServiceErr.WithMessage("获取历史消息失败"))
		return
	}

	s.sendMessageToUser(userID, ws.MessageTypeHistory, &ws.HistoryMessage{
		Messages: s.MessagesDaoToDto(historyMessages),
	}, "failed to send chat history to user")
}

func (s *ChatService) ensureFriend(userID int64, receiverID int64) bool {
	isFriend, err := s.followerDao.IsExistFriend(s.ctx, userID, receiverID)
	if err != nil {
		s.SendErr(userID, errno.ServiceErr.WithMessage("查询好友关系失败"))
		return false
	}
	if !isFriend {
		s.SendErr(userID, errno.ChatNotFriendErr)
		return false
	}

	return true
}

func (s *ChatService) sendMessageToUser(userID int64, msgType int, body any, logPrefix string) bool {
	user, online := s.manager.GetOnlineUser(userID)
	if !online {
		return false
	}

	if err := user.SendMessage(msgType, body); err != nil {
		log.Println(logPrefix+":", err)
	}

	return true
}

func (s *ChatService) saveChatMessage(senderID int64, receiverID int64, content string, isRead bool, isAI bool) bool {
	if err := s.chatDao.AddMessage(s.ctx, senderID, receiverID, content, isRead, isAI); err != nil {
		s.SendErr(senderID, errno.ServiceErr.WithMessage("消息保存失败"))
		return false
	}

	return true
}
