package likedao

import (
	"context"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestIsCommentLikeExists(t *testing.T) {
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

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*LikeDao).IsCommentLikeExists).Return(tc.mockRet, tc.mockErr).Build()

			ok, err := dao.IsCommentLikeExists(context.Background(), tc.userID, tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, ok)
			}
		})
	}
}
