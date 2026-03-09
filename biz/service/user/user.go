package service

import (
	"github.com/ACaiCat/tiktok-go/pkg/cache"
	userCache "github.com/ACaiCat/tiktok-go/pkg/cache/user"
	"github.com/ACaiCat/tiktok-go/pkg/db"
	userDao "github.com/ACaiCat/tiktok-go/pkg/db/user"
)

type UserService struct {
	dao   *userDao.UserDao
	cache *userCache.UserCache
}

func NewUserService() *UserService {
	return &UserService{
		dao:   userDao.NewUserDao(db.DB),
		cache: userCache.NewUserCache(cache.Cache),
	}
}
