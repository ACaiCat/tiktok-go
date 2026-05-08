package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/west2-online/jwch"

	"github.com/ACaiCat/tiktok-go/biz/model/user"
	"github.com/ACaiCat/tiktok-go/pkg/crypt"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/ACaiCat/tiktok-go/pkg/utils"
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
		return errors.WithMessagef(err, "service.BindJwch: crypt.Encrypt failed, userID=%d", userID)
	}

	err = s.dao.UpdateUserJwch(s.ctx, userID, req.JwchID, password)
	if err != nil {
		return errors.WithMessagef(err, "service.BindJwch: db.UpdateUserJwch failed, userID=%d", userID)
	}

	go func() {
		err := s.cache.CleanJwchSession(context.Background(), userID)
		if err != nil {
			hlog.Errorf("service.BindJwch cache clean failed: %v", err)
			return
		}
	}()

	return nil
}

func (s *UserService) GetJwchIdentifierAndCookies(userID int64) (string, string, error) {
	var err error

	idCache, cookieCache, err := s.cache.GetJwchSession(s.ctx, userID)

	if err == nil {
		stu := jwch.NewStudent().WithLoginData(idCache, utils.ParseCookies(cookieCache))
		err := stu.CheckSession()
		if err == nil {
			return idCache, cookieCache, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		hlog.CtxErrorf(s.ctx, "service.GetJwchIdentifierAndCookies cache read failed: %v", err)
	}

	usr, err := s.dao.GetByID(s.ctx, userID)
	if err != nil {
		return "", "", errors.WithMessagef(err, "service.GetJwchIdentifierAndCookies: db.GetByID failed, userID=%d", userID)
	}

	if usr == nil {
		return "", "", errno.UserIsNotExistErr
	}

	if usr.JwchID == nil || usr.JwchPassword == nil {
		return "", "", errno.JwchNotBindErr
	}

	password, err := crypt.Decrypt(*usr.JwchPassword)
	if err != nil {
		return "", "", errors.WithMessagef(err, "service.GetJwchIdentifierAndCookies: crypt.Decrypt failed, userID=%d", userID)
	}

	stu := jwch.NewStudent().WithUser(*usr.JwchID, password)
	err = stu.Login()
	if err != nil {
		return "", "", errno.JwchLoginFailedErr.WithError(err)
	}

	id, cookies, err := stu.GetIdentifierAndCookies()
	if err != nil {
		return "", "", err
	}

	cookiesStr := utils.ParseCookiesToString(cookies)

	go func() {
		err := s.cache.SetJwchSession(context.Background(), userID, id, cookiesStr)
		if err != nil {
			hlog.Errorf("service.GetJwchIdentifierAndCookies cache write failed: %v", err)
		}
	}()

	return id, cookiesStr, nil
}
