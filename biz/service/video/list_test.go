package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/biz/model/interaction"
	"github.com/ACaiCat/tiktok-go/biz/model/video"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	videoDao "github.com/ACaiCat/tiktok-go/pkg/db/video"
)

func TestVideoService_GetVideoList(t *testing.T) {
	type testCase struct {
		req         *video.ListReq
		mockVideos  []*modelDao.Video
		mockTotal   int64
		mockErr     error
		expectError string
	}

	testCases := map[string]testCase{
		"success": {
			req:        &video.ListReq{UserID: "1", PageNum: 0, PageSize: 10},
			mockVideos: []*modelDao.Video{{ID: 1, UserID: 100, Title: "test"}},
			mockTotal:  1,
		},
		"invalid user id": {
			req:         &video.ListReq{UserID: "not-a-number", PageNum: 0, PageSize: 10},
			expectError: "参数错误",
		},
		"db error": {
			req:         &video.ListReq{UserID: "1", PageNum: 0, PageSize: 10},
			mockErr:     assert.AnError,
			expectError: assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*gorm.DB).Transaction).To(
				func(_ *gorm.DB, fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
					return fc(&gorm.DB{})
				}).Build()
			mockey.Mock((*videoDao.VideoDao).GetVideosByUserID).To(
				func(_ *videoDao.VideoDao, ctx context.Context, userID int64, pageSize int, pageNum int) ([]*modelDao.Video, error) {
					return tc.mockVideos, tc.mockErr
				}).Build()
			mockey.Mock((*videoDao.VideoDao).GetVideoCountByUserID).To(
				func(_ *videoDao.VideoDao, ctx context.Context, userID int64) (int64, error) {
					return tc.mockTotal, tc.mockErr
				}).Build()
			mockey.Mock((*videoDao.VideoDao).WithTx).To(
				func(_ *videoDao.VideoDao, tx *gorm.DB) *videoDao.VideoDao {
					return &videoDao.VideoDao{}
				}).Build()

			mockey.Mock(NewVideoService).To(func(_ context.Context) *VideoService {
				return &VideoService{}
			}).Build()

			result, total, err := NewVideoService(t.Context()).GetVideoList(tc.req)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.mockTotal, total)
			assert.Equal(t, VideosDaoToDto(tc.mockVideos), result)
		})
	}
}

func TestVideoService_GetLikedVideos(t *testing.T) {
	type testCase struct {
		req         *interaction.ListLikeReq
		mockVideos  []*modelDao.Video
		mockErr     error
		expectError string
	}

	testCases := map[string]testCase{
		"success": {
			req:        &interaction.ListLikeReq{UserID: "1", PageNum: 0, PageSize: 10},
			mockVideos: []*modelDao.Video{{ID: 1, UserID: 100, Title: "liked"}},
		},
		"invalid user id": {
			req:         &interaction.ListLikeReq{UserID: "nan", PageNum: 0, PageSize: 10},
			expectError: "参数错误",
		},
		"db error": {
			req:         &interaction.ListLikeReq{UserID: "1", PageNum: 0, PageSize: 10},
			mockErr:     assert.AnError,
			expectError: assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*videoDao.VideoDao).GetUserLikeList).To(
				func(_ *videoDao.VideoDao, ctx context.Context, userID int64, pageSize int, pageNum int) ([]*modelDao.Video, error) {
					return tc.mockVideos, tc.mockErr
				}).Build()

			mockey.Mock(NewVideoService).To(func(_ context.Context) *VideoService {
				return &VideoService{}
			}).Build()

			result, err := NewVideoService(t.Context()).GetLikedVideos(tc.req)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, VideosDaoToDto(tc.mockVideos), result)
		})
	}
}
