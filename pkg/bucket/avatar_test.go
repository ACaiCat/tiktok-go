package bucket

import (
	"context"
	"io"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func TestUploadAvatar(t *testing.T) {
	type testCase struct {
		userID  int64
		data    []byte
		putErr  error
		wantErr bool
	}

	testCases := map[string]testCase{
		"upload avatar success": {
			userID:  1,
			data:    []byte("fake-jpeg"),
			putErr:  nil,
			wantErr: false,
		},
		"upload avatar minio error": {
			userID:  2,
			data:    []byte("fake-jpeg"),
			putErr:  assert.AnError,
			wantErr: true,
		},
		"upload empty data": {
			userID:  3,
			data:    []byte{},
			putErr:  nil,
			wantErr: false,
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*minio.Client).PutObject).To(func(_ context.Context, _ string, _ string, _ io.Reader, _ int64, _ minio.PutObjectOptions) (minio.UploadInfo, error) {
				return minio.UploadInfo{}, tc.putErr
			}).Build()

			err := UploadAvatar(context.Background(), tc.userID, tc.data)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestGetAvatarURL(t *testing.T) {
	type testCase struct {
		userID      int64
		useSSL      bool
		wantPrefix  string
		wantContain string
	}

	testCases := map[string]testCase{
		"http url for userID 1": {
			userID:      1,
			useSSL:      false,
			wantPrefix:  "http://",
			wantContain: constants.AvatarBucketName + "/avatar_1",
		},
		"https url when SSL enabled": {
			userID:      42,
			useSSL:      true,
			wantPrefix:  "https://",
			wantContain: constants.AvatarBucketName + "/avatar_42",
		},
	}

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			setSSL(tc.useSSL)
			url := GetAvatarURL(tc.userID)
			assert.Contains(t, url, tc.wantPrefix)
			assert.Contains(t, url, tc.wantContain)
		})
	}
}
