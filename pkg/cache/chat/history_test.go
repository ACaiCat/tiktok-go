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

func TestChatCache_GetChatHistory(t *testing.T) {
	type testCase struct {
		userID      int64
		otherUserID int64
		pageSize    int
		pageNum     int
		stored      []*model.ChatMessage
		mockErr     error
		wantErr     bool
	}

	testCases := map[string]testCase{
		"get history success": {
			userID:      1,
			otherUserID: 2,
			pageSize:    10,
			pageNum:     0,
			stored:      testMessages,
			wantErr:     false,
		},
		"cache miss returns error": {
			userID:      1,
			otherUserID: 2,
			pageSize:    10,
			pageNum:     0,
			mockErr:     assert.AnError,
			wantErr:     true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewChatCache(db)

			key := getHistoryKey(tc.userID, tc.otherUserID, tc.pageSize, tc.pageNum)
			if tc.mockErr != nil {
				mock.ExpectGet(key).SetErr(tc.mockErr)
			} else {
				data, _ := json.Marshal(tc.stored)
				mock.ExpectGet(key).SetVal(string(data))
			}

			msgs, err := cache.GetChatHistory(context.Background(), tc.userID, tc.otherUserID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, msgs, len(tc.stored))
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNormalizeConversationUserIDs(t *testing.T) {
	type testCase struct {
		userID      int64
		otherUserID int64
		wantLeft    int64
		wantRight   int64
	}

	testCases := map[string]testCase{
		"smaller id is left": {
			userID:      5,
			otherUserID: 10,
			wantLeft:    5,
			wantRight:   10,
		},
		"larger first still normalizes": {
			userID:      10,
			otherUserID: 5,
			wantLeft:    5,
			wantRight:   10,
		},
		"equal ids": {
			userID:      7,
			otherUserID: 7,
			wantLeft:    7,
			wantRight:   7,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			left, right := normalizeConversationUserIDs(tc.userID, tc.otherUserID)
			assert.Equal(t, tc.wantLeft, left)
			assert.Equal(t, tc.wantRight, right)
		})
	}
}
