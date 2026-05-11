package commentdao

import (
	"context"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestIncrCommentCount(t *testing.T) {
	type testCase struct {
		commentID int64
		mockErr   error
		wantErr   bool
	}

	testCases := map[string]testCase{
		"incr comment count success": {commentID: 1},
		"db error returns error":     {commentID: 1, mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*CommentDao).IncrCommentCount).Return(tc.mockErr).Build()

			err := dao.IncrCommentCount(context.Background(), tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDecrCommentCount(t *testing.T) {
	type testCase struct {
		commentID int64
		mockErr   error
		wantErr   bool
	}

	testCases := map[string]testCase{
		"decr comment count success": {commentID: 1},
		"db error returns error":     {commentID: 1, mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*CommentDao).DecrCommentCount).Return(tc.mockErr).Build()

			err := dao.DecrCommentCount(context.Background(), tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
