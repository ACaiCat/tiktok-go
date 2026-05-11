package bucket

import (
	"context"
	"io"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func TestUploadCover(t *testing.T) {
	type testCase struct {
		videoID int64
		data    []byte
		putErr  error
		wantErr bool
	}

	testCases := map[string]testCase{
		"upload cover success": {
			videoID: 1,
			data:    []byte("fake-jpeg"),
			putErr:  nil,
			wantErr: false,
		},
		"upload cover minio error": {
			videoID: 2,
			data:    []byte("fake-jpeg"),
			putErr:  assert.AnError,
			wantErr: true,
		},
		"upload empty data": {
			videoID: 3,
			data:    []byte{},
			putErr:  nil,
			wantErr: false,
		},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			Mock((*minio.Client).PutObject).To(func(_ context.Context, _ string, _ string, _ io.Reader, _ int64, _ minio.PutObjectOptions) (minio.UploadInfo, error) {
				return minio.UploadInfo{}, tc.putErr
			}).Build()

			err := UploadCover(context.Background(), tc.videoID, tc.data)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestGetCoverURL(t *testing.T) {
	type testCase struct {
		videoID     int64
		useSSL      bool
		wantPrefix  string
		wantContain string
	}

	testCases := map[string]testCase{
		"http url for videoID 1": {
			videoID:     1,
			useSSL:      false,
			wantPrefix:  "http://",
			wantContain: constants.CoverBucketName + "/cover_1",
		},
		"https url when SSL enabled": {
			videoID:     7,
			useSSL:      true,
			wantPrefix:  "https://",
			wantContain: constants.CoverBucketName + "/cover_7",
		},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			setSSL(tc.useSSL)
			url := GetCoverURL(tc.videoID)
			assert.Contains(t, url, tc.wantPrefix)
			assert.Contains(t, url, tc.wantContain)
		})
	}
}
