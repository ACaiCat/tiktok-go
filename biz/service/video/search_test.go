package service

import (
	"context"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/biz/model/video"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	videoDao "github.com/ACaiCat/tiktok-go/pkg/db/video"
)

func TestVideoService_SearchVideo(t *testing.T) {
	type testCase struct {
		req         *video.SearchReq
		mockVideos  []*modelDao.Video
		mockErr     error
		expectError string
	}

	testCases := map[string]testCase{
		"success": {
			req:        &video.SearchReq{Keywords: "golang test", PageNum: 0, PageSize: 10},
			mockVideos: []*modelDao.Video{{ID: 1, Title: "golang test video"}},
		},
		"success with dates": {
			req:        &video.SearchReq{Keywords: "test", PageNum: 0, PageSize: 10, FromDate: new("1778583551826"), ToDate: new("1778583551826")},
			mockVideos: []*modelDao.Video{},
		},
		"invalid from date": {
			req:         &video.SearchReq{Keywords: "test", PageNum: 0, PageSize: 10, FromDate: new("not-a-date")},
			expectError: "参数错误",
		},
		"invalid to date": {
			req:         &video.SearchReq{Keywords: "test", PageNum: 0, PageSize: 10, ToDate: new("not-a-date")},
			expectError: "参数错误",
		},
		"db error": {
			req:         &video.SearchReq{Keywords: "golang", PageNum: 0, PageSize: 10},
			mockErr:     assert.AnError,
			expectError: assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*videoDao.VideoDao).SearchVideo).To(
				func(_ *videoDao.VideoDao, ctx context.Context, keywords []string, pageSize int, pageNum int,
					fromDate time.Time, toDate time.Time, username string) ([]*modelDao.Video, error) {
					return tc.mockVideos, tc.mockErr
				}).Build()

			mockey.Mock(NewVideoService).To(func(_ context.Context) *VideoService {
				return &VideoService{}
			}).Build()

			result, err := NewVideoService(t.Context()).SearchVideo(tc.req)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, VideosDaoToDto(tc.mockVideos), result)
		})
	}
}
