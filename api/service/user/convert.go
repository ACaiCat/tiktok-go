package service

import (
	"strconv"

	dto "github.com/ACaiCat/tiktok-go/api/model/model"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func UserDaoToDto(user *model.User) *dto.User {
	return &dto.User{
		ID:        strconv.FormatInt(user.ID, 10),
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
		CreatedAt: strconv.FormatInt(user.CreatedAt.UnixMilli(), 10),
	}
}
