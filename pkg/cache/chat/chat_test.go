package chatcache

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestClearChatHistory(t *testing.T) {
	type testCase struct {
		userID      int64
		otherUserID int64
		scanKeys    []string
		scanErr     error
		delErr      error
		wantErr     bool
	}

	testCases := map[string]testCase{
		"clear history success with keys": {
			userID:      1,
			otherUserID: 2,
			scanKeys:    []string{"chat:history:1:2:10:0", "chat:history:1:2:10:1"},
			wantErr:     false,
		},
		"clear history success with no keys": {
			userID:      1,
			otherUserID: 2,
			scanKeys:    []string{},
			wantErr:     false,
		},
		"scan error returns error": {
			userID:      1,
			otherUserID: 2,
			scanErr:     assert.AnError,
			wantErr:     true,
		},
		"del error returns error": {
			userID:      1,
			otherUserID: 2,
			scanKeys:    []string{"chat:history:1:2:10:0"},
			delErr:      assert.AnError,
			wantErr:     true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewChatCache(db)

			pattern := getHistoryKeyPattern(tc.userID, tc.otherUserID)
			if tc.scanErr != nil {
				mock.ExpectScan(0, pattern, int64(100)).SetErr(tc.scanErr)
			} else {
				mock.ExpectScan(0, pattern, int64(100)).SetVal(tc.scanKeys, 0)
				if len(tc.scanKeys) > 0 {
					if tc.delErr != nil {
						mock.ExpectDel(tc.scanKeys...).SetErr(tc.delErr)
					} else {
						mock.ExpectDel(tc.scanKeys...).SetVal(int64(len(tc.scanKeys)))
					}
				}
			}

			err := cache.ClearChatHistory(context.Background(), tc.userID, tc.otherUserID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
