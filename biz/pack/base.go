package pack

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
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

func RespError(ctx context.Context, c *app.RequestContext, err error) {
	_errno := errno.ConvertErr(err)

	if _errno.ErrCode == errno.ServiceErrCode {
		hlog.CtxErrorf(ctx,
			"[%s] %s: %+v",
			string(c.Method()),
			string(c.Path()),
			err,
		)
	}

	c.JSON(consts.StatusOK, utils.H{
		"base": common.Base{
			Msg:  _errno.ErrMsg,
			Code: _errno.ErrCode,
		},
	})
}
