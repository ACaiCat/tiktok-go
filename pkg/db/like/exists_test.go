package likedao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func TestLikeDao_IsCommentLikeExists(t *testing.T) {
	type testCase struct {
		userID    int64
		commentID int64
		mockRet   bool
		mockErr   error
		wantErr   bool
	}

	testCases := map[string]testCase{
		"like exists":            {userID: 1, commentID: 5, mockRet: true},
		"like not exists":        {userID: 1, commentID: 5, mockRet: false},
		"db error returns error": {userID: 1, commentID: 5, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockLikeQueryChain()
			dao := newTestDao()
			count := int64(0)
			if tc.mockRet {
				count = 1
			}
			dbtestutil.MockCount(count, tc.mockErr)

			ok, err := dao.IsCommentLikeExists(context.Background(), tc.userID, tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "IsCommentLikeExists failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, ok)
			}
		})
	}
}
