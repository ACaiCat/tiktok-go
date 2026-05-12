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

func TestVideoService_GetFeed(t *testing.T) {
	type testCase struct {
		req         *video.FeedReq
		mockResult  []*modelDao.Video
		mockErr     error
		expectError string
	}

	testCases := map[string]testCase{
		"success": {
			req: &video.FeedReq{
				LatestTime: nil,
			},
			mockResult: []*modelDao.Video{
				{
					ID:           1,
					UserID:       100,
					VideoURL:     "https://www.example.com",
					CoverURL:     "https://www.example.com",
					Title:        "114514",
					Description:  "1919810",
					VisitCount:   114,
					LikeCount:    514,
					CommentCount: 1919,
				},
			},
		},
		"with latest time": {
			req: &video.FeedReq{
				LatestTime: new("1778583551826"),
			},
			mockResult: []*modelDao.Video{
				{
					ID:           1,
					UserID:       100,
					VideoURL:     "https://www.example.com",
					CoverURL:     "https://www.example.com",
					Title:        "114514",
					Description:  "1919810",
					VisitCount:   114,
					LikeCount:    514,
					CommentCount: 1919,
				},
			},
		},
		"invalid latest time": {
			req: &video.FeedReq{
				LatestTime: new("what can i say"),
			},
			mockResult:  []*modelDao.Video{},
			expectError: `参数错误: strconv.ParseInt: parsing "what can i say": invalid syntax`,
		},
		"video dao error": {
			req:         &video.FeedReq{},
			mockResult:  []*modelDao.Video{},
			mockErr:     assert.AnError,
			expectError: assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*videoDao.VideoDao).GetFeedByLatestTime).To(
				func(ctx context.Context, latestTime time.Time, limit int) ([]*modelDao.Video, error) {
					return tc.mockResult, tc.mockErr
				}).Build()

			mockey.Mock(NewVideoService).To(func(_ context.Context) *VideoService {
				return &VideoService{}
			}).Build()

			result, err := NewVideoService(t.Context()).GetFeed(tc.req)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, VideosDaoToDto(tc.mockResult), result)
		})
	}
}
