package chatcache

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func TestChatCache_SetChatHistory(t *testing.T) {
	type testCase struct {
		userID      int64
		otherUserID int64
		pageSize    int
		pageNum     int
		messages    []*model.ChatMessage
		mockErr     error
		wantErr     bool
	}

	testCases := map[string]testCase{
		"set history success": {
			userID:      1,
			otherUserID: 2,
			pageSize:    10,
			pageNum:     0,
			messages:    testMessages,
			wantErr:     false,
		},
		"redis error returns error": {
			userID:      1,
			otherUserID: 2,
			pageSize:    10,
			pageNum:     0,
			messages:    testMessages,
			mockErr:     assert.AnError,
			wantErr:     true,
		},
		"key is normalized by user order": {
			userID:      2,
			otherUserID: 1,
			pageSize:    10,
			pageNum:     0,
			messages:    testMessages,
			wantErr:     false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewChatCache(db)

			key := getHistoryKey(tc.userID, tc.otherUserID, tc.pageSize, tc.pageNum)
			data, _ := json.Marshal(tc.messages)
			if tc.mockErr != nil {
				mock.ExpectSet(key, data, constants.ChatHistoryCacheExpiration).SetErr(tc.mockErr)
			} else {
				mock.ExpectSet(key, data, constants.ChatHistoryCacheExpiration).SetVal("OK")
			}

			err := cache.SetChatHistory(context.Background(), tc.userID, tc.otherUserID, tc.pageSize, tc.pageNum, tc.messages)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestChatCache_ClearChatHistory(t *testing.T) {
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
