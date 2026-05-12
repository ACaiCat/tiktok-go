package userdao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gen"
)

func TestUserDao_IsUserExists(t *testing.T) {
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
		mockey.PatchConvey(name, t, func() {
			mockUserQueryChain()
			dao := newTestDao()

			mockey.Mock((*gen.DO).Count).To(func(_ *gen.DO) (int64, error) {
				if tc.mockRet {
					return 1, tc.mockErr
				}
				return 0, tc.mockErr
			}).Build()

			ok, err := dao.IsUserExists(context.Background(), tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "IsUserExists failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, ok)
			}
		})
	}
}
