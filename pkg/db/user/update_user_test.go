package userdao

import (
	"context"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserMFA(t *testing.T) {
	type testCase struct {
		userID  int64
		secret  string
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"update mfa success":     {userID: 1, secret: "TOTP_SECRET"},
		"db error returns error": {userID: 1, secret: "x", mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*UserDao).UpdateUserMFA).Return(tc.mockErr).Build()

			err := dao.UpdateUserMFA(context.Background(), tc.userID, tc.secret)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateUserAvatarURL(t *testing.T) {
	type testCase struct {
		userID    int64
		avatarURL string
		mockErr   error
		wantErr   bool
	}

	testCases := map[string]testCase{
		"update avatar url success": {userID: 1, avatarURL: "http://example.com/avatar.jpg"},
		"db error returns error":    {userID: 1, avatarURL: "x", mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*UserDao).UpdateUserAvatarURL).Return(tc.mockErr).Build()

			err := dao.UpdateUserAvatarURL(context.Background(), tc.userID, tc.avatarURL)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateUserJwch(t *testing.T) {
	type testCase struct {
		userID       int64
		jwchID       string
		jwchPassword string
		mockErr      error
		wantErr      bool
	}

	testCases := map[string]testCase{
		"update jwch success":    {userID: 1, jwchID: "202100000001", jwchPassword: "pw"},
		"db error returns error": {userID: 1, jwchID: "x", jwchPassword: "y", mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*UserDao).UpdateUserJwch).Return(tc.mockErr).Build()

			err := dao.UpdateUserJwch(context.Background(), tc.userID, tc.jwchID, tc.jwchPassword)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
