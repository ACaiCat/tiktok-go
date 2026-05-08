package service

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/user"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/ACaiCat/tiktok-go/pkg/jwt"
	totp "github.com/ACaiCat/tiktok-go/pkg/totp"
)

func (s *UserService) UserLogin(req *user.LoginReq) (*model.User, string, string, error) {
	var err error

	usr, err := s.dao.GetByUsername(s.ctx, req.Username)
	if err != nil {
		return nil, "", "", errors.WithMessagef(err, "service.UserLogin: db.GetByUsername failed, username=%q", req.Username)
	}

	if usr == nil {
		return nil, "", "", errno.PasswordIsNotVerified
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(req.Password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, "", "", errno.PasswordIsNotVerified
		}
		return nil, "", "", errors.WithMessagef(err, "service.UserLogin: bcrypt.CompareHashAndPassword failed, username=%q", req.Username)
	}

	if usr.TotpSecret != nil {
		if req.Code == nil || *req.Code == "" {
			return nil, "", "", errno.MFAMissingErr
		}

		ok, err := totp.ValidateCode(*usr.TotpSecret, *req.Code)

		if err != nil {
			return nil, "", "", errors.WithMessagef(err, "service.UserLogin: totp.ValidateCode failed, userID=%d", usr.ID)
		}

		if !ok {
			return nil, "", "", errno.MFACodeInvalidErr
		}
	}

	accessToken, err := jwt.CreateToken(constants.TypeAccessToken, usr.ID)

	if err != nil {
		return nil, "", "", errors.WithMessagef(err, "service.UserLogin: create access token failed, userID=%d", usr.ID)
	}

	refreshToken, err := jwt.CreateToken(constants.TypeRefreshToken, usr.ID)

	if err != nil {
		return nil, "", "", errors.WithMessagef(err, "service.UserLogin: create refresh token failed, userID=%d", usr.ID)
	}

	return UserDaoToDto(usr), accessToken, refreshToken, nil
}
