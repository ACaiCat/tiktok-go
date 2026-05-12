package service

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	userDao "github.com/ACaiCat/tiktok-go/pkg/db/user"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func TestUserService_GetUserInfo(t *testing.T) {
	type testCase struct {
		mockUser    *modelDao.User
		mockErr     error
		expectError string
	}

	testCases := map[string]testCase{
		"success": {
			mockUser: &modelDao.User{ID: 1, Username: "testuser"},
		},
		"user not found": {
			mockUser:    nil,
			expectError: errno.UserIsNotExistErr.ErrMsg,
		},
		"db error": {
			mockErr:     assert.AnError,
			expectError: assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*userDao.UserDao).GetByID).To(
				func(_ *userDao.UserDao, ctx context.Context, id int64) (*modelDao.User, error) {
					return tc.mockUser, tc.mockErr
				}).Build()

			mockey.Mock(NewUserService).To(func(_ context.Context) *UserService {
				return &UserService{}
			}).Build()

			result, err := NewUserService(t.Context()).GetUserInfo(1)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, UserDaoToDto(tc.mockUser), result)
		})
	}
}
