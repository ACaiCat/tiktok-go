package likedao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestLikeDao_GetLikeCounts(t *testing.T) {
	type testCase struct {
		videoIDs []int64
		mockRet  map[int64]int64
		mockErr  error
		wantErr  bool
	}

	testCases := map[string]testCase{
		"get like counts success": {videoIDs: []int64{1, 2}, mockRet: map[int64]int64{1: 5, 2: 3}},
		"empty video ids":         {videoIDs: []int64{}, mockRet: map[int64]int64{}},
		"db error returns error":  {videoIDs: []int64{1}, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*LikeDao).GetLikeCounts).Return(tc.mockRet, tc.mockErr).Build()

			m, err := dao.GetLikeCounts(context.Background(), tc.videoIDs)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, m)
			}
		})
	}
}

func TestLikeDao_GetUserLikes(t *testing.T) {
	type testCase struct {
		userID  int64
		mockRet []int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"get user likes success": {userID: 1, mockRet: []int64{10, 20}},
		"no likes":               {userID: 1, mockRet: []int64{}},
		"db error returns error": {userID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*LikeDao).GetUserLikes).Return(tc.mockRet, tc.mockErr).Build()

			ids, err := dao.GetUserLikes(context.Background(), tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, ids)
			}
		})
	}
}
