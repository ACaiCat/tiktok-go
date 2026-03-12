package errno

const (
	SuccessCode                   = 10000
	ServiceErrCode                = 10001
	ParamErrCode                  = 10002
	AuthErrCode                   = 10004
	AuthMissingErrCode            = 10005
	AuthAccessExpiredErrCode      = 10006
	AuthRefreshExpiredErrCode     = 10007
	UserAlreadyExistErrCode       = 11000
	UserIsNotExistErrCode         = 11001
	PasswordIsNotVerifiedErrCode  = 11002
	PasswordTooShortErrCode       = 11003
	MFACodeInvalidErrCode         = 11004
	MFAMissingErrCode             = 11005
	AvatarTooLargeErrCode         = 12000
	AvatarFormatErrCode           = 12001
	UsernameTooShortErrCode       = 12002
	UsernameTooLongErrCode        = 12003
	PasswordTooLongErrCode        = 12004
	NotSupportActionErrCode       = 13000
	LikeAlreadyExistErrCode       = 14000
	LikeNotExistErrCode           = 14001
	VideoNotExistErrCode          = 15000
	CommentNotExistErrCode        = 15001
	CommentNotBelongToUserErrCode = 15002
	FollowAlreadyExistErrCode     = 16000
	FollowNotExistErrCode         = 16001
	FollowSelfErrCode             = 16002
	ChatMsgParseErrCode           = 17000
	ChatMsgTypeErrCode            = 17001
	ChatNotFriendErrCode          = 17002
)
