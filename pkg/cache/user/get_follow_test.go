package usercache

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestUserCache_GetFollowing(t *testing.T) {
	type testCase struct {
		userID  int64
		stored  []string
		mockErr error
		want    []int64
		wantErr bool
	}

	testCases := map[string]testCase{
		"get following success": {
			userID: 1,
			stored: []string{
				"2",
				"3",
			},
			want: []int64{2, 3},
		},
		"empty set returns empty slice": {
			userID: 1,
			stored: []string{},
			want:   []int64{},
		},
		"redis error returns error": {
			userID:  1,
			mockErr: assert.AnError,
			wantErr: true,
		},
		"parse error returns error": {
			userID:  1,
			stored:  []string{"bad"},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewUserCache(db)

			key := getFollowingKey(tc.userID)
			if tc.mockErr != nil {
				mock.ExpectSMembers(key).SetErr(tc.mockErr)
			} else {
				mock.ExpectSMembers(key).SetVal(tc.stored)
			}

			ids, err := cache.GetFollowing(context.Background(), tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.ElementsMatch(t, tc.want, ids)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserCache_IsFollowed(t *testing.T) {
	type testCase struct {
		userID      int64
		followingID int64
		result      bool
		mockErr     error
		wantErr     bool
	}

	testCases := map[string]testCase{
		"user is followed": {
			userID:      1,
			followingID: 2,
			result:      true,
		},
		"user is not followed": {
			userID:      1,
			followingID: 3,
			result:      false,
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
			mock.ExpectExists(key).SetVal(1)
			if tc.mockErr != nil {
				mock.ExpectSIsMember(key, tc.followingID).SetErr(tc.mockErr)
			} else {
				mock.ExpectSIsMember(key, tc.followingID).SetVal(tc.result)
			}

			followed, err := cache.IsFollowed(context.Background(), tc.userID, tc.followingID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.result, followed)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
