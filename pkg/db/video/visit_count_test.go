package videodao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestVideoDao_IncrVisitCount(t *testing.T) {
	type testCase struct {
		videoID int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"incr visit count success": {videoID: 1},
		"db error":                 {videoID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*VideoDao).IncrVisitCount).Return(tc.mockErr).Build()

			err := dao.IncrVisitCount(context.Background(), tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
