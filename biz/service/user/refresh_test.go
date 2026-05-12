package service

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/biz/model/user"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/jwt"
)

func TestUserService_RefreshToken(t *testing.T) {
	type testCase struct {
		req             *user.RefreshReq
		mockValidateErr error
		mockCreateErr   error
		expectError     string
	}

	testCases := map[string]testCase{
		"success": {
			req: &user.RefreshReq{RefreshToken: "valid-refresh-token"},
		},
		"invalid token": {
			req:             &user.RefreshReq{RefreshToken: "invalid-token"},
			mockValidateErr: assert.AnError,
			expectError:     assert.AnError.Error(),
		},
		"create token error": {
			req:           &user.RefreshReq{RefreshToken: "valid-refresh-token"},
			mockCreateErr: assert.AnError,
			expectError:   assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(jwt.ValidateToken).To(
				func(token string, tokenType int8) (int64, error) {
					return 1, tc.mockValidateErr
				}).Build()
			mockey.Mock(jwt.CreateToken).To(
				func(tokenType int8, userID int64) (string, error) {
					return "new-token", tc.mockCreateErr
				}).Build()

			mockey.Mock(NewUserService).To(func(_ context.Context) *UserService {
				return &UserService{}
			}).Build()

			_ = constants.TypeRefreshToken // ensure import used
			access, refresh, err := NewUserService(t.Context()).RefreshToken(tc.req)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, access)
			assert.NotEmpty(t, refresh)
		})
	}
}
