package likedao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestLikeDao_AddVideoLike(t *testing.T) {
	type testCase struct {
		userID  int64
		videoID int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"add video like success": {userID: 1, videoID: 10},
		"db error returns error": {userID: 1, videoID: 10, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*LikeDao).AddVideoLike).Return(tc.mockErr).Build()

			err := dao.AddVideoLike(context.Background(), tc.userID, tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLikeDao_AddCommentLike(t *testing.T) {
	type testCase struct {
		userID    int64
		commentID int64
		mockErr   error
		wantErr   bool
	}

	testCases := map[string]testCase{
		"add comment like success": {userID: 1, commentID: 5},
		"db error returns error":   {userID: 1, commentID: 5, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*LikeDao).AddCommentLike).Return(tc.mockErr).Build()

			err := dao.AddCommentLike(context.Background(), tc.userID, tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
