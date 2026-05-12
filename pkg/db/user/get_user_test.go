package userdao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func TestUserDao_GetByID(t *testing.T) {
	type testCase struct {
		userID  int64
		mockRet *model.User
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"get user by id success":     {userID: 1, mockRet: &model.User{ID: 1, Username: "alice"}},
		"user not found returns nil": {userID: 99, mockRet: nil, mockErr: gorm.ErrRecordNotFound},
		"db error returns error":     {userID: 1, mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockUserQueryChain()
			dao := newTestDao()

			mockey.Mock((*gen.DO).First).Return(tc.mockRet, tc.mockErr).Build()

			u, err := dao.GetByID(context.Background(), tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetUserByID failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, u)
			}
		})
	}
}

func TestUserDao_GetByUsername(t *testing.T) {
	type testCase struct {
		username string
		mockRet  *model.User
		mockErr  error
		wantErr  bool
	}

	testCases := map[string]testCase{
		"get user by username success": {username: "alice", mockRet: &model.User{ID: 1, Username: "alice"}},
		"user not found returns nil":   {username: "ghost", mockRet: nil, mockErr: gorm.ErrRecordNotFound},
		"db error returns error":       {username: "alice", mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockUserQueryChain()
			dao := newTestDao()

			mockey.Mock((*gen.DO).First).Return(tc.mockRet, tc.mockErr).Build()

			u, err := dao.GetByUsername(context.Background(), tc.username)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetUserByUsername failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, u)
			}
		})
	}
}
