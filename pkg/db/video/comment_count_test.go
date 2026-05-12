package videodao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func TestVideoDao_IncrCommentCount(t *testing.T) {
	type testCase struct {
		videoID int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"incr comment count success": {videoID: 1},
		"db error":                   {videoID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockVideoQueryChain()
			dao := newTestDao()
			dbtestutil.MockUpdateColumn(tc.mockErr)

			err := dao.IncrCommentCount(context.Background(), tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "IncrCommentCount failed")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVideoDao_DecrCommentCount(t *testing.T) {
	type testCase struct {
		videoID int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"decr comment count success": {videoID: 1},
		"db error":                   {videoID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockVideoQueryChain()
			dao := newTestDao()
			dbtestutil.MockUpdateColumn(tc.mockErr)

			err := dao.DecrCommentCount(context.Background(), tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "DecrCommentCount failed")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
