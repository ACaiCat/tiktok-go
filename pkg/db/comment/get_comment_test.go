package commentdao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func TestCommentDao_GetCommentByID(t *testing.T) {
	type testCase struct {
		commentID int64
		mockRet   *model.Comment
		mockErr   error
		wantErr   bool
	}

	testCases := map[string]testCase{
		"get comment success":           {commentID: 1, mockRet: &model.Comment{ID: 1, Content: "nice"}},
		"comment not found returns nil": {commentID: 99, mockRet: nil, mockErr: gorm.ErrRecordNotFound},
		"db error returns error":        {commentID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockCommentQueryChain()
			dao := newTestDao()
			dbtestutil.MockFirst(tc.mockRet, tc.mockErr)

			c, err := dao.GetCommentByID(context.Background(), tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetCommentByID failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, c)
			}
		})
	}
}

func TestCommentDao_GetCommentsByVideoID(t *testing.T) {
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

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockCommentQueryChain()
			dao := newTestDao()
			dbtestutil.MockFind(tc.mockRet, tc.mockErr)

			cs, err := dao.GetCommentsByVideoID(context.Background(), tc.videoID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetCommentsByVideoID failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, cs)
			}
		})
	}
}

func TestCommentDao_GetCommentsByCommentID(t *testing.T) {
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

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockCommentQueryChain()
			dao := newTestDao()
			dbtestutil.MockFind(tc.mockRet, tc.mockErr)

			cs, err := dao.GetCommentsByCommentID(context.Background(), tc.commentID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetCommentsByCommentID failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, cs)
			}
		})
	}
}

func TestCommentDao_IsCommentExists(t *testing.T) {
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

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockCommentQueryChain()
			dao := newTestDao()
			count := int64(0)
			if tc.mockRet {
				count = 1
			}
			dbtestutil.MockCount(count, tc.mockErr)

			ok, err := dao.IsCommentExists(context.Background(), tc.commentID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "IsCommentExists failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, ok)
			}
		})
	}
}
