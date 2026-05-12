package commentdao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func TestCommentDao_IncrCommentCount(t *testing.T) {
	type testCase struct {
		commentID int64
		mockErr   error
		wantErr   bool
	}

	testCases := map[string]testCase{
		"incr comment count success": {commentID: 1},
		"db error returns error":     {commentID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockCommentQueryChain()
			dao := newTestDao()
			dbtestutil.MockUpdateColumn(tc.mockErr)

			err := dao.IncrCommentCount(context.Background(), tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "IncrCommentCount failed")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCommentDao_DecrCommentCount(t *testing.T) {
	type testCase struct {
		commentID int64
		mockErr   error
		wantErr   bool
	}

	testCases := map[string]testCase{
		"decr comment count success": {commentID: 1},
		"db error returns error":     {commentID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockCommentQueryChain()
			dao := newTestDao()
			dbtestutil.MockUpdateColumn(tc.mockErr)

			err := dao.DecrCommentCount(context.Background(), tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "DecrCommentCount failed")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
