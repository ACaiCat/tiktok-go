package commentdao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func TestCommentDao_IncrLikeCount(t *testing.T) {
	type testCase struct {
		commentID int64
		mockErr   error
		wantErr   bool
	}

	testCases := map[string]testCase{
		"incr like count success": {commentID: 1},
		"db error returns error":  {commentID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockCommentQueryChain()
			dao := newTestDao()
			dbtestutil.MockUpdateColumn(tc.mockErr)

			err := dao.IncrLikeCount(context.Background(), tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "IncrLikeCount failed")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCommentDao_DecrLikeCount(t *testing.T) {
	type testCase struct {
		commentID int64
		mockErr   error
		wantErr   bool
	}

	testCases := map[string]testCase{
		"decr like count success": {commentID: 1},
		"db error returns error":  {commentID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockCommentQueryChain()
			dao := newTestDao()
			dbtestutil.MockUpdateColumn(tc.mockErr)

			err := dao.DecrLikeCount(context.Background(), tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "DecrLikeCount failed")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
