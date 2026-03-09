package service

import (
	userDao "github.com/ACaiCat/tiktok-go/pkg/db/user"
)

type UserService struct {
	dao *userDao.UserDao
}

func NewUserService() *UserService {
	return &UserService{
		dao: userDao.NewUserDao(),
	}
}
