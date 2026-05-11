package videodao

import (
	"context"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestIsVideoExists(t *testing.T) {
	type testCase struct {
		videoID int64
		mockRet bool
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"video exists":    {videoID: 1, mockRet: true},
		"video not exist": {videoID: 99, mockRet: false},
		"db error":        {videoID: 1, mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*VideoDao).IsVideoExists).Return(tc.mockRet, tc.mockErr).Build()

			ok, err := dao.IsVideoExists(context.Background(), tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, ok)
			}
		})
	}
}
