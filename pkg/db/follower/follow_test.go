package followerdao

import (
	"context"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestAddFollow(t *testing.T) {
	type testCase struct {
		userID     int64
		followerID int64
		mockErr    error
		wantErr    bool
	}

	testCases := map[string]testCase{
		"add follow success": {userID: 1, followerID: 2},
		"db error":           {userID: 1, followerID: 2, mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*FollowerDao).AddFollow).Return(tc.mockErr).Build()

			err := dao.AddFollow(context.Background(), tc.userID, tc.followerID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteFollow(t *testing.T) {
	type testCase struct {
		userID     int64
		followerID int64
		mockErr    error
		wantErr    bool
	}

	testCases := map[string]testCase{
		"delete follow success": {userID: 1, followerID: 2},
		"db error":              {userID: 1, followerID: 2, mockErr: assert.AnError, wantErr: true},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*FollowerDao).DeleteFollow).Return(tc.mockErr).Build()

			err := dao.DeleteFollow(context.Background(), tc.userID, tc.followerID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
