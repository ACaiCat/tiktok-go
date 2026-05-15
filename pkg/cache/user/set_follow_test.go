package usercache

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func TestUserCache_SetFollowings(t *testing.T) {
	type testCase struct {
		userID       int64
		followingIDs []int64
		mockErr      error
		wantErr      bool
	}

	testCases := map[string]testCase{
		"set followings success": {
			userID:       1,
			followingIDs: []int64{2, 3, 4},
		},
		"empty followings success": {
			userID:       1,
			followingIDs: []int64{},
		},
		"pipeline error returns error": {
			userID:       1,
			followingIDs: []int64{2},
			mockErr:      assert.AnError,
			wantErr:      true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewUserCache(db)

			key := getFollowingKey(tc.userID)
			args := make([]interface{}, len(tc.followingIDs))
			for i, v := range tc.followingIDs {
				args[i] = v
			}
			if tc.mockErr != nil {
				mock.ExpectSAdd(key, args...).SetErr(tc.mockErr)
			} else {
				mock.ExpectSAdd(key, args...).SetVal(int64(len(tc.followingIDs)))
				mock.ExpectExpire(key, constants.LikeCacheExpiration).SetVal(true)
			}

			err := cache.SetFollowings(context.Background(), tc.userID, tc.followingIDs)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserCache_SetFollow(t *testing.T) {
	type testCase struct {
		userID      int64
		followingID int64
		mockErr     error
		wantErr     bool
	}

	testCases := map[string]testCase{
		"set follow success": {
			userID:      1,
			followingID: 2,
		},
		"pipeline error returns error": {
			userID:      1,
			followingID: 2,
			mockErr:     assert.AnError,
			wantErr:     true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewUserCache(db)

			key := getFollowingKey(tc.userID)
			if tc.mockErr != nil {
				mock.ExpectSAdd(key, tc.followingID).SetErr(tc.mockErr)
			} else {
				mock.ExpectSAdd(key, tc.followingID).SetVal(1)
				mock.ExpectExpire(key, constants.LikeCacheExpiration).SetVal(true)
			}

			err := cache.SetFollow(context.Background(), tc.userID, tc.followingID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserCache_SetUnfollow(t *testing.T) {
	type testCase struct {
		userID      int64
		followingID int64
		mockErr     error
		wantErr     bool
	}

	testCases := map[string]testCase{
		"set unfollow success": {
			userID:      1,
			followingID: 2,
		},
		"redis error returns error": {
			userID:      1,
			followingID: 2,
			mockErr:     assert.AnError,
			wantErr:     true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewUserCache(db)

			key := getFollowingKey(tc.userID)
			if tc.mockErr != nil {
				mock.ExpectSRem(key, tc.followingID).SetErr(tc.mockErr)
			} else {
				mock.ExpectSRem(key, tc.followingID).SetVal(1)
			}

			err := cache.SetUnfollow(context.Background(), tc.userID, tc.followingID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserCache_ClearFollowing(t *testing.T) {
	type testCase struct {
		userID  int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"clear following success": {
			userID: 1,
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

			key := getFollowingKey(tc.userID)
			if tc.mockErr != nil {
				mock.ExpectDel(key).SetErr(tc.mockErr)
			} else {
				mock.ExpectDel(key).SetVal(1)
			}

			err := cache.ClearFollowing(context.Background(), tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
