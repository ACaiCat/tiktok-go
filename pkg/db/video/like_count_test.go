package videodao

import (
	"context"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestIncrLikeCount(t *testing.T) {
	type testCase struct {
		videoID int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"incr like count success": {videoID: 1},
		"db error":                {videoID: 1, mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*VideoDao).IncrLikeCount).Return(tc.mockErr).Build()

			err := dao.IncrLikeCount(context.Background(), tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDecrLikeCount(t *testing.T) {
	type testCase struct {
		videoID int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"decr like count success": {videoID: 1},
		"db error":                {videoID: 1, mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*VideoDao).DecrLikeCount).Return(tc.mockErr).Build()

			err := dao.DecrLikeCount(context.Background(), tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
