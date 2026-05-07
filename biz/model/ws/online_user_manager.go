package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type OnlineUserManager struct {
	mu          sync.RWMutex
	onlineUsers map[int64]OnlineUser
}

func NewOnlineUserManager() *OnlineUserManager {
	return &OnlineUserManager{
		onlineUsers: make(map[int64]OnlineUser),
	}
}

func (m *OnlineUserManager) IsUserOnline(userID int64) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, online := m.onlineUsers[userID]
	return online
}

func (m *OnlineUserManager) GetOnlineUser(userID int64) (OnlineUser, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, online := m.onlineUsers[userID]
	return user, online
}

func (m *OnlineUserManager) AddOnlineUser(userID int64, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.onlineUsers[userID] = OnlineUser{
		UserID: userID,
		Conn:   conn,
	}
}

func (m *OnlineUserManager) RemoveOnlineUser(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.onlineUsers, userID)
}
