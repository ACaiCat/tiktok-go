package mw

import (
	"context"

	"github.com/ACaiCat/tiktok-go/biz/pack"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/jwt"
	"github.com/cloudwego/hertz/pkg/app"
)

func Auth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := string(c.GetHeader(constants.AccessTokenHeader))
		userID, err := jwt.ValidateToken(token, constants.TypeAccessToken)
		if err != nil {
			pack.RespError(c, err)
		}

		c.Set(constants.UserIdKey, userID)
		c.Next(ctx)

	}
}
