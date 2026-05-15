package usercache

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func TestUserCache_SetLikeVideos(t *testing.T) {
	type testCase struct {
		userID   int64
		videoIDs []int64
		mockErr  error
		wantErr  bool
	}

	testCases := map[string]testCase{
		"set liked videos success": {
			userID:   1,
			videoIDs: []int64{10, 20, 30},
			wantErr:  false,
		},
		"pipeline error returns error": {
			userID:   1,
			videoIDs: []int64{10},
			mockErr:  assert.AnError,
			wantErr:  true,
		},
		"empty video list success": {
			userID:   1,
			videoIDs: []int64{},
			wantErr:  false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewUserCache(db)

			key := getLikedVideosKey(tc.userID)
			args := make([]interface{}, len(tc.videoIDs))
			for i, v := range tc.videoIDs {
				args[i] = v
			}
			if tc.mockErr != nil {
				mock.ExpectSAdd(key, args...).SetErr(tc.mockErr)
			} else {
				mock.ExpectSAdd(key, args...).SetVal(int64(len(tc.videoIDs)))
				mock.ExpectExpire(key, constants.LikeCacheExpiration).SetVal(true)
			}

			err := cache.SetLikeVideos(context.Background(), tc.userID, tc.videoIDs)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserCache_SetLikeVideo(t *testing.T) {
	type testCase struct {
		userID  int64
		videoID int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"set single like success": {
			userID:  1,
			videoID: 10,
			wantErr: false,
		},
		"pipeline error returns error": {
			userID:  1,
			videoID: 10,
			mockErr: assert.AnError,
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewUserCache(db)

			key := getLikedVideosKey(tc.userID)
			if tc.mockErr != nil {
				mock.ExpectSAdd(key, tc.videoID).SetErr(tc.mockErr)
			} else {
				mock.ExpectSAdd(key, tc.videoID).SetVal(1)
				mock.ExpectExpire(key, constants.LikeCacheExpiration).SetVal(true)
			}

			err := cache.SetLikeVideo(context.Background(), tc.userID, tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserCache_SetUnlikeVideo(t *testing.T) {
	type testCase struct {
		userID  int64
		videoID int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"unlike video success": {
			userID:  1,
			videoID: 10,
			wantErr: false,
		},
		"redis error returns error": {
			userID:  1,
			videoID: 10,
			mockErr: assert.AnError,
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewUserCache(db)

			key := getLikedVideosKey(tc.userID)
			if tc.mockErr != nil {
				mock.ExpectSRem(key, tc.videoID).SetErr(tc.mockErr)
			} else {
				mock.ExpectSRem(key, tc.videoID).SetVal(1)
			}

			err := cache.SetUnlikeVideo(context.Background(), tc.userID, tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserCache_ClearLikedVideos(t *testing.T) {
	type testCase struct {
		userID  int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"clear liked videos success": {
			userID:  1,
			wantErr: false,
		},
		"redis error returns error": {
			userID:  1,
			mockErr: assert.AnError,
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewUserCache(db)

			key := getLikedVideosKey(tc.userID)
			if tc.mockErr != nil {
				mock.ExpectDel(key).SetErr(tc.mockErr)
			} else {
				mock.ExpectDel(key).SetVal(1)
			}

			err := cache.ClearLikedVideos(context.Background(), tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
