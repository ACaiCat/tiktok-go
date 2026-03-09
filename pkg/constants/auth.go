package constants

import "time"

const (
	TypeAccessToken        int8 = 0
	TypeRefreshToken       int8 = 1
	RefreshTokenExpiration      = 7 * 24 * time.Hour
	AccessTokenExpiration       = 2 * time.Hour
	TokenIssuer                 = "Cai"
	AccessTokenHeader           = "Access-Token"
	RefreshTokenHeader          = "Refresh-Token"
	UserIdKey                   = "user_id"

	MinPasswordLength = 6
)
