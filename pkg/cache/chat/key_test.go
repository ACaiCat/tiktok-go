package chatcache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
