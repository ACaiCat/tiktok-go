package social

import (
	"bytes"
	"context"
	"testing"

	"github.com/ACaiCat/tiktok-go/biz/model/social"
	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	socialService "github.com/ACaiCat/tiktok-go/biz/service/social"
)

func TestFollow(t *testing.T) {
	type testCase struct {
		body           string
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			body:           `{"to_user_id":"2","action_type":1}`,
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"missing param": {
			body:           `{}`,
			expectContains: `"base":{"code":10002,"msg":"参数错误:`,
		},
		"internal server error": {
			body:           `{"to_user_id":"2","action_type":1}`,
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/relation/action", Follow)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(socialService.NewSocialService).To(func(_ context.Context) *socialService.SocialService {
				return &socialService.SocialService{}
			}).Build()
			mockey.Mock((*socialService.SocialService).FollowAction).To(func(_ *social.FollowReq, _ int64) error {
				return tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodPost, "/relation/action",
				&ut.Body{Body: bytes.NewBufferString(tc.body), Len: len(tc.body)},
				ut.Header{Key: "Content-Type", Value: "application/json"},
			)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestListFollowing(t *testing.T) {
	type testCase struct {
		url            string
		mockResult     []*model.SocialUser
		mockTotal      int
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			url:            "/following/list?user_id=1&page_num=0&page_size=10",
			mockResult:     []*model.SocialUser{},
			mockTotal:      0,
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"missing param": {
			url:            "/following/list",
			expectContains: `"base":{"code":10002,"msg":"参数错误:`,
		},
		"internal server error": {
			url:            "/following/list?user_id=1&page_num=0&page_size=10",
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/following/list", ListFollowing)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(socialService.NewSocialService).To(func(_ context.Context) *socialService.SocialService {
				return &socialService.SocialService{}
			}).Build()
			mockey.Mock((*socialService.SocialService).ListFollowing).To(func(_ *social.ListFollowingReq) ([]*model.SocialUser, int, error) {
				return tc.mockResult, tc.mockTotal, tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestListFollower(t *testing.T) {
	type testCase struct {
		url            string
		mockResult     []*model.SocialUser
		mockTotal      int
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			url:            "/follower/list?user_id=1&page_num=0&page_size=10",
			mockResult:     []*model.SocialUser{},
			mockTotal:      0,
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"missing param": {
			url:            "/follower/list",
			expectContains: `"base":{"code":10002,"msg":"参数错误:`,
		},
		"internal server error": {
			url:            "/follower/list?user_id=1&page_num=0&page_size=10",
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/follower/list", ListFollower)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(socialService.NewSocialService).To(func(_ context.Context) *socialService.SocialService {
				return &socialService.SocialService{}
			}).Build()
			mockey.Mock((*socialService.SocialService).ListFollower).To(func(_ *social.ListFollowerReq) ([]*model.SocialUser, int, error) {
				return tc.mockResult, tc.mockTotal, tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestListFriend(t *testing.T) {
	type testCase struct {
		url            string
		mockResult     []*model.SocialUser
		mockTotal      int
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			url:            "/friends/list?page_num=0&page_size=10",
			mockResult:     []*model.SocialUser{},
			mockTotal:      0,
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"missing param": {
			url:            "/friends/list",
			expectContains: `"base":{"code":10002,"msg":"参数错误:`,
		},
		"internal server error": {
			url:            "/friends/list?page_num=0&page_size=10",
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/friends/list", ListFriend)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(socialService.NewSocialService).To(func(_ context.Context) *socialService.SocialService {
				return &socialService.SocialService{}
			}).Build()
			mockey.Mock((*socialService.SocialService).ListFriend).To(func(_ *social.ListFriendReq, _ int64) ([]*model.SocialUser, int, error) {
				return tc.mockResult, tc.mockTotal, tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}
