package likedao

import (
	"context"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestDeleteVideoLike(t *testing.T) {
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

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*LikeDao).DeleteVideoLike).Return(tc.mockErr).Build()

			err := dao.DeleteVideoLike(context.Background(), tc.userID, tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteCommentLike(t *testing.T) {
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

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*LikeDao).DeleteCommentLike).Return(tc.mockErr).Build()

			err := dao.DeleteCommentLike(context.Background(), tc.userID, tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
