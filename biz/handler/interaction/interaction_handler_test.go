package interaction

import (
	"bytes"
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/biz/model/interaction"
	"github.com/ACaiCat/tiktok-go/biz/model/model"
	interactionService "github.com/ACaiCat/tiktok-go/biz/service/interaction"
	videoService "github.com/ACaiCat/tiktok-go/biz/service/video"
)

func TestLike(t *testing.T) {
	type testCase struct {
		body           string
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			body:           `{"action_type":1}`,
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"missing param": {
			body:           `{}`,
			expectContains: `"base":{"code":10002,"msg":"参数错误:`,
		},
		"internal server error": {
			body:           `{"action_type":1}`,
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/like/action", Like)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(interactionService.NewInteractionService).To(func(_ context.Context) *interactionService.InteractionService {
				return &interactionService.InteractionService{}
			}).Build()
			mockey.Mock((*interactionService.InteractionService).LikeVideo).To(func(_ *interaction.LikeReq, _ int64) error {
				return tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodPost, "/like/action",
				&ut.Body{Body: bytes.NewBufferString(tc.body), Len: len(tc.body)},
				ut.Header{Key: "Content-Type", Value: "application/json"},
			)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestListLike(t *testing.T) {
	type testCase struct {
		url            string
		mockResult     []*model.Video
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			url:            "/like/list?user_id=1&page_num=0&page_size=10",
			mockResult:     []*model.Video{},
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"missing param": {
			url:            "/like/list",
			expectContains: `"base":{"code":10002,"msg":"参数错误:`,
		},
		"internal server error": {
			url:            "/like/list?user_id=1&page_num=0&page_size=10",
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/like/list", ListLike)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(videoService.NewVideoService).To(func(_ context.Context) *videoService.VideoService {
				return &videoService.VideoService{}
			}).Build()
			mockey.Mock((*videoService.VideoService).GetLikedVideos).To(func(_ *interaction.ListLikeReq) ([]*model.Video, error) {
				return tc.mockResult, tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestComment(t *testing.T) {
	type testCase struct {
		body           string
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			body:           `{"content":"great video!"}`,
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"missing param": {
			body:           `{}`,
			expectContains: `"base":{"code":10002,"msg":"参数错误:`,
		},
		"internal server error": {
			body:           `{"content":"great video!"}`,
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/comment/publish", Comment)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(interactionService.NewInteractionService).To(func(_ context.Context) *interactionService.InteractionService {
				return &interactionService.InteractionService{}
			}).Build()
			mockey.Mock((*interactionService.InteractionService).CommentAction).To(func(_ *interaction.CommentReq, _ int64) error {
				return tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodPost, "/comment/publish",
				&ut.Body{Body: bytes.NewBufferString(tc.body), Len: len(tc.body)},
				ut.Header{Key: "Content-Type", Value: "application/json"},
			)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestListComment(t *testing.T) {
	type testCase struct {
		url            string
		mockResult     []*model.Comment
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			url:            "/comment/list?page_num=0&page_size=10",
			mockResult:     []*model.Comment{},
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"missing param": {
			url:            "/comment/list",
			expectContains: `"base":{"code":10002,"msg":"参数错误:`,
		},
		"internal server error": {
			url:            "/comment/list?page_num=0&page_size=10",
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/comment/list", ListComment)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(interactionService.NewInteractionService).To(func(_ context.Context) *interactionService.InteractionService {
				return &interactionService.InteractionService{}
			}).Build()
			mockey.Mock((*interactionService.InteractionService).ListComment).To(func(_ *interaction.ListCommentReq) ([]*model.Comment, error) {
				return tc.mockResult, tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestDeleteComment(t *testing.T) {
	type testCase struct {
		body           string
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			body:           `{"comment_id":"123"}`,
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"missing param": {
			body:           `{}`,
			expectContains: `"base":{"code":10002,"msg":"参数错误:`,
		},
		"internal server error": {
			body:           `{"comment_id":"123"}`,
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.DELETE("/comment/delete", DeleteComment)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(interactionService.NewInteractionService).To(func(_ context.Context) *interactionService.InteractionService {
				return &interactionService.InteractionService{}
			}).Build()
			mockey.Mock((*interactionService.InteractionService).DeleteComment).To(func(_ *interaction.DeleteCommentReq, _ int64) error {
				return tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodDelete, "/comment/delete",
				&ut.Body{Body: bytes.NewBufferString(tc.body), Len: len(tc.body)},
				ut.Header{Key: "Content-Type", Value: "application/json"},
			)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}
