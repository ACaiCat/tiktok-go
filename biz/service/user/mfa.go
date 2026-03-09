package service

import (
	"log"

	"github.com/ACaiCat/tiktok-go/biz/model/tiktok-go/user"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	totp "github.com/ACaiCat/tiktok-go/pkg/totp"
)

func (s *UserService) GetMFA(req *user.MFAQRCodeReq, userID int64) (string, error) {
	var err error

	usr, err := s.dao.GetByID(userID)
	if err != nil {
		return "", errno.ServiceErr
	}

	if usr == nil {
		return "", errno.UserIsNotExistErr
	}

	secret, err := totp.CreateSecret(usr.Username)
	if err != nil {
		log.Println("failed to create secret for user", usr.Username, ":", err)
		return "", errno.ServiceErr
	}

	return secret, nil
}

func (s *UserService) BindMFA(req *user.BindMFAReq, userID int64) error {
	var err error

	if !totp.ValidateCode(req.Secret, req.Code) {
		return errno.MFACodeInvalidErr
	}

	err = s.dao.UpdateUserMFA(userID, req.Secret)
	if err != nil {
		return errno.ServiceErr
	}

	return nil

}
