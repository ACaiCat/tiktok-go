package chatdao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func TestChatDao_GetUnreadMessages(t *testing.T) {
	type testCase struct {
		userID   int64
		senderID int64
		mockRet  []*model.ChatMessage
		mockErr  error
		wantErr  bool
	}

	msgs := []*model.ChatMessage{{ID: 1, Content: "test"}}

	testCases := map[string]testCase{
		"get unread messages success": {userID: 2, senderID: 1, mockRet: msgs},
		"no unread messages":          {userID: 2, senderID: 1, mockRet: []*model.ChatMessage{}},
		"db error returns error":      {userID: 2, senderID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*ChatDao).GetUnreadMessages).Return(tc.mockRet, tc.mockErr).Build()

			got, err := dao.GetUnreadMessages(context.Background(), tc.userID, tc.senderID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, got)
			}
		})
	}
}

func TestChatDao_MarkMessagesAsRead(t *testing.T) {
	type testCase struct {
		userID   int64
		senderID int64
		mockErr  error
		wantErr  bool
	}

	testCases := map[string]testCase{
		"mark as read success":   {userID: 2, senderID: 1},
		"db error returns error": {userID: 2, senderID: 1, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*ChatDao).MarkMessagesAsRead).Return(tc.mockErr).Build()

			err := dao.MarkMessagesAsRead(context.Background(), tc.userID, tc.senderID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestChatDao_GetChatHistory(t *testing.T) {
	type testCase struct {
		userID      int64
		otherUserID int64
		pageSize    int
		pageNum     int
		mockRet     []*model.ChatMessage
		mockErr     error
		wantErr     bool
	}

	msgs := []*model.ChatMessage{{ID: 1, Content: "hi"}, {ID: 2, Content: "hello"}}

	testCases := map[string]testCase{
		"get history success":    {userID: 1, otherUserID: 2, pageSize: 10, pageNum: 0, mockRet: msgs},
		"empty history":          {userID: 1, otherUserID: 2, pageSize: 10, pageNum: 0, mockRet: []*model.ChatMessage{}},
		"db error returns error": {userID: 1, otherUserID: 2, pageSize: 10, pageNum: 0, mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*ChatDao).GetChatHistory).Return(tc.mockRet, tc.mockErr).Build()

			got, err := dao.GetChatHistory(context.Background(), tc.userID, tc.otherUserID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, got)
			}
		})
	}
}
