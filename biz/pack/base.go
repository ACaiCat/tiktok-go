package pack

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/ACaiCat/tiktok-go/biz/model/common"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

var (
	SuccessBase = &common.Base{
		Msg:  errno.SuccessMessage,
		Code: errno.SuccessCode,
	}
)

func RespError(c *app.RequestContext, err error) {
	_errno := errno.ConvertErr(err)

	c.JSON(consts.StatusOK, utils.H{
		"base": common.Base{
			Msg:  _errno.ErrMsg,
			Code: _errno.ErrCode,
		},
	})
}
