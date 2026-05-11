package ffmpeg

import (
	"bytes"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	ffmpeggo "github.com/u2takey/ffmpeg-go"
)

func TestGetVideoCover(t *testing.T) {
	type testCase struct {
		filePath string
		mockErr  error
		mockData []byte
		wantErr  bool
	}

	testCases := map[string]testCase{
		"successfully extract cover": {
			filePath: "/tmp/test.mp4",
			mockErr:  nil,
			mockData: []byte{0xFF, 0xD8, 0xFF}, // JPEG magic bytes
			wantErr:  false,
		},
		"ffmpeg returns error": {
			filePath: "/tmp/nonexistent.mp4",
			mockErr:  assert.AnError,
			wantErr:  true,
		},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			Mock((*ffmpeggo.Stream).Run).Return(tc.mockErr).Build()
			Mock((*bytes.Buffer).Bytes).Return([]byte{}).Build()

			result, err := GetVideoCover(tc.filePath)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}
