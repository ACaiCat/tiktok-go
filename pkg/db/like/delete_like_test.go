package likedao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func TestLikeDao_DeleteVideoLike(t *testing.T) {
	type testCase struct {
		userID  int64
		videoID int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"delete video like success": {userID: 1, videoID: 10},
		"db error returns error":    {userID: 1, videoID: 10, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockLikeQueryChain()
			dao := newTestDao()
			dbtestutil.MockDelete(tc.mockErr)

			err := dao.DeleteVideoLike(context.Background(), tc.userID, tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "DeleteVideoLike failed")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLikeDao_DeleteCommentLike(t *testing.T) {
	type testCase struct {
		userID    int64
		commentID int64
		mockErr   error
		wantErr   bool
	}

	testCases := map[string]testCase{
		"delete comment like success": {userID: 1, commentID: 5},
		"db error returns error":      {userID: 1, commentID: 5, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockLikeQueryChain()
			dao := newTestDao()
			dbtestutil.MockDelete(tc.mockErr)

			err := dao.DeleteCommentLike(context.Background(), tc.userID, tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "DeleteCommentLike failed")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
