package service

import (
	"encoding/base64"

	"github.com/pkg/errors"
	"github.com/skip2/go-qrcode"

	"github.com/ACaiCat/tiktok-go/biz/model/user"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/crypt"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/ACaiCat/tiktok-go/pkg/totp"
)

func (s *UserService) GetMFA(userID int64) (string, string, error) {
	var err error

	usr, err := s.dao.GetByID(s.ctx, userID)
	if err != nil {
		return "", "", err
	}

	if usr == nil {
		return "", "", errno.UserIsNotExistErr
	}

	key, err := totp.CreateKey(usr.Username)
	if err != nil {
		return "", "", errors.WithMessagef(err, "service.GetMFA: totp.CreateKey failed, userID=%d", userID)
	}

	rawQrcode, err := qrcode.Encode(key.String(), qrcode.Low, constants.TotpQRCodeSize)
	if err != nil {
		return "", "", errors.WithMessagef(err, "service.GetMFA: qrcode.Encode failed, userID=%d", userID)
	}

	base64Qrcode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(rawQrcode)

	return key.Secret(), base64Qrcode, nil
}

func (s *UserService) BindMFA(req *user.BindMFAReq, userID int64) error {
	var err error

	ok, err := totp.ValidateCode(req.Secret, req.Code)

	if err != nil {
		return errors.WithMessagef(err, "service.BindMFA: totp.ValidateCode failed, userID=%d", userID)
	}

	if !ok {
		return errno.MFACodeInvalidErr
	}

	secret, err := crypt.Encrypt(req.Secret)
	if err != nil {
		return errors.Wrapf(err, "service.BindMFA: crypt.Encrypt failed, userID=%d", userID)
	}

	err = s.dao.UpdateUserMFA(s.ctx, userID, secret)
	if err != nil {
		return errors.WithMessagef(err, "service.BindMFA: db.UpdateUserMFA failed, userID=%d", userID)
	}

	return nil
}
