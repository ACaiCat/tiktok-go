package usercache

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func TestSetLikeVideos(t *testing.T) {
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

func TestGetLikedVideos(t *testing.T) {
	type testCase struct {
		userID  int64
		stored  []string
		mockErr error
		wantLen int
		wantErr bool
	}

	testCases := map[string]testCase{
		"get liked videos success": {
			userID:  1,
			stored:  []string{"10", "20", "30"},
			wantLen: 3,
			wantErr: false,
		},
		"redis error returns error": {
			userID:  1,
			mockErr: assert.AnError,
			wantErr: true,
		},
		"empty set returns empty slice": {
			userID:  1,
			stored:  []string{},
			wantLen: 0,
			wantErr: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewUserCache(db)

			key := getLikedVideosKey(tc.userID)
			if tc.mockErr != nil {
				mock.ExpectSMembers(key).SetErr(tc.mockErr)
			} else {
				mock.ExpectSMembers(key).SetVal(tc.stored)
			}

			ids, err := cache.GetLikedVideos(context.Background(), tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, ids, tc.wantLen)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSetLikeVideo(t *testing.T) {
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

func TestSetUnlikeVideo(t *testing.T) {
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
				mock.ExpectDel(key).SetErr(tc.mockErr)
			} else {
				mock.ExpectDel(key).SetVal(1)
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

func TestIsVideoLiked(t *testing.T) {
	type testCase struct {
		userID  int64
		videoID int64
		result  bool
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"video is liked": {
			userID:  1,
			videoID: 10,
			result:  true,
			wantErr: false,
		},
		"video is not liked": {
			userID:  1,
			videoID: 99,
			result:  false,
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
				mock.ExpectSIsMember(key, tc.videoID).SetErr(tc.mockErr)
			} else {
				mock.ExpectSIsMember(key, tc.videoID).SetVal(tc.result)
			}

			liked, err := cache.IsVideoLiked(context.Background(), tc.userID, tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.result, liked)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestClearLikedVideos(t *testing.T) {
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
