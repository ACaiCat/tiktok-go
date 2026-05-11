package video

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/video"
	videoService "github.com/ACaiCat/tiktok-go/biz/service/video"
)

func TestFeed(t *testing.T) {
	type testCase struct {
		url            string
		mockResult     []*model.Video
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			url:            "/video/feed",
			mockResult:     []*model.Video{},
			expectContains: `"base":{"msg":"OK","code":10000},"data":{"items":[]}`,
		},
		"internal server error": {
			url:            "/video/feed",
			mockErr:        assert.AnError,
			expectContains: `"base":{"msg":"服务器内部错误","code":10001}`,
		},
	}
	router := route.NewEngine(&config.Options{})
	router.GET("/video/feed", Feed)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(videoService.NewVideoService).To(func(_ context.Context) *videoService.VideoService {
				return &videoService.VideoService{}
			}).Build()
			mockey.Mock((*videoService.VideoService).GetFeed).To(func(req *video.FeedReq) ([]*model.Video, error) {
				return tc.mockResult, tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}
