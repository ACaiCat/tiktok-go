package commentdao

import (
	"context"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func TestGetCommentByID(t *testing.T) {
	type testCase struct {
		commentID int64
		mockRet   *model.Comment
		mockErr   error
		wantErr   bool
	}

	testCases := map[string]testCase{
		"get comment success":           {commentID: 1, mockRet: &model.Comment{ID: 1, Content: "nice"}},
		"comment not found returns nil": {commentID: 99, mockRet: nil},
		"db error returns error":        {commentID: 1, mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*CommentDao).GetCommentByID).Return(tc.mockRet, tc.mockErr).Build()

			c, err := dao.GetCommentByID(context.Background(), tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, c)
			}
		})
	}
}

func TestGetCommentsByVideoID(t *testing.T) {
	type testCase struct {
		videoID  int64
		pageSize int
		pageNum  int
		mockRet  []*model.Comment
		mockErr  error
		wantErr  bool
	}

	comments := []*model.Comment{{ID: 1, Content: "nice"}, {ID: 2, Content: "great"}}

	testCases := map[string]testCase{
		"get comments success": {videoID: 10, pageSize: 10, pageNum: 0, mockRet: comments},
		"empty result":         {videoID: 10, pageSize: 10, pageNum: 0, mockRet: []*model.Comment{}},
		"db error":             {videoID: 10, pageSize: 10, pageNum: 0, mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*CommentDao).GetCommentsByVideoID).Return(tc.mockRet, tc.mockErr).Build()

			cs, err := dao.GetCommentsByVideoID(context.Background(), tc.videoID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, cs)
			}
		})
	}
}

func TestGetCommentsByCommentID(t *testing.T) {
	type testCase struct {
		commentID int64
		pageSize  int
		pageNum   int
		mockRet   []*model.Comment
		mockErr   error
		wantErr   bool
	}

	replies := []*model.Comment{{ID: 3, Content: "reply"}}

	testCases := map[string]testCase{
		"get replies success": {commentID: 1, pageSize: 10, pageNum: 0, mockRet: replies},
		"empty replies":       {commentID: 1, pageSize: 10, pageNum: 0, mockRet: []*model.Comment{}},
		"db error":            {commentID: 1, pageSize: 10, pageNum: 0, mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*CommentDao).GetCommentsByCommentID).Return(tc.mockRet, tc.mockErr).Build()

			cs, err := dao.GetCommentsByCommentID(context.Background(), tc.commentID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, cs)
			}
		})
	}
}
