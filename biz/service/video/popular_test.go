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
		req         *video.PopularReq
		mockVideos  []*modelDao.Video
		mockErr     error
		expectError string
	}

	mockVideos := []*modelDao.Video{
		{ID: 1, UserID: 100, Title: "popular video"},
	}

	testCases := map[string]testCase{
		"success from cache": {
			req:        &video.PopularReq{PageNum: 0, PageSize: 10},
			mockVideos: mockVideos,
		},
		"db error": {
			req:         &video.PopularReq{PageNum: 0, PageSize: 10},
			mockErr:     assert.AnError,
			expectError: assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*videoCache.VideoCache).GetPopularVideos).To(
				func(_ *videoCache.VideoCache, ctx context.Context) ([]*modelDao.Video, error) {
					return tc.mockVideos, tc.mockErr
				}).Build()
			mockey.Mock((*videoDao.VideoDao).GetPopularVideos).To(
				func(_ *videoDao.VideoDao, ctx context.Context, pageSize int, pageNum int) ([]*modelDao.Video, error) {
					return tc.mockVideos, tc.mockErr
				}).Build()

			mockey.Mock(NewVideoService).To(func(_ context.Context) *VideoService {
				return &VideoService{}
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
