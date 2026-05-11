package commentdao

import (
	"context"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestIsCommentExists(t *testing.T) {
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

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*CommentDao).IsCommentExists).Return(tc.mockRet, tc.mockErr).Build()

			ok, err := dao.IsCommentExists(context.Background(), tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, ok)
			}
		})
	}
}
