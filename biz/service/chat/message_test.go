package service

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/biz/model/ws"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func TestChatService_HandleChatMessage(t *testing.T) {
	type testCase struct {
		mockEnsureFriend bool
		mockSendOnline   bool
		mockSendErr      error
		mockSaveResult   bool
		expectError      string
	}

	testCases := map[string]testCase{
		"not friend":      {mockEnsureFriend: false},
		"send self error": {mockEnsureFriend: true, mockSendErr: assert.AnError, expectError: assert.AnError.Error()},
		"success":         {mockEnsureFriend: true, mockSendOnline: true, mockSaveResult: false},
	}

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			callCount := 0
			mockey.Mock((*ChatService).ensureFriend).To(func(userID int64, receiverID int64) bool {
				return tc.mockEnsureFriend
			}).Build()
			mockey.Mock((*ChatService).sendMessage).To(func(userID int64, msgType int, body any) (bool, error) {
				callCount++
				if callCount == 1 {
					return true, tc.mockSendErr
				}
				return tc.mockSendOnline, nil
			}).Build()
			mockey.Mock((*ChatService).saveChatMessage).To(func(senderID int64, receiverID int64, content string, isRead bool, isAI bool) bool {
				return tc.mockSaveResult
			}).Build()
			mockey.Mock((*ChatService).replyWithAI).To(func(userID int64, receiverID int64) {}).Build()

			mockey.Mock(NewChatService).To(func(_ context.Context, manager *ws.OnlineUserManager) *ChatService {
				return &ChatService{}
			}).Build()

			err := NewChatService(t.Context(), ws.NewOnlineUserManager()).HandleChatMessage(1, &ws.ChatMessage{ReceiverID: 2, Content: "hello"})
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestChatService_SendErr(t *testing.T) {
	type testCase struct {
		online        bool
		wantSendError bool
	}

	testCases := map[string]testCase{
		"online user":  {online: true, wantSendError: true},
		"offline user": {online: false, wantSendError: false},
	}

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			called := false
			manager := ws.NewOnlineUserManager()
			if tc.online {
				manager.AddOnlineUser(1, nil)
			}

			mockey.Mock((*ws.OnlineUser).SendError).To(func(code int, message string) {
				called = true
			}).Build()
			mockey.Mock(NewChatService).To(func(_ context.Context, _ *ws.OnlineUserManager) *ChatService {
				return &ChatService{manager: manager}
			}).Build()

			NewChatService(t.Context(), manager).SendErr(1, errno.ServiceErr)

			assert.Equal(t, tc.wantSendError, called)
		})
	}
}
