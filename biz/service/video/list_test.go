package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/biz/model/interaction"
	"github.com/ACaiCat/tiktok-go/biz/model/video"
	userCache "github.com/ACaiCat/tiktok-go/pkg/cache/user"
	videoCache "github.com/ACaiCat/tiktok-go/pkg/cache/video"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	userDao "github.com/ACaiCat/tiktok-go/pkg/db/user"
	videoDao "github.com/ACaiCat/tiktok-go/pkg/db/video"
)

func TestVideoService_GetVideoList(t *testing.T) {
	type testCase struct {
		req            *video.ListReq
		mockVideos     []*modelDao.Video
		mockTotal      int64
		mockErr        error
		userExists     bool
		userExistsErr  error
		cacheEnabled   bool
		cacheVideoIDs  []int64
		cacheTotal     int64
		cacheListErr   error
		cacheDetails   []*modelDao.Video
		expectError    string
		expectDBCalled bool
	}

	testCases := map[string]testCase{
		"success": {
			req:            &video.ListReq{UserID: "1", PageNum: 0, PageSize: 10},
			mockVideos:     []*modelDao.Video{{ID: 1, UserID: 100, Title: "test"}},
			mockTotal:      1,
			userExists:     true,
			expectDBCalled: true,
		},
		"success from cache": {
			req:           &video.ListReq{UserID: "1", PageNum: 0, PageSize: 10},
			cacheEnabled:  true,
			cacheVideoIDs: []int64{1},
			cacheTotal:    1,
			cacheDetails:  []*modelDao.Video{{ID: 1, UserID: 100, Title: "cached"}},
		},
		"invalid user id": {
			req:         &video.ListReq{UserID: "not-a-number", PageNum: 0, PageSize: 10},
			expectError: "参数错误",
		},
		"user not exists": {
			req:            &video.ListReq{UserID: "1", PageNum: 0, PageSize: 10},
			userExists:     false,
			expectError:    "用户不存在",
			expectDBCalled: true,
		},
		"user exists check error": {
			req:            &video.ListReq{UserID: "1", PageNum: 0, PageSize: 10},
			userExistsErr:  assert.AnError,
			expectError:    assert.AnError.Error(),
			expectDBCalled: true,
		},
		"db error": {
			req:            &video.ListReq{UserID: "1", PageNum: 0, PageSize: 10},
			userExists:     true,
			mockErr:        assert.AnError,
			expectError:    assert.AnError.Error(),
			expectDBCalled: true,
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			if !tc.cacheEnabled && tc.cacheListErr == nil {
				tc.cacheListErr = redis.Nil
			}

			dbCalled := false

			mockey.Mock((*gorm.DB).Transaction).To(
				func(_ *gorm.DB, fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
					dbCalled = true
					return fc(&gorm.DB{})
				}).Build()
			mockey.Mock((*userDao.UserDao).WithTx).To(
				func(_ *userDao.UserDao, tx *gorm.DB) *userDao.UserDao {
					return &userDao.UserDao{}
				}).Build()
			mockey.Mock((*userDao.UserDao).IsUserExists).To(
				func(_ *userDao.UserDao, ctx context.Context, userID int64) (bool, error) {
					return tc.userExists, tc.userExistsErr
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
			mockey.Mock((*videoCache.VideoCache).GetUserVideoList).To(
				func(_ *videoCache.VideoCache, ctx context.Context, userID int64, pageSize int, pageNum int) ([]int64, int64, error) {
					return tc.cacheVideoIDs, tc.cacheTotal, tc.cacheListErr
				}).Build()
			mockey.Mock((*videoCache.VideoCache).GetVideos).To(
				func(_ *videoCache.VideoCache, ctx context.Context, videoIDs []int64) ([]*modelDao.Video, error) {
					return tc.cacheDetails, nil
				}).Build()
			mockey.Mock((*videoCache.VideoCache).SetUserVideoList).Return(nil).Build()
			mockey.Mock((*videoCache.VideoCache).SetVideos).Return(nil).Build()

			svc := &VideoService{
				videoDao:  &videoDao.VideoDao{},
				userDao:   &userDao.UserDao{},
				userCache: &userCache.UserCache{},
				ctx:       context.Background(),
			}
			if tc.cacheEnabled {
				svc.videoCache = &videoCache.VideoCache{}
			}

			result, total, err := svc.GetVideoList(tc.req)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				assert.Equal(t, tc.expectDBCalled, dbCalled)
				return
			}

			assert.NoError(t, err)
			if tc.cacheEnabled {
				assert.Equal(t, tc.cacheTotal, total)
				assert.Equal(t, VideosDaoToDto(tc.cacheDetails), result)
			} else {
				assert.Equal(t, tc.mockTotal, total)
				assert.Equal(t, VideosDaoToDto(tc.mockVideos), result)
			}
			assert.Equal(t, tc.expectDBCalled, dbCalled)
		})
	}

	time.Sleep(1 * time.Second) // 等待所有mock调用完成，避免测试结束后调用未完成的mock导致panic
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
			req: &interaction.ListLikeReq{UserID: "1", PageNum: 0, PageSize: 10},
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

			mockey.Mock((*videoCache.VideoCache).SetVideo).Return(nil).Build()
			mockey.Mock((*userCache.UserCache).GetLikedVideos).Return([]int64{}, nil).Build()
			mockey.Mock((*videoCache.VideoCache).SetVideos).Return(nil).Build()

			mockey.Mock(NewVideoService).To(func(_ context.Context) *VideoService {
				return &VideoService{videoCache: &videoCache.VideoCache{}}
			}).Build()

			result, err := NewVideoService(context.Background()).GetLikedVideos(tc.req)

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
