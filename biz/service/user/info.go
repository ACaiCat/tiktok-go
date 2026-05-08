package service

import (
	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func (s *UserService) GetUserInfo(userID int64) (*model.User, error) {
	usr, err := s.dao.GetByID(s.ctx, userID)
	if err != nil {
		return nil, errors.WithMessagef(err, "service.GetUserInfo: db.GetByID failed, userID=%d", userID)
	}

	if usr == nil {
		return nil, errno.UserIsNotExistErr
	}

	return UserDaoToDto(usr), nil
}
