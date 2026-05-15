package service

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	videoCache "github.com/ACaiCat/tiktok-go/pkg/cache/video"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	videoDao "github.com/ACaiCat/tiktok-go/pkg/db/video"
)

func TestVideoService_getVideosByIDs(t *testing.T) {
	type testCase struct {
		videoIDs          []int64
		cacheEnabled      bool
		mockCacheVideos   []*modelDao.Video
		mockCacheErr      error
		mockDBVideos      []*modelDao.Video
		mockDBErr         error
		expectVideos      []*modelDao.Video
		expectError       string
		expectCacheCalls  int
		expectDBCalls     int
		expectSetCalls    int
		expectDBVideoIDs  [][]int64
		expectSetVideoIDs [][]int64
	}

	video1 := &modelDao.Video{ID: 1, UserID: 10, Title: "video 1"}
	video2 := &modelDao.Video{ID: 2, UserID: 20, Title: "video 2"}
	video3 := &modelDao.Video{ID: 3, UserID: 30, Title: "video 3"}

	testCases := map[string]testCase{
		"empty ids returns empty without dependencies": {
			videoIDs:     []int64{},
			cacheEnabled: true,
			expectVideos: []*modelDao.Video{},
		},
		"nil cache reads db directly": {
			videoIDs:         []int64{1, 2},
			mockCacheErr:     assert.AnError,
			mockDBVideos:     []*modelDao.Video{video1, video2},
			expectVideos:     []*modelDao.Video{video1, video2},
			expectCacheCalls: 1,
			expectDBCalls:    1,
			expectSetCalls:   1,
			expectDBVideoIDs: [][]int64{{1, 2}},
			expectSetVideoIDs: [][]int64{
				{1, 2},
			},
		},
		"cache hit returns cached videos": {
			videoIDs:         []int64{1, 2},
			cacheEnabled:     true,
			mockCacheVideos:  []*modelDao.Video{video1, video2},
			expectVideos:     []*modelDao.Video{video1, video2},
			expectCacheCalls: 1,
		},
		"cache error falls back to db and warms cache": {
			videoIDs:          []int64{1, 2},
			cacheEnabled:      true,
			mockCacheErr:      assert.AnError,
			mockDBVideos:      []*modelDao.Video{video1, video2},
			expectVideos:      []*modelDao.Video{video1, video2},
			expectCacheCalls:  1,
			expectDBCalls:     1,
			expectSetCalls:    1,
			expectDBVideoIDs:  [][]int64{{1, 2}},
			expectSetVideoIDs: [][]int64{{1, 2}},
		},
		"cache error returns wrapped db error": {
			videoIDs:         []int64{1, 2},
			cacheEnabled:     true,
			mockCacheErr:     assert.AnError,
			mockDBErr:        assert.AnError,
			expectError:      "service.getVideosByIDs: db.GetVideosByIDs failed",
			expectCacheCalls: 1,
			expectDBCalls:    1,
			expectDBVideoIDs: [][]int64{{1, 2}},
		},
		"partial cache miss fetches missing ids and preserves requested order": {
			videoIDs:          []int64{1, 2, 3},
			cacheEnabled:      true,
			mockCacheVideos:   []*modelDao.Video{video1, nil, nil},
			mockDBVideos:      []*modelDao.Video{video3},
			expectVideos:      []*modelDao.Video{video1, video3},
			expectCacheCalls:  1,
			expectDBCalls:     1,
			expectSetCalls:    1,
			expectDBVideoIDs:  [][]int64{{2, 3}},
			expectSetVideoIDs: [][]int64{{3}},
		},
		"partial cache miss returns wrapped fallback db error": {
			videoIDs:         []int64{1, 2},
			cacheEnabled:     true,
			mockCacheVideos:  []*modelDao.Video{video1, nil},
			mockDBErr:        assert.AnError,
			expectError:      "service.getVideosByIDs: db.GetVideosByIDs fallback failed",
			expectCacheCalls: 1,
			expectDBCalls:    1,
			expectDBVideoIDs: [][]int64{{2}},
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			cacheCalls := 0
			dbCalls := 0
			setCalls := 0
			var dbVideoIDs [][]int64
			var setVideoIDs [][]int64

			mockey.Mock((*videoCache.VideoCache).GetVideos).To(
				func(_ *videoCache.VideoCache, ctx context.Context, videoIDs []int64) ([]*modelDao.Video, error) {
					cacheCalls++
					return tc.mockCacheVideos, tc.mockCacheErr
				}).Build()
			mockey.Mock((*videoCache.VideoCache).SetVideos).To(
				func(_ *videoCache.VideoCache, ctx context.Context, videos []*modelDao.Video) error {
					setCalls++
					ids := make([]int64, 0, len(videos))
					for _, video := range videos {
						ids = append(ids, video.ID)
					}
					setVideoIDs = append(setVideoIDs, ids)
					return nil
				}).Build()
			mockey.Mock((*videoDao.VideoDao).GetVideosByIDs).To(
				func(_ *videoDao.VideoDao, ctx context.Context, videoIDs []int64) ([]*modelDao.Video, error) {
					dbCalls++
					ids := append([]int64(nil), videoIDs...)
					dbVideoIDs = append(dbVideoIDs, ids)
					return tc.mockDBVideos, tc.mockDBErr
				}).Build()

			svc := &VideoService{
				videoDao: &videoDao.VideoDao{},
				ctx:      context.Background(),
			}
			if tc.cacheEnabled {
				svc.videoCache = &videoCache.VideoCache{}
			}

			result, err := svc.getVideosByIDs(tc.videoIDs)

			if tc.expectError != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectVideos, result)
			}

			assert.Equal(t, tc.expectCacheCalls, cacheCalls)
			assert.Equal(t, tc.expectDBCalls, dbCalls)
			assert.Equal(t, tc.expectSetCalls, setCalls)
			assert.Equal(t, tc.expectDBVideoIDs, dbVideoIDs)
			assert.Equal(t, tc.expectSetVideoIDs, setVideoIDs)
		})
	}
}
