package userdao

import (
	"context"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestIsUserExists(t *testing.T) {
	type testCase struct {
		userID  int64
		mockRet bool
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"user exists":            {userID: 1, mockRet: true},
		"user not exists":        {userID: 99, mockRet: false},
		"db error returns error": {userID: 1, mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*UserDao).IsUserExists).Return(tc.mockRet, tc.mockErr).Build()

			ok, err := dao.IsUserExists(context.Background(), tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, ok)
			}
		})
	}
}
