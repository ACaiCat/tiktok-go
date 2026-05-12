package userdao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gen"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func TestUserDao_CreateUser(t *testing.T) {
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
		mockey.PatchConvey(name, t, func() {
			mockUserQueryChain()
			dao := newTestDao()

			mockey.Mock((*gen.DO).Create).To(func(_ *gen.DO, value interface{}) error {
				if users, ok := value.([]*model.User); ok && len(users) > 0 && !tc.wantErr {
					users[0].ID = tc.mockID
				}
				return tc.mockErr
			}).Build()

			id, err := dao.CreateUser(context.Background(), tc.username, tc.password)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "CreateUser failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockID, id)
			}
		})
	}
}
