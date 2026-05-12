package likedao

import (
	"context"
	"reflect"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
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
			mockLikeQueryChain()
			dao := newTestDao()
			dbtestutil.MockScan(func(dest interface{}) {
				destValue := reflect.ValueOf(dest)
				if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Slice {
					return
				}
				sliceValue := destValue.Elem()
				elemType := sliceValue.Type().Elem()
				for videoID, count := range tc.mockRet {
					item := reflect.New(elemType).Elem()
					item.FieldByName("VideoID").SetInt(videoID)
					item.FieldByName("Count").SetInt(count)
					sliceValue = reflect.Append(sliceValue, item)
				}
				destValue.Elem().Set(sliceValue)
			}, tc.mockErr)

			m, err := dao.GetLikeCounts(context.Background(), tc.videoIDs)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetLikeCount failed")
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
			mockLikeQueryChain()
			dao := newTestDao()
			dbtestutil.MockScan(func(dest interface{}) {
				dbtestutil.FillValue(dest, tc.mockRet)
			}, tc.mockErr)

			ids, err := dao.GetUserLikes(context.Background(), tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetUserLikes failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, ids)
			}
		})
	}
}
