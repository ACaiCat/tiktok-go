package service

import (
	"github.com/ACaiCat/tiktok-go/api/model/ws"
	"github.com/ACaiCat/tiktok-go/pkg/db"
	chatDao "github.com/ACaiCat/tiktok-go/pkg/db/chat"
	followerDao "github.com/ACaiCat/tiktok-go/pkg/db/follower"
	userDao "github.com/ACaiCat/tiktok-go/pkg/db/user"
)

type ChatService struct {
	userDao     *userDao.UserDao
	followerDao *followerDao.FollowerDao
	chatDao     *chatDao.ChatDao
	manager     *ws.OnlineUserManager
}

func NewChatService(manger *ws.OnlineUserManager) *ChatService {
	return &ChatService{
		userDao:     userDao.NewUserDao(db.DB),
		followerDao: followerDao.NewFollowerDao(db.DB),
		chatDao:     chatDao.NewChatDao(db.DB),
		manager:     manger,
	}
}
