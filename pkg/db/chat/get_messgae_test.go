package chatdao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
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
			mockChatQueryChain()
			dao := newTestDao()
			dbtestutil.MockFind(tc.mockRet, tc.mockErr)

			got, err := dao.GetUnreadMessages(context.Background(), tc.userID, tc.senderID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetUnreadMessages failed")
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
			mockChatQueryChain()
			dao := newTestDao()
			dbtestutil.MockUpdate(tc.mockErr)

			err := dao.MarkMessagesAsRead(context.Background(), tc.userID, tc.senderID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "MarkMessagesAsRead failed")
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
			mockChatQueryChain()
			dao := newTestDao()
			dbtestutil.MockFind(tc.mockRet, tc.mockErr)

			got, err := dao.GetChatHistory(context.Background(), tc.userID, tc.otherUserID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "GetChatHistory failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, got)
			}
		})
	}
}

func TestChatDao_GetChatHistory_ExcludesMissingRecord(t *testing.T) {
	mockey.PatchConvey("history handles empty result", t, func() {
		mockChatQueryChain()
		dao := newTestDao()
		dbtestutil.MockFind([]*model.ChatMessage{}, nil)

		got, err := dao.GetChatHistory(context.Background(), 1, 2, 10, 0)
		assert.NoError(t, err)
		assert.Empty(t, got)
	})
}
