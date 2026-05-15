package followerdao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
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
			mockFollowerQueryChain()
			dao := newTestDao()
			dbtestutil.MockTransaction(nil)
			dbtestutil.MockScan(func(dest interface{}) {
				ids := make([]int64, 0, tc.mockCnt)
				for _, user := range tc.mockRet {
					ids = append(ids, user.ID)
				}
				dbtestutil.FillValue(dest, ids)
			}, tc.mockErr)
			dbtestutil.MockFind(tc.mockRet, tc.mockErr)

			us, cnt, err := dao.GetFollower(context.Background(), tc.userID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetFollower failed")
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
			mockFollowerQueryChain()
			dao := newTestDao()
			dbtestutil.MockTransaction(nil)
			dbtestutil.MockScan(func(dest interface{}) {
				ids := make([]int64, 0, tc.mockCnt)
				for _, user := range tc.mockRet {
					ids = append(ids, user.ID)
				}
				dbtestutil.FillValue(dest, ids)
			}, tc.mockErr)
			dbtestutil.MockFind(tc.mockRet, tc.mockErr)

			us, cnt, err := dao.GetFollowing(context.Background(), tc.userID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetFollowing failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, us)
				assert.Equal(t, tc.mockCnt, cnt)
			}
		})
	}
}

func TestFollowerDao_GetFollowingIDs(t *testing.T) {
	type testCase struct {
		userID  int64
		mockIDs []int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"get following ids success": {userID: 1, mockIDs: []int64{2, 3}},
		"no following ids":          {userID: 1, mockIDs: []int64{}},
		"db error":                  {userID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockFollowerQueryChain()
			dao := newTestDao()
			dbtestutil.MockScan(func(dest interface{}) {
				dbtestutil.FillValue(dest, tc.mockIDs)
			}, tc.mockErr)

			ids, err := dao.GetFollowingIDs(context.Background(), tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetFollowingIDs failed")
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.mockIDs, ids)
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
			mockFollowerQueryChain()
			dao := newTestDao()
			dbtestutil.MockTransaction(nil)
			scanCall := 0
			dbtestutil.MockScan(func(dest interface{}) {
				scanCall++
				ids := make([]int64, 0, tc.mockCnt)
				for _, user := range tc.mockRet {
					ids = append(ids, user.ID)
				}
				if scanCall == 1 {
					dbtestutil.FillValue(dest, ids)
					return
				}
				dbtestutil.FillValue(dest, ids)
			}, tc.mockErr)
			dbtestutil.MockFind(tc.mockRet, tc.mockErr)

			us, cnt, err := dao.GetFriends(context.Background(), tc.userID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetFriend failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, us)
				assert.Equal(t, tc.mockCnt, cnt)
			}
		})
	}
}
