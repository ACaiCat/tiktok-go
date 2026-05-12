package followerdao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func TestFollowerDao_GetFollower(t *testing.T) {
	type testCase struct {
		userID   int64
		pageSize int
		pageNum  int
		mockRet  []*model.User
		mockCnt  int
		mockErr  error
		wantErr  bool
	}

	users := []*model.User{{ID: 2, Username: "bob"}}

	testCases := map[string]testCase{
		"get followers success": {userID: 1, pageSize: 10, pageNum: 0, mockRet: users, mockCnt: 1},
		"no followers":          {userID: 1, pageSize: 10, pageNum: 0, mockRet: []*model.User{}, mockCnt: 0},
		"db error":              {userID: 1, pageSize: 10, pageNum: 0, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*FollowerDao).GetFollower).Return(tc.mockRet, tc.mockCnt, tc.mockErr).Build()

			us, cnt, err := dao.GetFollower(context.Background(), tc.userID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, us)
				assert.Equal(t, tc.mockCnt, cnt)
			}
		})
	}
}

func TestFollowerDao_GetFollowing(t *testing.T) {
	type testCase struct {
		userID   int64
		pageSize int
		pageNum  int
		mockRet  []*model.User
		mockCnt  int
		mockErr  error
		wantErr  bool
	}

	users := []*model.User{{ID: 3, Username: "carol"}}

	testCases := map[string]testCase{
		"get following success": {userID: 1, pageSize: 10, pageNum: 0, mockRet: users, mockCnt: 1},
		"no following":          {userID: 1, pageSize: 10, pageNum: 0, mockRet: []*model.User{}, mockCnt: 0},
		"db error":              {userID: 1, pageSize: 10, pageNum: 0, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*FollowerDao).GetFollowing).Return(tc.mockRet, tc.mockCnt, tc.mockErr).Build()

			us, cnt, err := dao.GetFollowing(context.Background(), tc.userID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, us)
				assert.Equal(t, tc.mockCnt, cnt)
			}
		})
	}
}

func TestFollowerDao_GetFriends(t *testing.T) {
	type testCase struct {
		userID   int64
		pageSize int
		pageNum  int
		mockRet  []*model.User
		mockCnt  int
		mockErr  error
		wantErr  bool
	}

	users := []*model.User{{ID: 2, Username: "bob"}}

	testCases := map[string]testCase{
		"get friends success": {userID: 1, pageSize: 10, pageNum: 0, mockRet: users, mockCnt: 1},
		"no friends":          {userID: 1, pageSize: 10, pageNum: 0, mockRet: []*model.User{}, mockCnt: 0},
		"db error":            {userID: 1, pageSize: 10, pageNum: 0, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*FollowerDao).GetFriends).Return(tc.mockRet, tc.mockCnt, tc.mockErr).Build()

			us, cnt, err := dao.GetFriends(context.Background(), tc.userID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, us)
				assert.Equal(t, tc.mockCnt, cnt)
			}
		})
	}
}
