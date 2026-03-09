package pack

import (
	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/tiktok-go/user"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
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
