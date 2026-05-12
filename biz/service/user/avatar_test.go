package service

import (
	"context"
	"mime/multipart"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/bucket"
	userDao "github.com/ACaiCat/tiktok-go/pkg/db/user"
	"github.com/ACaiCat/tiktok-go/pkg/img"
	"github.com/ACaiCat/tiktok-go/pkg/utils"
)

func TestUserService_UploadAvatar(t *testing.T) {
	type testCase struct {
		mockReadErr   error
		mockCheckErr  error
		mockUploadErr error
		mockUpdateErr error
		expectError   string
	}

	testCases := map[string]testCase{
		"success": {},
		"read file error": {
			mockReadErr: assert.AnError,
			expectError: assert.AnError.Error(),
		},
		"check avatar error": {
			mockCheckErr: assert.AnError,
			expectError:  assert.AnError.Error(),
		},
		"upload error": {
			mockUploadErr: assert.AnError,
			expectError:   assert.AnError.Error(),
		},
		"update url error": {
			mockUpdateErr: assert.AnError,
			expectError:   assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock(utils.FileHeaderToBytes).To(
				func(fileHeader *multipart.FileHeader) ([]byte, error) {
					return []byte("mock-image"), tc.mockReadErr
				}).Build()
			mockey.Mock(img.CheckAvatar).To(
				func(data []byte) (string, error) {
					return "jpeg", tc.mockCheckErr
				}).Build()
			mockey.Mock(bucket.UploadAvatar).To(
				func(ctx context.Context, userID int64, data []byte) error {
					return tc.mockUploadErr
				}).Build()
			mockey.Mock(bucket.GetAvatarURL).To(
				func(userID int64) string {
					return "https://example.com/avatar.jpg"
				}).Build()
			mockey.Mock((*userDao.UserDao).UpdateUserAvatarURL).To(
				func(_ *userDao.UserDao, ctx context.Context, userID int64, avatarURL string) error {
					return tc.mockUpdateErr
				}).Build()

			mockey.Mock(NewUserService).To(func(_ context.Context) *UserService {
				return &UserService{}
			}).Build()

			fh := &multipart.FileHeader{Filename: "avatar.jpg"}
			err := NewUserService(t.Context()).UploadAvatar(fh, 1)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
		})
	}
}
