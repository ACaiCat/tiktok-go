package errno

var (
	SuccessMessage = "OK"
)

var (
	Success               = NewErrNo(SuccessCode, "成功")
	ServiceErr            = NewErrNo(ServiceErrCode, "服务内部错误")
	ParamErr              = NewErrNo(ParamErrCode, "参数错误")
	AuthErr               = NewErrNo(AuthErrCode, "认证失败")
	AuthMissingErr        = NewErrNo(AuthMissingErrCode, "认证信息缺失")
	AuthAccessExpiredErr  = NewErrNo(AuthAccessExpiredErrCode, "访问令牌过期")
	AuthRefreshExpiredErr = NewErrNo(AuthRefreshExpiredErrCode, "刷新令牌过期")
	TooManyRequestsErr    = NewErrNo(TooManyRequestsErrCode, "服务器过载，请稍后重试")

	UserAlreadyExistErr   = NewErrNo(UserAlreadyExistErrCode, "用户已存在")
	UserIsNotExistErr     = NewErrNo(UserIsNotExistErrCode, "用户不存在")
	PasswordIsNotVerified = NewErrNo(PasswordIsNotVerifiedErrCode, "用户名或密码错误")
	PasswordTooShortErr   = NewErrNo(PasswordTooShortErrCode, "密码长度太短")
	MFACodeInvalidErr     = NewErrNo(MFACodeInvalidErrCode, "多因素认证码无效")
	MFAMissingErr         = NewErrNo(MFAMissingErrCode, "缺少多因素认证码")
	AvatarTooLargeErr     = NewErrNo(AvatarTooLargeErrCode, "头像文件过大")
	AvatarFormatErr       = NewErrNo(AvatarFormatErrCode, "头像文件格式错误")

	UsernameTooShortErr = NewErrNo(UsernameTooShortErrCode, "用户名长度太短")

	UsernameTooLongErr        = NewErrNo(UsernameTooLongErrCode, "用户名长度太长")
	PasswordTooLongErr        = NewErrNo(PasswordTooLongErrCode, "密码长度太长")
	NotSupportActionErr       = NewErrNo(NotSupportActionErrCode, "不支持的操作")
	LikeAlreadyExistErr       = NewErrNo(LikeAlreadyExistErrCode, "已经点赞过了")
	LikeNotExistErr           = NewErrNo(LikeNotExistErrCode, "没有点赞过")
	VideoNotExistErr          = NewErrNo(VideoNotExistErrCode, "视频不存在")
	CommentNotExistErr        = NewErrNo(CommentNotExistErrCode, "评论不存在")
	CommentNotBelongToUserErr = NewErrNo(CommentNotBelongToUserErrCode, "评论不属于用户")

	FollowAlreadyExistErr = NewErrNo(FollowAlreadyExistErrCode, "已经关注过了")
	FollowNotExistErr     = NewErrNo(FollowNotExistErrCode, "没有关注过")
	FollowSelfErr         = NewErrNo(FollowSelfErrCode, "不能关注自己")

	ChatMsgParseErr  = NewErrNo(ChatMsgParseErrCode, "消息格式错误")
	ChatMsgTypeErr   = NewErrNo(ChatMsgTypeErrCode, "未知的消息类型")
	ChatNotFriendErr = NewErrNo(ChatNotFriendErrCode, "你们还不是好友")

	JwchLoginFailedErr  = NewErrNo(JwchLoginFailedErrCode, "教务处登录失败")
	JwchLoginExpiredErr = NewErrNo(JwchLoginExpiredErrCode, "教务处登录过期")
	JwchNotBindErr      = NewErrNo(JwchNotBindErrCode, "未绑定教务处账号")
)
