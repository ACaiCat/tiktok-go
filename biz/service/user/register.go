package service

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/ACaiCat/tiktok-go/biz/model/user"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func (s *UserService) UserRegister(req *user.RegisterReq) error {
	var err error

	if len(req.Username) < constants.MinUsernameLength {
		return errno.UsernameTooShortErr.WithMessage(fmt.Sprintf("username must be at least %d characters", constants.MinUsernameLength))
	}

	if len(req.Username) > constants.MaxUsernameLength {
		return errno.UsernameTooLongErr.WithMessage(fmt.Sprintf("username must be at most %d characters", constants.MaxUsernameLength))
	}

	if len(req.Password) < constants.MinPasswordLength {
		return errno.PasswordTooShortErr.WithMessage(fmt.Sprintf("password must be at least %d characters", constants.MinPasswordLength))
	}

	if len(req.Password) > constants.MaxPasswordLength {
		return errno.PasswordTooLongErr.WithMessage(fmt.Sprintf("password must be at most %d characters", constants.MaxPasswordLength))
	}

	usr, err := s.dao.GetByUsername(s.ctx, req.Username)
	if err != nil {
		return errors.WithMessagef(err, "service.UserRegister: db.GetByUsername failed, username=%q", req.Username)
	}

	if usr != nil {
		return errno.UserAlreadyExistErr
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return errors.WithMessagef(err, "service.UserRegister: bcrypt.GenerateFromPassword failed, username=%q", req.Username)
	}

	_, err = s.dao.CreateUser(s.ctx, req.Username, string(hashPassword))
	if err != nil {
		return errors.WithMessagef(err, "service.UserRegister: db.CreateUser failed, username=%q", req.Username)
	}

	return nil
}
