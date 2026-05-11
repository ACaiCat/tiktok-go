package followerdao

import (
	"context"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestIsExistFollow(t *testing.T) {
	type testCase struct {
		userID     int64
		followerID int64
		mockRet    bool
		mockErr    error
		wantErr    bool
	}

	testCases := map[string]testCase{
		"follow exists":    {userID: 1, followerID: 2, mockRet: true},
		"follow not exist": {userID: 1, followerID: 2, mockRet: false},
		"db error":         {userID: 1, followerID: 2, mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*FollowerDao).IsExistFollow).Return(tc.mockRet, tc.mockErr).Build()

			ok, err := dao.IsExistFollow(context.Background(), tc.userID, tc.followerID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, ok)
			}
		})
	}
}

func TestIsExistFriend(t *testing.T) {
	type testCase struct {
		userID   int64
		friendID int64
		mockRet  bool
		mockErr  error
		wantErr  bool
	}

	testCases := map[string]testCase{
		"friend exists":    {userID: 1, friendID: 2, mockRet: true},
		"friend not exist": {userID: 1, friendID: 2, mockRet: false},
		"db error":         {userID: 1, friendID: 2, mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*FollowerDao).IsExistFriend).Return(tc.mockRet, tc.mockErr).Build()

			ok, err := dao.IsExistFriend(context.Background(), tc.userID, tc.friendID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, ok)
			}
		})
	}
}
