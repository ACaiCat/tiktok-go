package service

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/biz/model/video"
	videoDao "github.com/ACaiCat/tiktok-go/pkg/db/video"
)

func TestVideoService_VisitVideo(t *testing.T) {
	type testCase struct {
		req         *video.VisitVideoReq
		mockErr     error
		expectError string
	}

	testCases := map[string]testCase{
		"success": {
			req: &video.VisitVideoReq{VideoID: "1"},
		},
		"invalid video id": {
			req:         &video.VisitVideoReq{VideoID: "not-a-number"},
			expectError: "参数错误",
		},
		"db error": {
			req:         &video.VisitVideoReq{VideoID: "1"},
			mockErr:     assert.AnError,
			expectError: assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*videoDao.VideoDao).IncrVisitCount).To(
				func(_ *videoDao.VideoDao, ctx context.Context, videoID int64) error {
					return tc.mockErr
				}).Build()

			mockey.Mock(NewVideoService).To(func(_ context.Context) *VideoService {
				return &VideoService{}
			}).Build()

			err := NewVideoService(context.Background()).VisitVideo(tc.req)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
		})
	}
}
