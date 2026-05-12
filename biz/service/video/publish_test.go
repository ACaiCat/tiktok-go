package service

import (
	"context"
	"mime/multipart"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	videoDao "github.com/ACaiCat/tiktok-go/pkg/db/video"
	"github.com/ACaiCat/tiktok-go/pkg/ffmpeg"
	"github.com/ACaiCat/tiktok-go/pkg/utils"
)

func TestVideoService_PublishVideo(t *testing.T) {
	type testCase struct {
		mockReadErr      error
		mockTranscodeErr error
		mockPublishErr   error
		expectError      string
	}

	testCases := map[string]testCase{
		"success": {},
		"read file error": {
			mockReadErr: assert.AnError,
			expectError: assert.AnError.Error(),
		},
		"transcode error": {
			mockTranscodeErr: assert.AnError,
			expectError:      assert.AnError.Error(),
		},
		"publish db error": {
			mockPublishErr: assert.AnError,
			expectError:    assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(utils.FileHeaderToBytes).To(
				func(fileHeader *multipart.FileHeader) ([]byte, error) {
					return []byte("mock-data"), tc.mockReadErr
				}).Build()
			mockey.Mock(ffmpeg.TranscodeVideo).To(
				func(path string) ([]byte, error) {
					return []byte("transcoded"), tc.mockTranscodeErr
				}).Build()
			mockey.Mock(ffmpeg.GetVideoCover).To(
				func(path string) ([]byte, error) {
					return []byte("cover"), nil
				}).Build()
			mockey.Mock((*videoDao.VideoDao).PublishVideo).To(
				func(_ *videoDao.VideoDao, ctx context.Context, userID int64, title string, description string,
					uploadFn func(videoID int64) error,
					getVideoURL func(videoID int64) string,
					getCoverURL func(videoID int64) string,
				) error {
					return tc.mockPublishErr
				}).Build()

			mockey.Mock(NewVideoService).To(func(_ context.Context) *VideoService {
				return &VideoService{}
			}).Build()

			fh := &multipart.FileHeader{Filename: "test.mp4"}
			err := NewVideoService(context.Background()).PublishVideo(1, "title", "desc", fh)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
		})
	}
}
