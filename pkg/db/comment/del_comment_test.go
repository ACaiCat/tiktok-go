package commentdao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func TestCommentDao_DeleteComment(t *testing.T) {
	type testCase struct {
		commentID int64
		mockErr   error
		wantErr   bool
	}

	testCases := map[string]testCase{
		"delete comment success": {commentID: 1},
		"db error returns error": {commentID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockCommentQueryChain()
			dao := newTestDao()
			dbtestutil.MockDelete(tc.mockErr)

			err := dao.DeleteComment(context.Background(), tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "DeleteComment failed")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCommentDao_DeleteCommentReply(t *testing.T) {
	type testCase struct {
		commentID int64
		mockErr   error
		wantErr   bool
	}

	testCases := map[string]testCase{
		"delete comment reply success": {commentID: 5},
		"db error returns error":       {commentID: 5, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockCommentQueryChain()
			dao := newTestDao()
			dbtestutil.MockDelete(tc.mockErr)

			err := dao.DeleteCommentReply(context.Background(), tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "DeleteCommentReply failed")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
