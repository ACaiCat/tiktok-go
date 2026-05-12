package service

import (
	"context"
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/biz/model/ws"
	chatCache "github.com/ACaiCat/tiktok-go/pkg/cache/chat"
	chatDao "github.com/ACaiCat/tiktok-go/pkg/db/chat"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func TestChatService_HandleHistoryMessage(t *testing.T) {
	type testCase struct {
		mockHistoryErr error
		mockSendErr    error
		expectError    string
	}

	testCases := map[string]testCase{
		"success":       {},
		"history error": {mockHistoryErr: assert.AnError, expectError: assert.AnError.Error()},
		"send error":    {mockSendErr: assert.AnError, expectError: assert.AnError.Error()},
	}

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*ChatService).getChatHistory).To(func(userID int64, otherUserID int64, pageSize int, pageNum int) ([]*modelDao.ChatMessage, error) {
				return []*modelDao.ChatMessage{}, tc.mockHistoryErr
			}).Build()
			mockey.Mock((*ChatService).sendMessage).To(func(userID int64, msgType int, body any) (bool, error) {
				return true, tc.mockSendErr
			}).Build()
			mockey.Mock(NewChatService).To(func(_ context.Context, manager *ws.OnlineUserManager) *ChatService { return &ChatService{} }).Build()
			err := NewChatService(context.Background(), ws.NewOnlineUserManager()).HandleHistoryMessage(1, &ws.HistoryRequest{Sender: 2})
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestChatService_getChatHistory(t *testing.T) {
	type testCase struct {
		mockCacheResult []*modelDao.ChatMessage
		mockCacheErr    error
		mockDAOResult   []*modelDao.ChatMessage
		mockDAOErr      error
		expectError     string
	}

	testCases := map[string]testCase{
		"success from cache": {mockCacheResult: []*modelDao.ChatMessage{{ID: 1}}},
		"success from dao":   {mockCacheErr: redis.Nil, mockDAOResult: []*modelDao.ChatMessage{{ID: 1}}},
		"dao error":          {mockCacheErr: redis.Nil, mockDAOErr: assert.AnError, expectError: assert.AnError.Error()},
	}

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*chatCache.ChatCache).GetChatHistory).To(func(ctx context.Context, userID int64,
				otherUserID int64, pageSize int, pageNum int) ([]*modelDao.ChatMessage, error) {
				return tc.mockCacheResult, tc.mockCacheErr
			}).Build()
			mockey.Mock((*chatDao.ChatDao).GetChatHistory).To(func(ctx context.Context, userID int64, otherUserID int64,
				pageSize int, pageNum int) ([]*modelDao.ChatMessage, error) {
				return tc.mockDAOResult, tc.mockDAOErr
			}).Build()
			mockey.Mock((*chatCache.ChatCache).SetChatHistory).Return(nil).Build()
			mockey.Mock(NewChatService).To(func(_ context.Context, manager *ws.OnlineUserManager) *ChatService {
				return &ChatService{cache: &chatCache.ChatCache{}, chatDao: &chatDao.ChatDao{}}
			}).Build()
			result, err := NewChatService(context.Background(), ws.NewOnlineUserManager()).
				getChatHistory(1, 2, 10, 0)
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}
			assert.NoError(t, err)
			if errors.Is(tc.mockCacheErr, redis.Nil) {
				assert.Equal(t, tc.mockDAOResult, result)
				return
			}
			assert.Equal(t, tc.mockCacheResult, result)
		})
	}
}
