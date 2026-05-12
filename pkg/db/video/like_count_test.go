package videodao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestVideoDao_IncrLikeCount(t *testing.T) {
	type testCase struct {
		videoID int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"incr like count success": {videoID: 1},
		"db error":                {videoID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*VideoDao).IncrLikeCount).Return(tc.mockErr).Build()

			err := dao.IncrLikeCount(context.Background(), tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVideoDao_DecrLikeCount(t *testing.T) {
	type testCase struct {
		videoID int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"decr like count success": {videoID: 1},
		"db error":                {videoID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*VideoDao).DecrLikeCount).Return(tc.mockErr).Build()

			err := dao.DecrLikeCount(context.Background(), tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
