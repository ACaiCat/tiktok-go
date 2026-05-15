package videocache

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func TestVideoCache_GetVideo(t *testing.T) {
	baseTime := time.UnixMilli(1710000000000)
	video := &model.Video{
		ID:           1,
		UserID:       10,
		VideoURL:     "https://example.com/video-1.mp4",
		CoverURL:     "https://example.com/video-1.jpg",
		Title:        "video-1",
		Description:  "desc-1",
		VisitCount:   100,
		LikeCount:    20,
		CommentCount: 3,
		CreatedAt:    baseTime,
	}

	testCases := map[string]struct {
		videoID       int64
		stored        map[string]string
		mockErr       error
		wantErr       bool
		wantRedisNil  bool
		wantCreatedAt time.Time
	}{
		"get video success": {
			videoID:       video.ID,
			stored:        hashForVideo(video),
			wantCreatedAt: video.CreatedAt,
		},
		"cache miss returns redis nil": {
			videoID:      video.ID,
			stored:       map[string]string{},
			wantErr:      true,
			wantRedisNil: true,
		},
		"redis error returns error": {
			videoID: video.ID,
			mockErr: assert.AnError,
			wantErr: true,
		},
		"invalid hash data returns error": {
			videoID: video.ID,
			stored: map[string]string{
				"id":            "1",
				"user_id":       "10",
				"video_url":     "https://example.com/video-1.mp4",
				"cover_url":     "https://example.com/video-1.jpg",
				"title":         "video-1",
				"description":   "desc-1",
				"visit_count":   "100",
				"like_count":    "20",
				"comment_count": "3",
				"created_at":    "invalid",
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewVideoCache(db)

			key := getVideoKey(tc.videoID)
			if tc.mockErr != nil {
				mock.ExpectHGetAll(key).SetErr(tc.mockErr)
			} else {
				mock.ExpectHGetAll(key).SetVal(tc.stored)
			}

			got, err := cache.GetVideo(context.Background(), tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
				if tc.wantRedisNil {
					assert.ErrorIs(t, err, redis.Nil)
				}
				assert.NoError(t, mock.ExpectationsWereMet())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, video.ID, got.ID)
			assert.Equal(t, video.UserID, got.UserID)
			assert.Equal(t, video.VideoURL, got.VideoURL)
			assert.Equal(t, video.CoverURL, got.CoverURL)
			assert.Equal(t, video.Title, got.Title)
			assert.Equal(t, video.Description, got.Description)
			assert.Equal(t, video.VisitCount, got.VisitCount)
			assert.Equal(t, video.LikeCount, got.LikeCount)
			assert.Equal(t, video.CommentCount, got.CommentCount)
			assert.True(t, tc.wantCreatedAt.Equal(got.CreatedAt))
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestVideoCache_GetVideos(t *testing.T) {
	baseTime := time.UnixMilli(1710000000000)
	videos := []*model.Video{
		{
			ID:           2,
			UserID:       20,
			VideoURL:     "https://example.com/video-2.mp4",
			CoverURL:     "https://example.com/video-2.jpg",
			Title:        "video-2",
			Description:  "desc-2",
			VisitCount:   200,
			LikeCount:    30,
			CommentCount: 4,
			CreatedAt:    baseTime,
		},
		{
			ID:           1,
			UserID:       10,
			VideoURL:     "https://example.com/video-1.mp4",
			CoverURL:     "https://example.com/video-1.jpg",
			Title:        "video-1",
			Description:  "desc-1",
			VisitCount:   100,
			LikeCount:    20,
			CommentCount: 3,
			CreatedAt:    baseTime.Add(time.Minute),
		},
	}

	t.Run("get videos success and preserve order", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		cache := NewVideoCache(db)

		for _, video := range videos {
			mock.ExpectHGetAll(getVideoKey(video.ID)).SetVal(hashForVideo(video))
		}

		got, err := cache.GetVideos(context.Background(), []int64{2, 1})
		require.NoError(t, err)
		require.Len(t, got, 2)
		require.NotNil(t, got[0])
		require.NotNil(t, got[1])
		assert.Equal(t, int64(2), got[0].ID)
		assert.Equal(t, int64(1), got[1].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("partial cache miss keeps index with nil", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		cache := NewVideoCache(db)

		mock.ExpectHGetAll(getVideoKey(2)).SetVal(hashForVideo(videos[0]))
		mock.ExpectHGetAll(getVideoKey(999)).SetVal(map[string]string{})
		mock.ExpectHGetAll(getVideoKey(1)).SetVal(hashForVideo(videos[1]))

		got, err := cache.GetVideos(context.Background(), []int64{2, 999, 1})
		require.NoError(t, err)
		require.Len(t, got, 3)
		assert.NotNil(t, got[0])
		assert.Nil(t, got[1])
		assert.NotNil(t, got[2])
		assert.Equal(t, int64(2), got[0].ID)
		assert.Equal(t, int64(1), got[2].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty ids returns empty slice", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		cache := NewVideoCache(db)

		got, err := cache.GetVideos(context.Background(), []int64{})
		require.NoError(t, err)
		assert.Empty(t, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("redis error returns error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		cache := NewVideoCache(db)

		mock.ExpectHGetAll(getVideoKey(2)).SetErr(assert.AnError)

		got, err := cache.GetVideos(context.Background(), []int64{2})
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("invalid hash data returns error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		cache := NewVideoCache(db)

		broken := hashForVideo(videos[0])
		broken["visit_count"] = "broken"
		mock.ExpectHGetAll(getVideoKey(2)).SetVal(broken)

		got, err := cache.GetVideos(context.Background(), []int64{2})
		assert.Error(t, err)
		assert.Nil(t, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func hashForVideo(video *model.Video) map[string]string {
	return map[string]string{
		"id":            strconv.FormatInt(video.ID, 10),
		"user_id":       strconv.FormatInt(video.UserID, 10),
		"video_url":     video.VideoURL,
		"cover_url":     video.CoverURL,
		"title":         video.Title,
		"description":   video.Description,
		"visit_count":   strconv.FormatInt(video.VisitCount, 10),
		"like_count":    strconv.FormatInt(video.LikeCount, 10),
		"comment_count": strconv.FormatInt(video.CommentCount, 10),
		"created_at":    strconv.FormatInt(video.CreatedAt.UnixMilli(), 10),
	}
}
