package mw

import (
	"context"
	"strings"

	"github.com/ACaiCat/tiktok-go/biz/pack"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/jwt"
	"github.com/cloudwego/hertz/pkg/app"
)

func Auth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := strings.TrimPrefix(string(c.GetHeader(constants.AuthHeader)), "Bearer ")

		userID, err := jwt.ValidateToken(token, constants.TypeAccessToken)
		if err != nil {
			pack.RespError(c, err)
			c.Abort()
			return
		}

		c.Set(constants.UserIdKey, userID)
		c.Next(ctx)

	}
}
