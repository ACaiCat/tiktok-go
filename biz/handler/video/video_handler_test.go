package video

import (
	"bytes"
	"context"
	"mime/multipart"
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

func buildPublishForm(fileName string, title string, description string) (*bytes.Buffer, string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_, _ = w.CreateFormFile("data", fileName)
	_ = w.WriteField("title", title)
	_ = w.WriteField("description", description)
	_ = w.Close()
	return &buf, w.FormDataContentType()
}

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
			expectContains: `"base":{"code":10000,"msg":"OK"},"data":{"items":[]}`,
		},
		"internal server error": {
			url:            "/video/feed",
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
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
			mockey.Mock((*videoService.VideoService).GetFeed).To(func(_ *video.FeedReq) ([]*model.Video, error) {
				return tc.mockResult, tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestPublish(t *testing.T) {
	type testCase struct {
		url            string
		title          string
		description    string
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			url:            "/video/publish",
			title:          "mock_title",
			description:    "mock_description",
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"missing param": {
			url:            "/video/publish",
			title:          "",
			description:    "",
			expectContains: `"base":{"code":10002,"msg":"参数错误:`,
		},
		"internal server error": {
			url:            "/video/publish",
			title:          "mock_title",
			description:    "mock_description",
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}
	router := route.NewEngine(&config.Options{})
	router.POST("/video/publish", Publish)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(videoService.NewVideoService).To(func(_ context.Context) *videoService.VideoService {
				return &videoService.VideoService{}
			}).Build()
			mockey.Mock((*videoService.VideoService).PublishVideo).To(func(_ int64, _ string, _ string, _ *multipart.FileHeader) error {
				return tc.mockErr
			}).Build()

			buf, contentType := buildPublishForm("mockfile", tc.title, tc.description)
			result := ut.PerformRequest(router, consts.MethodPost, tc.url, &ut.Body{
				Body: buf, Len: buf.Len(),
			}, ut.Header{Key: "Content-Type", Value: contentType})

			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}
