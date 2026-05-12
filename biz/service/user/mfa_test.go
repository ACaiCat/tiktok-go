package service

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/pquerna/otp"
	"github.com/skip2/go-qrcode"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/biz/model/user"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	userDao "github.com/ACaiCat/tiktok-go/pkg/db/user"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	totpTool "github.com/ACaiCat/tiktok-go/pkg/totp"
)

func TestUserService_GetMFA(t *testing.T) {
	type testCase struct {
		mockUser      *modelDao.User
		mockGetErr    error
		mockKeyErr    error
		mockQRCodeErr error
		expectError   string
	}

	testCases := map[string]testCase{
		"success": {
			mockUser: &modelDao.User{ID: 1, Username: "testuser"},
		},
		"user not found": {
			expectError: errno.UserIsNotExistErr.ErrMsg,
		},
		"db error": {
			mockGetErr:  assert.AnError,
			expectError: assert.AnError.Error(),
		},
		"create key error": {
			mockUser:    &modelDao.User{ID: 1, Username: "testuser"},
			mockKeyErr:  assert.AnError,
			expectError: assert.AnError.Error(),
		},
		"qrcode error": {
			mockUser:      &modelDao.User{ID: 1, Username: "testuser"},
			mockQRCodeErr: assert.AnError,
			expectError:   assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*userDao.UserDao).GetByID).To(
				func(_ *userDao.UserDao, ctx context.Context, userID int64) (*modelDao.User, error) {
					return tc.mockUser, tc.mockGetErr
				}).Build()
			mockey.Mock(totpTool.CreateKey).To(
				func(username string) (*otp.Key, error) {
					if tc.mockKeyErr != nil {
						return nil, tc.mockKeyErr
					}
					return &otp.Key{}, nil
				}).Build()
			mockey.Mock((*otp.Key).String).To(func(_ *otp.Key) string {
				return "otpauth://totp/test"
			}).Build()
			mockey.Mock((*otp.Key).Secret).To(func(_ *otp.Key) string {
				return "secret"
			}).Build()
			mockey.Mock(qrcode.Encode).To(
				func(content string, level qrcode.RecoveryLevel, size int) ([]byte, error) {
					return []byte("png"), tc.mockQRCodeErr
				}).Build()

			mockey.Mock(NewUserService).To(func(_ context.Context) *UserService {
				return &UserService{}
			}).Build()

			secret, rawQRCode, err := NewUserService(t.Context()).GetMFA(1)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, "secret", secret)
			assert.Contains(t, rawQRCode, "data:image/png;base64,")
		})
	}
}

func TestUserService_BindMFA(t *testing.T) {
	type testCase struct {
		req           *user.BindMFAReq
		mockValid     bool
		mockValidErr  error
		mockUpdateErr error
		expectError   string
	}

	testCases := map[string]testCase{
		"success": {
			req:       &user.BindMFAReq{Secret: "secret", Code: "123456"},
			mockValid: true,
		},
		"validate error": {
			req:          &user.BindMFAReq{Secret: "secret", Code: "123456"},
			mockValidErr: assert.AnError,
			expectError:  assert.AnError.Error(),
		},
		"invalid code": {
			req:         &user.BindMFAReq{Secret: "secret", Code: "123456"},
			mockValid:   false,
			expectError: errno.MFACodeInvalidErr.ErrMsg,
		},
		"update error": {
			req:           &user.BindMFAReq{Secret: "secret", Code: "123456"},
			mockValid:     true,
			mockUpdateErr: assert.AnError,
			expectError:   assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(totpTool.ValidateCode).To(
				func(secret, code string) (bool, error) {
					return tc.mockValid, tc.mockValidErr
				}).Build()
			mockey.Mock((*userDao.UserDao).UpdateUserMFA).To(
				func(_ *userDao.UserDao, ctx context.Context, userID int64, secret string) error {
					return tc.mockUpdateErr
				}).Build()

			mockey.Mock(NewUserService).To(func(_ context.Context) *UserService {
				return &UserService{}
			}).Build()

			err := NewUserService(t.Context()).BindMFA(tc.req, 1)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
		})
	}
}
