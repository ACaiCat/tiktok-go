package mw

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func AccessLog() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()

		c.Next(ctx)

		end := time.Now()
		latency := end.Sub(start)
		statusCode := c.Response.StatusCode()
		clientIP := c.ClientIP()
		method := string(c.Request.Header.Method())
		path := string(c.Request.URI().PathOriginal())

		hlog.CtxInfof(ctx, "[HERTZ] %3d | %8v | %12s | %-4s %s",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
