package ws

import (
	"github.com/gorilla/websocket"
)

type OnlineUserManager struct {
	onlineUsers map[int64]OnlineUser
}

func NewOnlineUserManager() *OnlineUserManager {
	return &OnlineUserManager{
		onlineUsers: make(map[int64]OnlineUser),
	}
}

func (m *OnlineUserManager) IsUserOnline(userID int64) bool {
	_, online := m.onlineUsers[userID]
	return online
}

func (m *OnlineUserManager) GetOnlineUser(userID int64) (OnlineUser, bool) {
	user, online := m.onlineUsers[userID]
	return user, online
}

func (m *OnlineUserManager) AddOnlineUser(userID int64, conn *websocket.Conn) {
	m.onlineUsers[userID] = OnlineUser{
		UserID: userID,
		Conn:   conn,
	}
}

func (m *OnlineUserManager) RemoveOnlineUser(userID int64) {
	delete(m.onlineUsers, userID)
}
