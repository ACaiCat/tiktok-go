package mw

import (
	"github.com/cloudwego/hertz/pkg/app"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func GetUserID(c *app.RequestContext) int64 {
	id, _ := c.Get(constants.UserIDKey)
	userID, ok := id.(int64)
	if !ok {
		return 0
	}
	return userID
}
