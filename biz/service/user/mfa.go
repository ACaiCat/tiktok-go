package service

import (
	"encoding/base64"
	"log"

	"github.com/skip2/go-qrcode"

	"github.com/ACaiCat/tiktok-go/biz/model/user"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/ACaiCat/tiktok-go/pkg/totp"
)

func (s *UserService) GetMFA(userID int64) (string, string, error) {
	var err error

	usr, err := s.dao.GetByID(userID)
	if err != nil {
		return "", "", errno.ServiceErr
	}

	if usr == nil {
		return "", "", errno.UserIsNotExistErr
	}

	key, err := totp.CreateKey(usr.Username)
	if err != nil {
		return "", "", errno.ServiceErr
	}

	rawQrcode, err := qrcode.Encode(key.String(), qrcode.Low, constants.TotpQRCodeSize)
	if err != nil {
		log.Println("failed to generate QR code for user", usr.Username, ":", err)
		return "", "", errno.ServiceErr
	}

	base64Qrcode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(rawQrcode)

	return key.Secret(), base64Qrcode, nil
}

func (s *UserService) BindMFA(req *user.BindMFAReq, userID int64) error {
	var err error

	ok, err := totp.ValidateCode(req.Secret, req.Code)

	if err != nil {
		return errno.ServiceErr
	}

	if !ok {
		return errno.MFACodeInvalidErr
	}

	err = s.dao.UpdateUserMFA(userID, req.Secret)
	if err != nil {
		return errno.ServiceErr
	}

	return nil
}
