package chat

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/biz/model/ws"
	chatService "github.com/ACaiCat/tiktok-go/biz/service/chat"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func encodeMessage(msgType int, body interface{}) string {
	rawBody, _ := json.Marshal(body)
	msg := ws.Message{
		Type: msgType,
		Body: rawBody,
	}
	data, _ := json.Marshal(msg)
	return string(data)
}

func TestHandleMessage_Chat(t *testing.T) {
	type testCase struct {
		messageText string
		mockErr     error
		wantSendErr bool
	}

	chatMsg := ws.ChatMessage{ReceiverID: 2, Content: "hello"}
	testCases := map[string]testCase{
		"success": {
			messageText: encodeMessage(ws.MessageTypeChat, chatMsg),
		},
		"service error": {
			messageText: encodeMessage(ws.MessageTypeChat, chatMsg),
			mockErr:     assert.AnError,
			wantSendErr: true,
		},
		"invalid body": {
			messageText: `{"type":1,"body":"not-json-object"}`,
			wantSendErr: true,
		},
	}

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			manager := ws.NewOnlineUserManager()
			sendErrCalled := false

			mockey.Mock(chatService.NewChatService).To(func(_ context.Context, _ *ws.OnlineUserManager) *chatService.ChatService {
				return &chatService.ChatService{}
			}).Build()
			mockey.Mock((*chatService.ChatService).HandleChatMessage).To(func(_ int64, _ *ws.ChatMessage) error {
				return tc.mockErr
			}).Build()
			mockey.Mock((*chatService.ChatService).SendErr).To(func(_ int64, _ errno.ErrNo) {
				sendErrCalled = true
			}).Build()

			handleMessage(context.Background(), 1, tc.messageText)

			_ = manager
			if tc.wantSendErr {
				assert.True(t, sendErrCalled)
			} else {
				assert.False(t, sendErrCalled)
			}
		})
	}
}

func TestHandleMessage_Unread(t *testing.T) {
	type testCase struct {
		messageText string
		mockErr     error
		wantSendErr bool
	}

	unreadReq := ws.UnreadRequest{Sender: 2}
	testCases := map[string]testCase{
		"success": {
			messageText: encodeMessage(ws.MessageTypeUnread, unreadReq),
		},
		"service error": {
			messageText: encodeMessage(ws.MessageTypeUnread, unreadReq),
			mockErr:     assert.AnError,
			wantSendErr: true,
		},
		"invalid body": {
			messageText: `{"type":3,"body":"not-json-object"}`,
			wantSendErr: true,
		},
	}

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			sendErrCalled := false

			mockey.Mock(chatService.NewChatService).To(func(_ context.Context, _ *ws.OnlineUserManager) *chatService.ChatService {
				return &chatService.ChatService{}
			}).Build()
			mockey.Mock((*chatService.ChatService).HandleUnreadMessage).To(func(_ int64, _ *ws.UnreadRequest) error {
				return tc.mockErr
			}).Build()
			mockey.Mock((*chatService.ChatService).SendErr).To(func(_ int64, _ errno.ErrNo) {
				sendErrCalled = true
			}).Build()

			handleMessage(context.Background(), 1, tc.messageText)

			if tc.wantSendErr {
				assert.True(t, sendErrCalled)
			} else {
				assert.False(t, sendErrCalled)
			}
		})
	}
}

func TestHandleMessage_History(t *testing.T) {
	type testCase struct {
		messageText string
		mockErr     error
		wantSendErr bool
	}

	historyReq := ws.HistoryRequest{Sender: 2}
	testCases := map[string]testCase{
		"success": {
			messageText: encodeMessage(ws.MessageTypeHistory, historyReq),
		},
		"service error": {
			messageText: encodeMessage(ws.MessageTypeHistory, historyReq),
			mockErr:     assert.AnError,
			wantSendErr: true,
		},
		"invalid body": {
			messageText: `{"type":2,"body":"not-json-object"}`,
			wantSendErr: true,
		},
	}

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			sendErrCalled := false

			mockey.Mock(chatService.NewChatService).To(func(_ context.Context, _ *ws.OnlineUserManager) *chatService.ChatService {
				return &chatService.ChatService{}
			}).Build()
			mockey.Mock((*chatService.ChatService).HandleHistoryMessage).To(func(_ int64, _ *ws.HistoryRequest) error {
				return tc.mockErr
			}).Build()
			mockey.Mock((*chatService.ChatService).SendErr).To(func(_ int64, _ errno.ErrNo) {
				sendErrCalled = true
			}).Build()

			handleMessage(context.Background(), 1, tc.messageText)

			if tc.wantSendErr {
				assert.True(t, sendErrCalled)
			} else {
				assert.False(t, sendErrCalled)
			}
		})
	}
}

func TestHandleMessage_InvalidJSON(t *testing.T) {
	defer mockey.UnPatchAll()
	mockey.PatchConvey("invalid json", t, func() {
		sendErrCalled := false

		mockey.Mock(chatService.NewChatService).To(func(_ context.Context, _ *ws.OnlineUserManager) *chatService.ChatService {
			return &chatService.ChatService{}
		}).Build()
		mockey.Mock((*chatService.ChatService).SendErr).To(func(_ int64, _ errno.ErrNo) {
			sendErrCalled = true
		}).Build()

		handleMessage(context.Background(), 1, "not-valid-json")
		assert.True(t, sendErrCalled)
	})
}

func TestHandleMessage_UnknownType(t *testing.T) {
	defer mockey.UnPatchAll()
	mockey.PatchConvey("unknown message type", t, func() {
		sendErrCalled := false

		mockey.Mock(chatService.NewChatService).To(func(_ context.Context, _ *ws.OnlineUserManager) *chatService.ChatService {
			return &chatService.ChatService{}
		}).Build()
		mockey.Mock((*chatService.ChatService).SendErr).To(func(_ int64, _ errno.ErrNo) {
			sendErrCalled = true
		}).Build()

		handleMessage(context.Background(), 1, encodeMessage(999, nil))
		assert.True(t, sendErrCalled)
	})
}
