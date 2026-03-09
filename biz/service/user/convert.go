package service

import (
	"strconv"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)
import dto "github.com/ACaiCat/tiktok-go/biz/model/model"

func UserDaoToDTO(user *model.User) *dto.User {
	return &dto.User{
		Id:        strconv.FormatInt(user.ID, 10),
		Username:  user.Username,
		AvatarUrl: user.AvatarURL,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		DeletedAt: user.DeletedAt.Time.String(),
	}
}
