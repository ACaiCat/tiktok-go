package pack

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/ACaiCat/tiktok-go/api/model/model"
	"github.com/ACaiCat/tiktok-go/api/model/user"
)

func RespRefreshToken(c *app.RequestContext, accessToken string, refreshToken string) {
	c.JSON(consts.StatusOK, user.RefreshResp{
		Base: SuccessBase,
		Data: &user.TokenData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

func RespLogin(c *app.RequestContext, usr *model.User) {
	c.JSON(consts.StatusOK, user.LoginResp{
		Base: SuccessBase,
		Data: usr,
	})
}

func RespRegister(c *app.RequestContext) {
	c.JSON(consts.StatusOK, user.RegisterResp{
		Base: SuccessBase,
	})
}

func RespMFA(c *app.RequestContext, secret string, base64Qrcode string) {
	c.JSON(consts.StatusOK, user.MFAQRCodeResp{
		Base: SuccessBase,
		Data: &user.MFAQRCodeData{
			Secret: secret,
			Qrcode: base64Qrcode,
		},
	})
}

func RespBindMFA(c *app.RequestContext) {
	c.JSON(consts.StatusOK, user.BindMFAResp{
		Base: SuccessBase,
	})
}

func RespUploadAvatar(c *app.RequestContext, usr *model.User) {
	c.JSON(consts.StatusOK, user.UploadAvatarResp{
		Base: SuccessBase,
		Data: usr,
	})
}

func RespUserInfo(c *app.RequestContext, usr *model.User) {
	c.JSON(consts.StatusOK, user.InfoResp{
		Base: SuccessBase,
		Data: usr,
	})
}
