package videodao

import (
	"context"
	"testing"
	"time"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func TestGetVideoByID(t *testing.T) {
	type testCase struct {
		videoID int64
		mockRet *model.Video
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"get video success":           {videoID: 1, mockRet: &model.Video{ID: 1, Title: "test"}},
		"video not found returns nil": {videoID: 99, mockRet: nil},
		"db error returns error":      {videoID: 1, mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*VideoDao).GetVideoByID).Return(tc.mockRet, tc.mockErr).Build()

			v, err := dao.GetVideoByID(context.Background(), tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, v)
			}
		})
	}
}

func TestGetFeedByLatestTime(t *testing.T) {
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

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*VideoDao).GetFeedByLatestTime).Return(tc.mockRet, tc.mockErr).Build()

			vs, err := dao.GetFeedByLatestTime(context.Background(), tc.latestTime, tc.limit)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, vs)
			}
		})
	}
}

func TestGetVideosByUserID(t *testing.T) {
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

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*VideoDao).GetVideosByUserID).Return(tc.mockRet, tc.mockErr).Build()

			vs, err := dao.GetVideosByUserID(context.Background(), tc.userID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, vs)
			}
		})
	}
}

func TestGetPopularVideos(t *testing.T) {
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

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*VideoDao).GetPopularVideos).Return(tc.mockRet, tc.mockErr).Build()

			vs, err := dao.GetPopularVideos(context.Background(), tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, vs)
			}
		})
	}
}

func TestGetVideoCountByUserID(t *testing.T) {
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

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*VideoDao).GetVideoCountByUserID).Return(tc.mockRet, tc.mockErr).Build()

			cnt, err := dao.GetVideoCountByUserID(context.Background(), tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, cnt)
			}
		})
	}
}

func TestGetUserLikeList(t *testing.T) {
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

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*VideoDao).GetUserLikeList).Return(tc.mockRet, tc.mockErr).Build()

			vs, err := dao.GetUserLikeList(context.Background(), tc.userID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, vs)
			}
		})
	}
}
