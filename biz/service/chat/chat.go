package service

import (
	"context"

	"github.com/ACaiCat/tiktok-go/biz/model/ws"
	"github.com/ACaiCat/tiktok-go/pkg/cache"
	chatCache "github.com/ACaiCat/tiktok-go/pkg/cache/chat"
	"github.com/ACaiCat/tiktok-go/pkg/db"
	chatDao "github.com/ACaiCat/tiktok-go/pkg/db/chat"
	followerDao "github.com/ACaiCat/tiktok-go/pkg/db/follower"
	userDao "github.com/ACaiCat/tiktok-go/pkg/db/user"
)

type ChatService struct {
	userDao     *userDao.UserDao
	followerDao *followerDao.FollowerDao
	chatDao     *chatDao.ChatDao
	cache       *chatCache.ChatCache
	manager     *ws.OnlineUserManager
	ctx         context.Context
}

func NewChatService(ctx context.Context, manger *ws.OnlineUserManager) *ChatService {
	return &ChatService{
		userDao:     userDao.NewUserDao(db.DB),
		followerDao: followerDao.NewFollowerDao(db.DB),
		chatDao:     chatDao.NewChatDao(db.DB),
		cache:       chatCache.NewChatCache(cache.Cache),
		manager:     manger,
		ctx:         ctx,
	}
}
