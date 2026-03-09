package service

import (
	"fmt"

	"github.com/ACaiCat/tiktok-go/biz/model/tiktok-go/user"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) UserRegister(req *user.RegisterReq) error {
	var err error

	usr, err := s.dao.GetByUsername(req.Username)
	if err != nil {
		return errno.ServiceErr
	}

	if usr != nil {
		return errno.UserAlreadyExistErr
	}

	if len(req.Password) < constants.MinPasswordLength {
		return errno.PasswordTooShortErr.WithMessage(fmt.Sprintf("密码长度必须至少为%d位", constants.MinPasswordLength))
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	_, err = s.dao.CreateUser(req.Username, string(hashPassword))
	if err != nil {
		return errno.ServiceErr
	}

	return nil

}
