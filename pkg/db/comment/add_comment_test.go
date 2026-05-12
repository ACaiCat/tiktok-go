package commentdao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestCommentDao_AddVideoComment(t *testing.T) {
	type testCase struct {
		userID  int64
		videoID int64
		content string
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"add video comment success": {userID: 1, videoID: 10, content: "nice"},
		"db error returns error":    {userID: 1, videoID: 10, content: "x", mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*CommentDao).AddVideoComment).Return(tc.mockErr).Build()

			err := dao.AddVideoComment(context.Background(), tc.userID, tc.videoID, tc.content)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCommentDao_AddCommentReply(t *testing.T) {
	type testCase struct {
		userID    int64
		videoID   int64
		commentID int64
		content   string
		mockErr   error
		wantErr   bool
	}

	testCases := map[string]testCase{
		"add comment reply success": {userID: 1, videoID: 10, commentID: 5, content: "agreed"},
		"db error returns error":    {userID: 1, videoID: 10, commentID: 5, content: "x", mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*CommentDao).AddCommentReply).Return(tc.mockErr).Build()

			err := dao.AddCommentReply(context.Background(), tc.userID, tc.videoID, tc.commentID, tc.content)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
