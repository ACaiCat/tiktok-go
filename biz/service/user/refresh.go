package service

import (
	"github.com/ACaiCat/tiktok-go/biz/model/tiktok-go/user"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/jwt"
)

func (s *UserService) RefreshToken(req *user.RefreshReq) (string, string, error) {
	userID, err := jwt.ValidateToken(req.RefreshToken, constants.TypeRefreshToken)
	if err != nil {
		return "", "", err
	}

	accessToken, err := jwt.CreateToken(constants.TypeAccessToken, userID)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.CreateToken(constants.TypeRefreshToken, userID)

	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
