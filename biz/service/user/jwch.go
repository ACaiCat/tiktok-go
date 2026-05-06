package service

import (
	"log"

	"github.com/ACaiCat/tiktok-go/biz/model/user"
	"github.com/ACaiCat/tiktok-go/pkg/crypt"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/west2-online/jwch"
)

func (s *UserService) BindJwch(req *user.BindJwchReq, userID int64) error {
	var err error

	stu := jwch.NewStudent()
	stu.ID = req.JwchID
	stu.Password = req.JwchPassword

	err = stu.Login()
	if err != nil {
		return errno.JwchLoginFailedErr.WithError(err)
	}

	password, err := crypt.Encrypt(req.JwchPassword)

	if err != nil {
		log.Printf("password encrypt err %v", err)
		return errno.ServiceErr
	}

	err = s.dao.UpdateUserJwch(s.ctx, userID, req.JwchID, password)
	if err != nil {
		return errno.ServiceErr
	}

	return nil
}

func (s *UserService) GetJwchIdentifierAndCookies(userID int64) (string, string, error) {
	var err error

	log.Println(userID)
	usr, err := s.dao.GetByID(s.ctx, userID)
	if err != nil {
		log.Println(err)
		return "", "", errno.ServiceErr
	}

	if usr == nil {
		return "", "", errno.UserIsNotExistErr
	}

	if usr.JwchID == nil || usr.JwchPassword == nil {
		return "", "", errno.JwchNotBindErr
	}

	password, err := crypt.Decrypt(*usr.JwchPassword)
	if err != nil {
		log.Printf("password decrypt err %v", err)
		return "", "", errno.ServiceErr
	}

	stu := jwch.NewStudent()
	stu.ID = *usr.JwchID
	stu.Password = password

	err = stu.Login()
	if err != nil {
		return "", "", errno.JwchLoginFailedErr.WithError(err)
	}

	id, cookies, err := stu.GetIdentifierAndCookies()
	if err != nil {
		return "", "", err
	}

	var cookiesStr string

	for _, cookie := range cookies {
		cookiesStr += cookie.Name + "=" + cookie.Value + ";"
	}

	return id, cookiesStr, nil
}
