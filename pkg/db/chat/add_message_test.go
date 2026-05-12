package chatdao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestChatDao_AddMessage(t *testing.T) {
	type testCase struct {
		senderID   int64
		receiverID int64
		content    string
		isRead     bool
		isAi       bool
		mockErr    error
		wantErr    bool
	}

	testCases := map[string]testCase{
		"add message success":     {senderID: 1, receiverID: 2, content: "hello"},
		"add message with isRead": {senderID: 1, receiverID: 2, content: "hi", isRead: true},
		"db error returns error":  {senderID: 1, receiverID: 2, content: "fail", mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			dao := newTestDao()
			mockey.Mock((*ChatDao).AddMessage).Return(tc.mockErr).Build()

			err := dao.AddMessage(context.Background(), tc.senderID, tc.receiverID, tc.content, tc.isRead, tc.isAi)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
