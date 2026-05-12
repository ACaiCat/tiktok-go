package service

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/ACaiCat/tiktok-go/biz/model/user"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	userDao "github.com/ACaiCat/tiktok-go/pkg/db/user"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/ACaiCat/tiktok-go/pkg/jwt"
	"github.com/ACaiCat/tiktok-go/pkg/totp"
)

func TestUserService_UserLogin(t *testing.T) {
	type testCase struct {
		req         *user.LoginReq
		mockUser    *modelDao.User
		mockErr     error
		expectError string
	}

	hashedPass := "$2a$10$XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	testCases := map[string]testCase{
		"user not found": {
			req:         &user.LoginReq{Username: "nouser", Password: "pass"},
			mockUser:    nil,
			expectError: errno.PasswordIsNotVerified.ErrMsg,
		},
		"db error": {
			req:         &user.LoginReq{Username: "testuser", Password: "pass"},
			mockErr:     assert.AnError,
			expectError: assert.AnError.Error(),
		},
		"mfa required but missing": {
			req: &user.LoginReq{Username: "mfauser", Password: "password123"},
			mockUser: &modelDao.User{
				ID:         1,
				Username:   "mfauser",
				Password:   hashedPass,
				TotpSecret: new("mysecret"),
			},
			expectError: errno.MFAMissingErr.ErrMsg,
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*userDao.UserDao).GetByUsername).To(
				func(_ *userDao.UserDao, ctx context.Context, username string) (*modelDao.User, error) {
					return tc.mockUser, tc.mockErr
				}).Build()
			mockey.Mock(bcrypt.CompareHashAndPassword).To(
				func(hashedPassword, password []byte) error {
					return nil
				}).Build()
			mockey.Mock(totp.ValidateCode).To(
				func(secret, code string) (bool, error) {
					return true, nil
				}).Build()
			mockey.Mock(jwt.CreateToken).To(
				func(tokenType int8, userID int64) (string, error) {
					return "mock-token", nil
				}).Build()

			mockey.Mock(NewUserService).To(func(_ context.Context) *UserService {
				return &UserService{}
			}).Build()

			_, _, _, err := NewUserService(context.Background()).UserLogin(tc.req)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
		})
	}
}
