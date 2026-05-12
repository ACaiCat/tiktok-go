package followerdao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func TestFollowerDao_AddFollow(t *testing.T) {
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

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockFollowerQueryChain()
			dao := newTestDao()
			dbtestutil.MockCreate(tc.mockErr)

			err := dao.AddFollow(context.Background(), tc.userID, tc.followerID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "AddFollow failed")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFollowerDao_DeleteFollow(t *testing.T) {
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

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockFollowerQueryChain()
			dao := newTestDao()
			dbtestutil.MockDelete(tc.mockErr)

			err := dao.DeleteFollow(context.Background(), tc.userID, tc.followerID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "DeleteFollow failed")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
