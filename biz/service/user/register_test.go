package service

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/biz/model/user"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	userDao "github.com/ACaiCat/tiktok-go/pkg/db/user"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func TestUserService_UserRegister(t *testing.T) {
	type testCase struct {
		req         *user.RegisterReq
		mockUser    *modelDao.User
		mockErr     error
		expectError string
	}

	testCases := map[string]testCase{
		"success": {
			req:      &user.RegisterReq{Username: "validuser", Password: "validpass"},
			mockUser: nil,
		},
		"username too short": {
			req:         &user.RegisterReq{Username: "ab", Password: "validpass"},
			expectError: "username must be at least",
		},
		"username too long": {
			req:         &user.RegisterReq{Username: "thisusernameiswaytoolongandexceedsthirtytwocharacters", Password: "validpass"},
			expectError: "username must be at most",
		},
		"password too short": {
			req:         &user.RegisterReq{Username: "validuser", Password: "ab"},
			expectError: "password must be at least",
		},
		"user already exists": {
			req:         &user.RegisterReq{Username: "existuser", Password: "validpass"},
			mockUser:    &modelDao.User{ID: 1, Username: "existuser"},
			expectError: errno.UserAlreadyExistErr.ErrMsg,
		},
		"db get user error": {
			req:         &user.RegisterReq{Username: "validuser", Password: "validpass"},
			mockErr:     assert.AnError,
			expectError: assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*userDao.UserDao).GetByUsername).To(
				func(_ *userDao.UserDao, ctx context.Context, username string) (*modelDao.User, error) {
					return tc.mockUser, tc.mockErr
				}).Build()
			mockey.Mock((*userDao.UserDao).CreateUser).To(
				func(_ *userDao.UserDao, ctx context.Context, username string, password string) (int64, error) {
					return 1, nil
				}).Build()

			mockey.Mock(NewUserService).To(func(_ context.Context) *UserService {
				return &UserService{}
			}).Build()

			err := NewUserService(context.Background()).UserRegister(tc.req)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
		})
	}
}
