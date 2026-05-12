package ai

import (
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestNewToolCallContext(t *testing.T) {
	type testCase struct {
		userIDs []int64
		wantLen int
	}

	testCases := map[string]testCase{
		"no users": {
			userIDs: []int64{},
			wantLen: 0,
		},
		"single user": {
			userIDs: []int64{1},
			wantLen: 1,
		},
		"multiple users": {
			userIDs: []int64{1, 2, 3},
			wantLen: 3,
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			ctx := NewToolCallContext(tc.userIDs...)
			assert.Equal(t, tc.wantLen, len(ctx.AllowedUserIDs))
		})
	}
}

func TestToolCallContext_CanAccessUser(t *testing.T) {
	type testCase struct {
		allowedUsers []int64
		checkUser    int64
		wantAccess   bool
	}

	testCases := map[string]testCase{
		"allowed user can access": {
			allowedUsers: []int64{1, 2, 3},
			checkUser:    2,
			wantAccess:   true,
		},
		"disallowed user cannot access": {
			allowedUsers: []int64{1, 2, 3},
			checkUser:    99,
			wantAccess:   false,
		},
		"empty allowlist denies all": {
			allowedUsers: []int64{},
			checkUser:    1,
			wantAccess:   false,
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			ctx := NewToolCallContext(tc.allowedUsers...)
			assert.Equal(t, tc.wantAccess, ctx.CanAccessUser(tc.checkUser))
		})
	}
}
