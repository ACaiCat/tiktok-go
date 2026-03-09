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

	UserAlreadyExistErr   = NewErrNo(UserAlreadyExistErrCode, "用户已存在")
	UserIsNotExistErr     = NewErrNo(UserIsNotExistErrCode, "用户不存在")
	PasswordIsNotVerified = NewErrNo(PasswordIsNotVerifiedErrCode, "用户名或密码错误")
	PasswordTooShortErr   = NewErrNo(PasswordTooShortErrCode, "密码长度太短")
)
