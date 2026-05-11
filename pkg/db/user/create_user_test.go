package userdao

import (
	"context"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	type testCase struct {
		username string
		password string
		mockID   int64
		mockErr  error
		wantErr  bool
	}

	testCases := map[string]testCase{
		"create user success":    {username: "alice", password: "hash", mockID: 1},
		"db error returns error": {username: "alice", password: "hash", mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*UserDao).CreateUser).Return(tc.mockID, tc.mockErr).Build()

			id, err := dao.CreateUser(context.Background(), tc.username, tc.password)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockID, id)
			}
		})
	}
}
