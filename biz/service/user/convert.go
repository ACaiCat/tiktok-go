package service

import (
	"strconv"

	dto "github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func UserDaoToDto(user *model.User) *dto.User {
	avatarURL := ""
	if user.AvatarURL != nil {
		avatarURL = *user.AvatarURL
	}
	return &dto.User{
		ID:        strconv.FormatInt(user.ID, 10),
		Username:  user.Username,
		AvatarURL: avatarURL,
		CreatedAt: strconv.FormatInt(user.CreatedAt.UnixMilli(), 10),
	}
}
