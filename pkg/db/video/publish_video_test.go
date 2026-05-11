package videodao

import (
	"context"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func TestPublishVideo(t *testing.T) {
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

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			dao := newTestDao()
			Mock((*VideoDao).PublishVideo).Return(tc.mockErr).Build()

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
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
