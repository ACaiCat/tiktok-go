namespace go user

include "model.thrift"
include "common.thrift"

// 用户注册请求
struct RegisterReq {
  // 用户名
  1: required string username (api.form = 'username');
  // 密码
  2: required string password (api.form = 'password');
}

// 用户注册响应
struct RegisterResp {
  // 响应状态
  1: required common.Base base;
}

// 用户登录请求
struct LoginReq {
  // 用户名
  1: required string username (api.form= 'username') ;
  // 密码
  2: required string password (api.form = 'password');
  // 验证码
  3: optional string code (api.form = 'code');
}

// 用户登录响应
struct LoginResp {
  // 响应状态
  1: required common.Base base;
  // 用户数据
  2: optional model.User data;
}

// 刷新Token请求
struct RefreshReq {
  // 刷新Token
  1: required string refresh_token (api.body = 'refresh_token');
}

// 刷新Token响应
struct RefreshResp {
  // 响应状态
  1: required common.Base base;
  // Token数据
  2: optional TokenData data;
}

// Token数据
struct TokenData {
  // 访问Token
  1: required string access_token;
  // 刷新Token
  2: required string refresh_token;
}

// 用户信息请求
struct InfoReq {
  // 用户ID
  1: optional string UserID (api.query = 'user_id');
}

// 用户信息响应
struct InfoResp {
  // 响应状态
  1: required common.Base base;
  // 用户数据
  2: optional model.User data;
}

// 上传头像请求
struct UploadAvatarReq {
  // 图片文件
  1: optional binary data (api.form = 'data');
}

// 上传头像响应
struct UploadAvatarResp {
  // 响应状态
  1: required common.Base base;
  // 用户数据
  2: optional model.User data;
}

// MFA二维码请求
struct MFAQRCodeReq {
}

// MFA二维码请求响应
struct MFAQRCodeResp {
  // 响应状态
  1: required common.Base base;
  // 数据
  2: optional MFAQRCodeData data;
}

// MFA二维码数据
struct MFAQRCodeData {
  // 密钥
  1: required string secret;
  // 二维码Base64
  2: required string qrcode;
}

// 绑定MFA请求
struct BindMFAReq {
  // 校验码
  1: required string code (api.body = 'code');
  // 密钥
  2: required string secret (api.body = 'secret');
}

// 绑定MFA响应
struct BindMFAResp {
  // 响应状态
  1: required common.Base base;
}

// 以图搜图请求
struct SearchImageReq {
  // 图片原始数据
  1: optional binary data (api.form = 'data');
}

// 以图搜图响应
struct SearchImageResp {
  // 响应状态
  1: required common.Base base;
  // 图片URL
  2: required string data;
}

service UserHandler {
  // 用户注册
  RegisterResp Register(1: RegisterReq req) (api.post = "/user/register")

  // 用户登录
  LoginResp Login(1: LoginReq req) (api.post = "/user/login")

  // 刷新Token
  RefreshResp Refresh(1: RefreshReq req) (api.post = "/auth/refresh")

  // 用户信息
  InfoResp Info(1: InfoReq req) (api.get = "/user/info")

  // 上传头像
  UploadAvatarResp UploadAvatar(1: UploadAvatarReq req) (api.put = "/user/avatar/upload")

  // MFA二维码
  MFAQRCodeResp MFAQRCode(1: MFAQRCodeReq req) (api.get = "/auth/mfa/qrcode")

  // 绑定MFA
  BindMFAResp BindMFA(1: BindMFAReq req) (api.post = "/auth/mfa/bind")

  // 以图搜图
  SearchImageResp SearchImage(1: SearchImageReq req) (api.post = "/user/image/search")
}