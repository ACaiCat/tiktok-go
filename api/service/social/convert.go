package service

import (
	"strconv"

	"github.com/ACaiCat/tiktok-go/api/model/model"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func UserToSocialUser(user *modelDao.User) *model.SocialUser {
	return &model.SocialUser{
		ID:        strconv.FormatInt(user.ID, 10),
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
	}
}

func UsersToSocialUsers(users []*modelDao.User) []*model.SocialUser {
	socialUsers := make([]*model.SocialUser, len(users))

	for i, user := range users {
		socialUsers[i] = UserToSocialUser(user)
	}

	return socialUsers
}
