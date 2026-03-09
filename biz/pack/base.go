package pack

import (
	"github.com/ACaiCat/tiktok-go/biz/model/common"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

var (
	SuccessBase = &common.Base{
		Msg:  errno.SuccessMessage,
		Code: errno.SuccessCode,
	}
)

func RespError(c *app.RequestContext, err error) {
	_errno := errno.ConvertErr(err)

	c.JSON(consts.StatusOK, common.Base{
		Msg:  _errno.ErrMsg,
		Code: _errno.ErrCode,
	})
}
