package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/ACaiCat/tiktok-go/biz/model/user"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func (s *UserService) UserRegister(req *user.RegisterReq) error {
	var err error

	if len(req.Username) < constants.MinUsernameLength {
		return errno.UsernameTooShortErr.WithMessage(fmt.Sprintf("用户名长度必须至少为%d位", constants.MinUsernameLength))
	}

	if len(req.Username) > constants.MaxUsernameLength {
		return errno.UsernameTooLongErr.WithMessage(fmt.Sprintf("用户名长度必须不超过%d位", constants.MaxUsernameLength))
	}

	if len(req.Password) < constants.MinPasswordLength {
		return errno.PasswordTooShortErr.WithMessage(fmt.Sprintf("密码长度必须至少为%d位", constants.MinPasswordLength))
	}

	if len(req.Password) > constants.MaxPasswordLength {
		return errno.PasswordTooLongErr.WithMessage(fmt.Sprintf("密码长度必须不超过%d位", constants.MaxPasswordLength))
	}

	usr, err := s.dao.GetByUsername(req.Username)
	if err != nil {
		return errno.ServiceErr
	}

	if usr != nil {
		return errno.UserAlreadyExistErr
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return errno.ServiceErr
	}

	_, err = s.dao.CreateUser(req.Username, string(hashPassword))
	if err != nil {
		return errno.ServiceErr
	}

	return nil
}
