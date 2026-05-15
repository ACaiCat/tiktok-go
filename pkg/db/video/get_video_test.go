package videodao

import (
	"context"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func TestVideoDao_GetVideoByID(t *testing.T) {
	type testCase struct {
		videoID int64
		mockRet *model.Video
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"get video success":           {videoID: 1, mockRet: &model.Video{ID: 1, Title: "test"}},
		"video not found returns nil": {videoID: 99, mockRet: nil, mockErr: gorm.ErrRecordNotFound},
		"db error returns error":      {videoID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockVideoQueryChain()
			dao := newTestDao()
			dbtestutil.MockFirst(tc.mockRet, tc.mockErr)

			v, err := dao.GetVideoByID(context.Background(), tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetVideoByID failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, v)
			}
		})
	}
}

func TestVideoDao_GetVideosByIDs(t *testing.T) {
	type testCase struct {
		videoIDs []int64
		mockRet  []*model.Video
		mockErr  error
		wantRet  []*model.Video
		wantErr  bool
	}

	videos := []*model.Video{{ID: 2}, {ID: 1}}

	testCases := map[string]testCase{
		"get videos success keeps input order": {
			videoIDs: []int64{1, 2},
			mockRet:  videos,
			wantRet:  []*model.Video{{ID: 1}, {ID: 2}},
		},
		"empty ids": {
			videoIDs: []int64{},
			wantRet:  []*model.Video{},
		},
		"db error returns error": {
			videoIDs: []int64{1},
			mockErr:  assert.AnError,
			wantErr:  true,
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockVideoQueryChain()
			dao := newTestDao()
			dbtestutil.MockFind(tc.mockRet, tc.mockErr)

			vs, err := dao.GetVideosByIDs(context.Background(), tc.videoIDs)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetVideosByIDs failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantRet, vs)
			}
		})
	}
}

func TestVideoDao_GetFeedByLatestTime(t *testing.T) {
	type testCase struct {
		latestTime time.Time
		limit      int
		mockRet    []*model.Video
		mockErr    error
		wantErr    bool
	}

	videos := []*model.Video{{ID: 1}, {ID: 2}}

	testCases := map[string]testCase{
		"get feed success":       {latestTime: time.Time{}, limit: 10, mockRet: videos},
		"with latest time":       {latestTime: time.Now(), limit: 5, mockRet: videos[:1]},
		"db error returns error": {limit: 10, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockVideoQueryChain()
			dao := newTestDao()
			dbtestutil.MockFind(tc.mockRet, tc.mockErr)

			vs, err := dao.GetFeedByLatestTime(context.Background(), tc.latestTime, tc.limit)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetFeedByLatestTime failed")
			} else {
				assert.NoError(t, err)
				assert.Len(t, vs, len(tc.mockRet))
			}
		})
	}
}

func TestVideoDao_GetVideosByUserID(t *testing.T) {
	type testCase struct {
		userID   int64
		pageSize int
		pageNum  int
		mockRet  []*model.Video
		mockErr  error
		wantErr  bool
	}

	videos := []*model.Video{{ID: 1}, {ID: 2}}

	testCases := map[string]testCase{
		"get videos by user success": {userID: 1, pageSize: 10, pageNum: 0, mockRet: videos},
		"no videos":                  {userID: 1, pageSize: 10, pageNum: 0, mockRet: []*model.Video{}},
		"db error":                   {userID: 1, pageSize: 10, pageNum: 0, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockVideoQueryChain()
			dao := newTestDao()
			dbtestutil.MockFind(tc.mockRet, tc.mockErr)

			vs, err := dao.GetVideosByUserID(context.Background(), tc.userID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetVideosByUserID failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, vs)
			}
		})
	}
}

func TestVideoDao_GetPopularVideos(t *testing.T) {
	type testCase struct {
		pageSize int
		pageNum  int
		mockRet  []*model.Video
		mockErr  error
		wantErr  bool
	}

	videos := []*model.Video{{ID: 1, VisitCount: 100}}

	testCases := map[string]testCase{
		"get popular videos success": {pageSize: 10, pageNum: 0, mockRet: videos},
		"empty result":               {pageSize: 10, pageNum: 0, mockRet: []*model.Video{}},
		"db error":                   {pageSize: 10, pageNum: 0, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockVideoQueryChain()
			dao := newTestDao()
			dbtestutil.MockFind(tc.mockRet, tc.mockErr)

			vs, err := dao.GetPopularVideos(context.Background(), tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetPopularVideos failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, vs)
			}
		})
	}
}

func TestVideoDao_GetVideoCountByUserID(t *testing.T) {
	type testCase struct {
		userID  int64
		mockRet int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"get count success":      {userID: 1, mockRet: 5},
		"db error returns error": {userID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockVideoQueryChain()
			dao := newTestDao()
			dbtestutil.MockCount(tc.mockRet, tc.mockErr)

			cnt, err := dao.GetVideoCountByUserID(context.Background(), tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetVideoCountByUserID failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, cnt)
			}
		})
	}
}

func TestVideoDao_GetUserLikeList(t *testing.T) {
	type testCase struct {
		userID   int64
		pageSize int
		pageNum  int
		mockRet  []*model.Video
		mockErr  error
		wantErr  bool
	}

	videos := []*model.Video{{ID: 1}, {ID: 3}}

	testCases := map[string]testCase{
		"get like list success": {userID: 1, pageSize: 10, pageNum: 0, mockRet: videos},
		"empty list":            {userID: 1, pageSize: 10, pageNum: 0, mockRet: []*model.Video{}},
		"db error":              {userID: 1, pageSize: 10, pageNum: 0, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockVideoQueryChain()
			dao := newTestDao()
			scanErr := error(nil)
			findErr := tc.mockErr
			if tc.wantErr {
				scanErr = tc.mockErr
				findErr = nil
			}
			dbtestutil.MockScan(func(dest interface{}) {
				ids := make([]int64, 0, len(tc.mockRet))
				for _, video := range tc.mockRet {
					ids = append(ids, video.ID)
				}
				dbtestutil.FillValue(dest, ids)
			}, scanErr)
			dbtestutil.MockFind(tc.mockRet, findErr)

			vs, err := dao.GetUserLikeList(context.Background(), tc.userID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetUserLikeList failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, vs)
			}
		})
	}
}

func TestVideoDao_IsVideoExists(t *testing.T) {
	type testCase struct {
		videoID int64
		mockRet bool
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"video exists":    {videoID: 1, mockRet: true},
		"video not exist": {videoID: 99, mockRet: false},
		"db error":        {videoID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockVideoQueryChain()
			dao := newTestDao()
			count := int64(0)
			if tc.mockRet {
				count = 1
			}
			dbtestutil.MockCount(count, tc.mockErr)

			ok, err := dao.IsVideoExists(context.Background(), tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "IsVideoExists failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, ok)
			}
		})
	}
}
