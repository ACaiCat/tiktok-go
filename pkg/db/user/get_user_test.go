package userdao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

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
		"user not found returns nil": {userID: 99, mockRet: nil},
		"db error returns error":     {userID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*UserDao).GetByID).Return(tc.mockRet, tc.mockErr).Build()

			u, err := dao.GetByID(context.Background(), tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
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
		"user not found returns nil":   {username: "ghost", mockRet: nil},
		"db error returns error":       {username: "alice", mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*UserDao).GetByUsername).Return(tc.mockRet, tc.mockErr).Build()

			u, err := dao.GetByUsername(context.Background(), tc.username)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, u)
			}
		})
	}
}
