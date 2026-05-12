package videodao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func TestVideoDao_IsVideoExists(t *testing.T) {
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

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockVideoQueryChain()
			dao := newTestDao()
			count := int64(0)
			if tc.mockRet {
				count = 1
			}
			dbtestutil.MockCount(count, tc.mockErr)

			ok, err := dao.IsVideoExists(context.Background(), tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "IsVideoExists failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, ok)
			}
		})
	}
}
