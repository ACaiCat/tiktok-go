package mw

import (
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/cloudwego/hertz/pkg/app"
)

func GetUserID(c *app.RequestContext) int64 {
	id, _ := c.Get(constants.UserIdKey)

	return id.(int64)
}
