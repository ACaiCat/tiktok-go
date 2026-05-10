package mw

import (
	"context"
	"sync"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzSentinel "github.com/hertz-contrib/opensergo/sentinel/adapter"

	"github.com/ACaiCat/tiktok-go/biz/pack"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

var once sync.Once

func Sentinel() app.HandlerFunc {
	once.Do(initSentinel)

	return hertzSentinel.SentinelServerMiddleware(
		hertzSentinel.WithServerResourceExtractor(func(c context.Context, ctx *app.RequestContext) string {
			return "api"
		}),
		hertzSentinel.WithServerBlockFallback(func(ctx context.Context, c *app.RequestContext) {
			pack.RespError(ctx, c, errno.TooManyRequestsErr)
			c.Abort()
		}),
	)
}

func initSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		hlog.Fatal(err)
	}

	qps := 80.0
	_, err = flow.LoadRules([]*flow.Rule{
		&flow.Rule{
			Resource:               "api",
			Threshold:              qps,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		},
	})
	if err != nil {
		hlog.Fatal(err)
	}
}
