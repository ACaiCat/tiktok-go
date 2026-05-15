package service

import (
	"context"

	"github.com/ACaiCat/tiktok-go/pkg/cache"
	userCache "github.com/ACaiCat/tiktok-go/pkg/cache/user"
	"github.com/ACaiCat/tiktok-go/pkg/db"
	followerDao "github.com/ACaiCat/tiktok-go/pkg/db/follower"
	userDao "github.com/ACaiCat/tiktok-go/pkg/db/user"
)

type SocialService struct {
	userDao     *userDao.UserDao
	followerDao *followerDao.FollowerDao
	userCache   *userCache.UserCache
	ctx         context.Context
}

func NewSocialService(ctx context.Context) *SocialService {
	return &SocialService{
		userDao:     userDao.NewUserDao(db.DB),
		followerDao: followerDao.NewFollowerDao(db.DB),
		userCache:   userCache.NewUserCache(cache.Cache),
		ctx:         ctx,
	}
}
