package commentdao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func TestCommentDao_IsCommentExists(t *testing.T) {
	type testCase struct {
		commentID int64
		mockRet   bool
		mockErr   error
		wantErr   bool
	}

	testCases := map[string]testCase{
		"comment exists":         {commentID: 1, mockRet: true},
		"comment not exists":     {commentID: 99, mockRet: false},
		"db error returns error": {commentID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockCommentQueryChain()
			dao := newTestDao()
			count := int64(0)
			if tc.mockRet {
				count = 1
			}
			dbtestutil.MockCount(count, tc.mockErr)

			ok, err := dao.IsCommentExists(context.Background(), tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "IsCommentExists failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, ok)
			}
		})
	}
}
