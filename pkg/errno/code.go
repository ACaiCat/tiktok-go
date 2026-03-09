package errno

const (
	SuccessCode                  = 10000
	ServiceErrCode               = 10001
	ParamErrCode                 = 10002
	AuthErrCode                  = 10004
	AuthMissingErrCode           = 10005
	AuthAccessExpiredErrCode     = 10006
	AuthRefreshExpiredErrCode    = 10007
	UserAlreadyExistErrCode      = 11000
	UserIsNotExistErrCode        = 11001
	PasswordIsNotVerifiedErrCode = 11002
	PasswordTooShortErrCode      = 11003
	MFACodeInvalidErrCode        = 11004
	MFAMissingErrCode            = 11005
)
