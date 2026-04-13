package constants

import "time"

const (
	TypeAccessToken        int8 = 0
	TypeRefreshToken       int8 = 1
	RefreshTokenExpiration      = 7 * 24 * time.Hour
	AccessTokenExpiration       = 2 * time.Hour
	TokenIssuer                 = "Cai"
	AuthHeader                  = "Authorization"
	AccessTokenHeader           = "Access-Token"
	RefreshTokenHeader          = "Refresh-Token"
	UserIDKey                   = "user_id"

	MinPasswordLength = 6
	MaxPasswordLength = 64

	MinUsernameLength = 3
	MaxUsernameLength = 32
)
