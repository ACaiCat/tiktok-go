package service

import (
	"errors"
	"log"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/tiktok-go/user"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/ACaiCat/tiktok-go/pkg/jwt"
	totp "github.com/ACaiCat/tiktok-go/pkg/totp"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) UserLogin(req *user.LoginReq) (*model.User, string, string, error) {
	var err error

	usr, err := s.dao.GetByUsername(req.Username)
	if err != nil {
		return nil, "", "", errno.ServiceErr
	}

	if usr == nil {
		return nil, "", "", errno.PasswordIsNotVerified
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(req.Password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, "", "", errno.PasswordIsNotVerified
		}
		log.Println("Error comparing password hash:", err)
		return nil, "", "", errno.ServiceErr
	}

	if usr.TotpSecret != "" {
		if req.Code == "" {
			return nil, "", "", errno.MFAMissingErr
		}

		ok, err := totp.ValidateCode(usr.TotpSecret, req.Code)

		if err != nil {
			return nil, "", "", errno.ServiceErr
		}

		if !ok {
			return nil, "", "", errno.MFACodeInvalidErr
		}
	}

	accessToken, err := jwt.CreateToken(constants.TypeAccessToken, usr.ID)

	if err != nil {
		return nil, "", "", errno.ServiceErr
	}

	refreshToken, err := jwt.CreateToken(constants.TypeRefreshToken, usr.ID)

	if err != nil {
		return nil, "", "", errno.ServiceErr
	}

	return UserDaoToDTO(usr), accessToken, refreshToken, nil

}
