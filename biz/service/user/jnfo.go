package service

import (
	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func (s *UserService) GetUserInfo(userID int64) (*model.User, error) {
	usr, err := s.dao.GetByID(userID)
	if err != nil {
		return nil, errno.ServiceErr
	}

	if usr == nil {
		return nil, errno.UserIsNotExistErr
	}

	return UserDaoToDto(usr), nil
}
