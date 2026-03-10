package service

import (
	"strconv"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)
import dto "github.com/ACaiCat/tiktok-go/biz/model/model"

func UserDaoToDTO(user *model.User) *dto.User {
	return &dto.User{
		ID:        strconv.FormatInt(user.ID, 10),
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
		CreatedAt: strconv.FormatInt(user.CreatedAt.UnixMilli(), 10),
	}
}
