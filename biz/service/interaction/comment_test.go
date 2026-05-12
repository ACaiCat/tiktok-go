package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/biz/model/interaction"
	"github.com/ACaiCat/tiktok-go/biz/model/model"
	commentDao "github.com/ACaiCat/tiktok-go/pkg/db/comment"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	videoDao "github.com/ACaiCat/tiktok-go/pkg/db/video"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func TestInteractionService_CommentAction(t *testing.T) {
	type testCase struct {
		req             *interaction.CommentReq
		mockVideoExists bool
		mockVideoErr    error
		mockComment     *modelDao.Comment
		mockCommentErr  error
		expectError     string
	}

	testCases := map[string]testCase{
		"missing target": {
			req:         &interaction.CommentReq{Content: "hi"},
			expectError: "视频ID或评论ID不能为空",
		},
		"both targets": {
			req:         &interaction.CommentReq{VideoID: new("1"), CommentID: new("2"), Content: "hi"},
			expectError: "视频ID和评论ID不能同时存在",
		},
		"video not exist": {
			req:             &interaction.CommentReq{VideoID: new("1"), Content: "hi"},
			mockVideoExists: false,
			expectError:     errno.VideoNotExistErr.ErrMsg,
		},
		"video success": {
			req:             &interaction.CommentReq{VideoID: new("1"), Content: "hi"},
			mockVideoExists: true,
		},
		"reply success": {
			req:         &interaction.CommentReq{CommentID: new("2"), Content: "hi"},
			mockComment: &modelDao.Comment{ID: 2, VideoID: 1},
		},
		"reply comment not exist": {
			req:         &interaction.CommentReq{CommentID: new("2"), Content: "hi"},
			expectError: errno.CommentNotExistErr.ErrMsg,
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*gorm.DB).Transaction).To(func(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
				return fc(&gorm.DB{})
			}).Build()
			mockey.Mock((*commentDao.CommentDao).WithTx).To(func(tx *gorm.DB) *commentDao.CommentDao { return &commentDao.CommentDao{} }).Build()
			mockey.Mock((*videoDao.VideoDao).WithTx).To(func(tx *gorm.DB) *videoDao.VideoDao { return &videoDao.VideoDao{} }).Build()
			mockey.Mock((*videoDao.VideoDao).IsVideoExists).To(func(ctx context.Context, videoID int64) (bool, error) {
				return tc.mockVideoExists, tc.mockVideoErr
			}).Build()
			mockey.Mock((*commentDao.CommentDao).AddVideoComment).Return(nil).Build()
			mockey.Mock((*videoDao.VideoDao).IncrCommentCount).Return(nil).Build()
			mockey.Mock((*commentDao.CommentDao).GetCommentByID).To(func(ctx context.Context, commentID int64) (*modelDao.Comment, error) {
				return tc.mockComment, tc.mockCommentErr
			}).Build()
			mockey.Mock((*commentDao.CommentDao).AddCommentReply).Return(nil).Build()
			mockey.Mock((*commentDao.CommentDao).IncrCommentCount).Return(nil).Build()

			mockey.Mock(NewInteractionService).To(func(_ context.Context) *InteractionService {
				return &InteractionService{}
			}).Build()

			err := NewInteractionService(t.Context()).CommentAction(tc.req, 1)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestInteractionService_DeleteComment(t *testing.T) {
	type testCase struct {
		req         *interaction.DeleteCommentReq
		mockComment *modelDao.Comment
		mockErr     error
		expectError string
	}

	testCases := map[string]testCase{
		"invalid comment id": {
			req:         &interaction.DeleteCommentReq{CommentID: "bad"},
			expectError: "invalid syntax",
		},
		"comment not exist": {
			req:         &interaction.DeleteCommentReq{CommentID: "1"},
			expectError: errno.CommentNotExistErr.ErrMsg,
		},
		"not owner": {
			req:         &interaction.DeleteCommentReq{CommentID: "1"},
			mockComment: &modelDao.Comment{ID: 1, UserID: 2, VideoID: 1},
			expectError: errno.CommentNotBelongToUserErr.ErrMsg,
		},
		"success": {
			req:         &interaction.DeleteCommentReq{CommentID: "1"},
			mockComment: &modelDao.Comment{ID: 1, UserID: 1, VideoID: 1},
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*gorm.DB).Transaction).To(func(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
				return fc(&gorm.DB{})
			}).Build()
			mockey.Mock((*commentDao.CommentDao).WithTx).To(func(tx *gorm.DB) *commentDao.CommentDao { return &commentDao.CommentDao{} }).Build()
			mockey.Mock((*videoDao.VideoDao).WithTx).To(func(tx *gorm.DB) *videoDao.VideoDao { return &videoDao.VideoDao{} }).Build()
			mockey.Mock((*commentDao.CommentDao).GetCommentByID).To(func(ctx context.Context, commentID int64) (*modelDao.Comment, error) {
				return tc.mockComment, tc.mockErr
			}).Build()
			mockey.Mock((*videoDao.VideoDao).DecrCommentCount).Return(nil).Build()
			mockey.Mock((*commentDao.CommentDao).DecrCommentCount).Return(nil).Build()
			mockey.Mock((*commentDao.CommentDao).DeleteComment).Return(nil).Build()

			mockey.Mock(NewInteractionService).To(func(_ context.Context) *InteractionService {
				return &InteractionService{}
			}).Build()

			err := NewInteractionService(t.Context()).DeleteComment(tc.req, 1)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestInteractionService_ListComment(t *testing.T) {
	type testCase struct {
		req              *interaction.ListCommentReq
		mockVideoExists  bool
		mockCommentExist bool
		mockResult       []*modelDao.Comment
		mockErr          error
		expectError      string
	}

	testCases := map[string]testCase{
		"missing target": {
			req:         &interaction.ListCommentReq{},
			expectError: "视频ID或评论ID不能为空",
		},
		"video success": {
			req:             &interaction.ListCommentReq{VideoID: new("1")},
			mockVideoExists: true,
			mockResult:      []*modelDao.Comment{},
		},
		"video not exist": {
			req:             &interaction.ListCommentReq{VideoID: new("1")},
			mockVideoExists: false,
			expectError:     errno.VideoNotExistErr.ErrMsg,
		},
		"comment success": {
			req:              &interaction.ListCommentReq{CommentID: new("2")},
			mockCommentExist: true,
			mockResult:       []*modelDao.Comment{},
		},
		"comment not exist": {
			req:              &interaction.ListCommentReq{CommentID: new("2")},
			mockCommentExist: false,
			expectError:      errno.CommentNotExistErr.ErrMsg,
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*videoDao.VideoDao).IsVideoExists).To(func(ctx context.Context, videoID int64) (bool, error) {
				return tc.mockVideoExists, tc.mockErr
			}).Build()
			mockey.Mock((*commentDao.CommentDao).GetCommentsByVideoID).To(func(ctx context.Context, videoID int64,
				pageSize int, pageNum int) ([]*modelDao.Comment, error) {
				return tc.mockResult, tc.mockErr
			}).Build()
			mockey.Mock((*commentDao.CommentDao).IsCommentExists).To(func(ctx context.Context,
				commentID int64) (bool, error) {
				return tc.mockCommentExist, tc.mockErr
			}).Build()
			mockey.Mock((*commentDao.CommentDao).GetCommentsByCommentID).To(func(ctx context.Context, commentID int64,
				pageSize int, pageNum int) ([]*modelDao.Comment, error) {
				return tc.mockResult, tc.mockErr
			}).Build()

			mockey.Mock(NewInteractionService).To(func(_ context.Context) *InteractionService {
				return &InteractionService{}
			}).Build()

			result, err := NewInteractionService(t.Context()).ListComment(tc.req)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, CommentsDaoToDto(tc.mockResult), result)
		})
	}

	_ = model.Comment{}
}
