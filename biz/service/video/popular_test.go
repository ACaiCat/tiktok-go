package service

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/biz/model/video"
	videoCache "github.com/ACaiCat/tiktok-go/pkg/cache/video"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	videoDao "github.com/ACaiCat/tiktok-go/pkg/db/video"
)

func TestVideoService_GetPopularVideos(t *testing.T) {
	type testCase struct {
		req             *video.PopularReq
		mockPopularIDs  []int64
		mockCacheVideos []*modelDao.Video
		mockPageVideos  []*modelDao.Video
		mockWarmVideos  []*modelDao.Video
		mockDBVideos    []*modelDao.Video
		mockCacheErr    error
		mockDBErr       error
		expectError     string
	}

	mockVideos := []*modelDao.Video{
		{ID: 1, UserID: 100, Title: "popular video"},
	}

	testCases := map[string]testCase{
		"success from cache": {
			req:             &video.PopularReq{PageNum: 0, PageSize: 10},
			mockPopularIDs:  []int64{1},
			mockCacheVideos: mockVideos,
		},
		"cache detail miss fallback to db": {
			req:             &video.PopularReq{PageNum: 0, PageSize: 10},
			mockPopularIDs:  []int64{1},
			mockCacheVideos: []*modelDao.Video{nil},
			mockDBVideos:    mockVideos,
		},
		"outside cache window reads db directly": {
			req:            &video.PopularReq{PageNum: 10, PageSize: 10},
			mockPageVideos: mockVideos,
		},
		"db error after cache miss": {
			req:          &video.PopularReq{PageNum: 0, PageSize: 10},
			mockCacheErr: assert.AnError,
			mockDBErr:    assert.AnError,
			expectError:  assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*videoCache.VideoCache).GetPopularVideos).To(
				func(_ *videoCache.VideoCache, ctx context.Context, pageSize int, pageNum int) ([]int64, error) {
					return tc.mockPopularIDs, tc.mockCacheErr
				}).Build()
			mockey.Mock((*videoCache.VideoCache).GetVideos).To(
				func(_ *videoCache.VideoCache, ctx context.Context, videoIDs []int64) ([]*modelDao.Video, error) {
					return tc.mockCacheVideos, nil
				}).Build()
			mockey.Mock((*videoCache.VideoCache).SetPopularVideos).Return(nil).Build()
			mockey.Mock((*videoCache.VideoCache).SetVideos).Return(nil).Build()
			mockey.Mock((*videoDao.VideoDao).GetPopularVideos).To(
				func(_ *videoDao.VideoDao, ctx context.Context, pageSize int, pageNum int) ([]*modelDao.Video, error) {
					if pageSize == 10 && pageNum == 0 && tc.mockPageVideos != nil {
						return tc.mockPageVideos, tc.mockDBErr
					}
					if pageSize == 10 && pageNum == 10 && tc.mockPageVideos != nil {
						return tc.mockPageVideos, tc.mockDBErr
					}
					if tc.mockWarmVideos != nil {
						return tc.mockWarmVideos, nil
					}
					return tc.mockDBVideos, tc.mockDBErr
				}).Build()
			mockey.Mock((*videoDao.VideoDao).GetVideosByIDs).To(
				func(_ *videoDao.VideoDao, ctx context.Context, videoIDs []int64) ([]*modelDao.Video, error) {
					return tc.mockDBVideos, tc.mockDBErr
				}).Build()

			mockey.Mock(NewVideoService).To(func(_ context.Context) *VideoService {
				return &VideoService{ctx: context.Background(), videoCache: &videoCache.VideoCache{}, videoDao: &videoDao.VideoDao{}}
			}).Build()

			result, err := NewVideoService(context.Background()).GetPopularVideos(tc.req)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}
