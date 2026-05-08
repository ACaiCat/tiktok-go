package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/ACaiCat/tiktok-go/config"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func CreateToken(tokenType int8, userID int64) (string, error) {
	var expireTime time.Time
	nowTime := time.Now()

	switch tokenType {
	case constants.TypeAccessToken:
		expireTime = nowTime.Add(constants.AccessTokenExpiration)
	case constants.TypeRefreshToken:
		expireTime = nowTime.Add(constants.RefreshTokenExpiration)
	default:
		return "", errno.AuthErr.WithMessage("invalid token type")
	}

	claims := &Claims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    constants.TokenIssuer,
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	var secret string

	if tokenType == constants.TypeAccessToken {
		secret = config.AppConfig.JWT.AccessSecret
	} else {
		secret = config.AppConfig.JWT.RefreshSecret
	}

	tokenObject := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenObject.SignedString([]byte(secret))
	if err != nil {
		return "", errno.AuthErr.WithMessage("令牌签名失败")
	}

	return token, nil
}

func ValidateToken(token string, tokenType int8) (int64, error) {
	if token == "" {
		return 0, errno.AuthMissingErr
	}

	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		switch claims.TokenType {
		case constants.TypeAccessToken:
			return []byte(config.AppConfig.JWT.AccessSecret), nil
		case constants.TypeRefreshToken:
			return []byte(config.AppConfig.JWT.RefreshSecret), nil
		default:
			return 0, errno.AuthErr.WithMessage("令牌无效")
		}
	})

	if claims.TokenType != tokenType {
		return 0, errno.AuthErr.WithMessage("令牌类型不匹配")
	}

	if err != nil || !parsedToken.Valid {
		if errors.Is(err, jwt.ErrTokenExpired) {
			if tokenType == constants.TypeAccessToken {
				return 0, errno.AuthAccessExpiredErr
			}
			return 0, errno.AuthRefreshExpiredErr
		}

		return 0, errno.AuthErr.WithMessage("令牌无效")
	}

	return claims.UserID, nil
}
