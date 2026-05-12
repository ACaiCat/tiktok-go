package videodao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
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
			mockVideoQueryChain()
			dao := newTestDao()
			dbtestutil.MockUpdateColumn(tc.mockErr)

			err := dao.IncrVisitCount(context.Background(), tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "IncrVisitCount failed")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
