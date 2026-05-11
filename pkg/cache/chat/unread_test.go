package chatcache

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

var testMessages = []*model.ChatMessage{
	{ID: 1, SenderID: 10, ReceiverID: 20, Content: "hello", CreatedAt: time.Now()},
	{ID: 2, SenderID: 10, ReceiverID: 20, Content: "world", CreatedAt: time.Now()},
}

func TestSetUnreadMessages(t *testing.T) {
	type testCase struct {
		userID   int64
		senderID int64
		messages []*model.ChatMessage
		mockErr  error
		wantErr  bool
	}

	testCases := map[string]testCase{
		"set unread messages success": {
			userID:   20,
			senderID: 10,
			messages: testMessages,
			mockErr:  nil,
			wantErr:  false,
		},
		"redis error returns error": {
			userID:   20,
			senderID: 10,
			messages: testMessages,
			mockErr:  assert.AnError,
			wantErr:  true,
		},
		"empty messages success": {
			userID:   20,
			senderID: 10,
			messages: []*model.ChatMessage{},
			mockErr:  nil,
			wantErr:  false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewChatCache(db)

			data, _ := json.Marshal(tc.messages)
			key := getUnreadKey(tc.userID, tc.senderID)
			if tc.mockErr != nil {
				mock.ExpectSet(key, data, constants.ChatUnreadCacheExpiration).SetErr(tc.mockErr)
			} else {
				mock.ExpectSet(key, data, constants.ChatUnreadCacheExpiration).SetVal("OK")
			}

			err := cache.SetUnreadMessages(context.Background(), tc.userID, tc.senderID, tc.messages)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUnreadMessages(t *testing.T) {
	type testCase struct {
		userID   int64
		senderID int64
		stored   []*model.ChatMessage
		mockErr  error
		wantErr  bool
	}

	testCases := map[string]testCase{
		"get unread messages success": {
			userID:   20,
			senderID: 10,
			stored:   testMessages,
			wantErr:  false,
		},
		"key not found returns error": {
			userID:   20,
			senderID: 10,
			mockErr:  assert.AnError,
			wantErr:  true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewChatCache(db)

			key := getUnreadKey(tc.userID, tc.senderID)
			if tc.mockErr != nil {
				mock.ExpectGet(key).SetErr(tc.mockErr)
			} else {
				data, _ := json.Marshal(tc.stored)
				mock.ExpectGet(key).SetVal(string(data))
			}

			msgs, err := cache.GetUnreadMessages(context.Background(), tc.userID, tc.senderID)
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

func TestClearUnreadMessages(t *testing.T) {
	type testCase struct {
		userID   int64
		senderID int64
		mockErr  error
		wantErr  bool
	}

	testCases := map[string]testCase{
		"clear unread messages success": {
			userID:   20,
			senderID: 10,
			wantErr:  false,
		},
		"redis error returns error": {
			userID:   20,
			senderID: 10,
			mockErr:  assert.AnError,
			wantErr:  true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewChatCache(db)

			key := getUnreadKey(tc.userID, tc.senderID)
			if tc.mockErr != nil {
				mock.ExpectDel(key).SetErr(tc.mockErr)
			} else {
				mock.ExpectDel(key).SetVal(1)
			}

			err := cache.ClearUnreadMessages(context.Background(), tc.userID, tc.senderID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
