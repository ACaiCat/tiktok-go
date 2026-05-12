package videodao

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func TestVideoDao_PublishVideo(t *testing.T) {
	type testCase struct {
		userID      int64
		title       string
		description string
		mockErr     error
		wantErr     bool
	}

	testCases := map[string]testCase{
		"publish success": {userID: 1, title: "my video", description: "desc"},
		"db error":        {userID: 1, title: "my video", mockErr: assert.AnError, wantErr: true},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockVideoQueryChain()
			dao := newTestDao()
			dbtestutil.MockTransaction(nil)
			dbtestutil.MockCreateWithHook(func(value interface{}) {
				if tc.wantErr {
					return
				}
				video, ok := value.(*model.Video)
				if ok {
					video.ID = 99
				}
			}, tc.mockErr)
			dbtestutil.MockUpdates(tc.mockErr)

			err := dao.PublishVideo(
				context.Background(),
				tc.userID,
				tc.title,
				tc.description,
				func(videoID int64) error { return nil },
				func(videoID int64) string { return "url" },
				func(videoID int64) string { return "cover" },
			)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "PublishVideo failed")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
