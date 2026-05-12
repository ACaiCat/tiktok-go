package user

import (
	"bytes"
	"context"
	"mime/multipart"
	"testing"

	"github.com/ACaiCat/tiktok-go/biz/model/user"
	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	userService "github.com/ACaiCat/tiktok-go/biz/service/user"
)

func buildAvatarForm(fileName string) (*bytes.Buffer, string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_, _ = w.CreateFormFile("data", fileName)
	_ = w.Close()
	return &buf, w.FormDataContentType()
}

func TestRegister(t *testing.T) {
	type testCase struct {
		body           string
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			body:           `{"username":"testuser","password":"testpass"}`,
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"missing param": {
			body:           `{}`,
			expectContains: `"base":{"code":10002,"msg":"参数错误:`,
		},
		"internal server error": {
			body:           `{"username":"testuser","password":"testpass"}`,
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/user/register", Register)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(userService.NewUserService).To(func(_ context.Context) *userService.UserService {
				return &userService.UserService{}
			}).Build()
			mockey.Mock((*userService.UserService).UserRegister).To(func(_ *user.RegisterReq) error {
				return tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodPost, "/user/register",
				&ut.Body{Body: bytes.NewBufferString(tc.body), Len: len(tc.body)},
				ut.Header{Key: "Content-Type", Value: "application/json"},
			)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestLogin(t *testing.T) {
	type testCase struct {
		body           string
		mockUser       *model.User
		mockAccess     string
		mockRefresh    string
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			body:           `{"username":"testuser","password":"testpass"}`,
			mockUser:       &model.User{ID: "1", Username: "testuser"},
			mockAccess:     "access-token",
			mockRefresh:    "refresh-token",
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"missing param": {
			body:           `{}`,
			expectContains: `"base":{"code":10002,"msg":"参数错误:`,
		},
		"internal server error": {
			body:           `{"username":"testuser","password":"testpass"}`,
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/user/login", Login)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(userService.NewUserService).To(func(_ context.Context) *userService.UserService {
				return &userService.UserService{}
			}).Build()
			mockey.Mock((*userService.UserService).UserLogin).To(func(_ *user.LoginReq) (*model.User, string, string, error) {
				return tc.mockUser, tc.mockAccess, tc.mockRefresh, tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodPost, "/user/login",
				&ut.Body{Body: bytes.NewBufferString(tc.body), Len: len(tc.body)},
				ut.Header{Key: "Content-Type", Value: "application/json"},
			)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestRefresh(t *testing.T) {
	type testCase struct {
		body           string
		mockAccess     string
		mockRefresh    string
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			body:           `{"refresh_token":"old-refresh"}`,
			mockAccess:     "new-access",
			mockRefresh:    "new-refresh",
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"missing param": {
			body:           `{}`,
			expectContains: `"base":{"code":10002,"msg":"参数错误:`,
		},
		"internal server error": {
			body:           `{"refresh_token":"old-refresh"}`,
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/auth/refresh", Refresh)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(userService.NewUserService).To(func(_ context.Context) *userService.UserService {
				return &userService.UserService{}
			}).Build()
			mockey.Mock((*userService.UserService).RefreshToken).To(func(_ *user.RefreshReq) (string, string, error) {
				return tc.mockAccess, tc.mockRefresh, tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodPost, "/auth/refresh",
				&ut.Body{Body: bytes.NewBufferString(tc.body), Len: len(tc.body)},
				ut.Header{Key: "Content-Type", Value: "application/json"},
			)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestInfo(t *testing.T) {
	type testCase struct {
		url            string
		mockUser       *model.User
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			url:            "/user/info?user_id=1",
			mockUser:       &model.User{ID: "1", Username: "testuser"},
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"internal server error": {
			url:            "/user/info?user_id=1",
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/user/info", Info)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(userService.NewUserService).To(func(_ context.Context) *userService.UserService {
				return &userService.UserService{}
			}).Build()
			mockey.Mock((*userService.UserService).GetUserInfo).To(func(_ int64) (*model.User, error) {
				return tc.mockUser, tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestUploadAvatar(t *testing.T) {
	type testCase struct {
		mockUploadErr  error
		mockInfoErr    error
		mockUser       *model.User
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			mockUser:       &model.User{ID: "1", Username: "testuser"},
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"upload error": {
			mockUploadErr:  assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
		"get info error": {
			mockInfoErr:    assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.PUT("/user/avatar/upload", UploadAvatar)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(userService.NewUserService).To(func(_ context.Context) *userService.UserService {
				return &userService.UserService{}
			}).Build()
			mockey.Mock((*userService.UserService).UploadAvatar).To(func(_ *multipart.FileHeader, userID int64) error {
				return tc.mockUploadErr
			}).Build()
			mockey.Mock((*userService.UserService).GetUserInfo).To(func(_ int64) (*model.User, error) {
				return tc.mockUser, tc.mockInfoErr
			}).Build()

			buf, contentType := buildAvatarForm("avatar.jpg")
			result := ut.PerformRequest(router, consts.MethodPut, "/user/avatar/upload",
				&ut.Body{Body: buf, Len: buf.Len()},
				ut.Header{Key: "Content-Type", Value: contentType},
			)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestMFAQRCode(t *testing.T) {
	type testCase struct {
		mockSecret     string
		mockQrcode     string
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			mockSecret:     "secret",
			mockQrcode:     "base64qrcode",
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"internal server error": {
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/auth/mfa/qrcode", MFAQRCode)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(userService.NewUserService).To(func(_ context.Context) *userService.UserService {
				return &userService.UserService{}
			}).Build()
			mockey.Mock((*userService.UserService).GetMFA).To(func(_ int64) (string, string, error) {
				return tc.mockSecret, tc.mockQrcode, tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodGet, "/auth/mfa/qrcode", nil)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestBindMFA(t *testing.T) {
	type testCase struct {
		body           string
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			body:           `{"code":"123456","secret":"mysecret"}`,
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"missing param": {
			body:           `{}`,
			expectContains: `"base":{"code":10002,"msg":"参数错误:`,
		},
		"internal server error": {
			body:           `{"code":"123456","secret":"mysecret"}`,
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/auth/mfa/bind", BindMFA)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(userService.NewUserService).To(func(_ context.Context) *userService.UserService {
				return &userService.UserService{}
			}).Build()
			mockey.Mock((*userService.UserService).BindMFA).To(func(_ *user.BindMFAReq, _ int64) error {
				return tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodPost, "/auth/mfa/bind",
				&ut.Body{Body: bytes.NewBufferString(tc.body), Len: len(tc.body)},
				ut.Header{Key: "Content-Type", Value: "application/json"},
			)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestBindJwch(t *testing.T) {
	type testCase struct {
		body           string
		mockErr        error
		expectContains string
	}

	testCases := map[string]testCase{
		"success": {
			body:           `{"jwch_id":"student123","jwch_password":"pass"}`,
			expectContains: `"base":{"code":10000,"msg":"OK"}`,
		},
		"missing param": {
			body:           `{}`,
			expectContains: `"base":{"code":10002,"msg":"参数错误:`,
		},
		"internal server error": {
			body:           `{"jwch_id":"student123","jwch_password":"pass"}`,
			mockErr:        assert.AnError,
			expectContains: `"base":{"code":10001,"msg":"服务器内部错误"}`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/user/jwch/bind", BindJwch)

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(userService.NewUserService).To(func(_ context.Context) *userService.UserService {
				return &userService.UserService{}
			}).Build()
			mockey.Mock((*userService.UserService).BindJwch).To(func(_ *user.BindJwchReq, _ int64) error {
				return tc.mockErr
			}).Build()

			result := ut.PerformRequest(router, consts.MethodPost, "/user/jwch/bind",
				&ut.Body{Body: bytes.NewBufferString(tc.body), Len: len(tc.body)},
				ut.Header{Key: "Content-Type", Value: "application/json"},
			)
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}
