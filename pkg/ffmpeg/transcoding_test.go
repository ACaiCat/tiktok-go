package ffmpeg

import (
	"bytes"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	ffmpeggo "github.com/u2takey/ffmpeg-go"
)

func TestTranscodeVideo(t *testing.T) {
	type testCase struct {
		filePath string
		mockErr  error
		wantErr  bool
	}

	testCases := map[string]testCase{
		"successfully transcode video": {
			filePath: "/tmp/test.avi",
			mockErr:  nil,
			wantErr:  false,
		},
		"ffmpeg returns error": {
			filePath: "/tmp/nonexistent.avi",
			mockErr:  assert.AnError,
			wantErr:  true,
		},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			Mock((*ffmpeggo.Stream).Run).Return(tc.mockErr).Build()
			Mock((*bytes.Buffer).Bytes).Return([]byte{}).Build()

			result, err := TranscodeVideo(tc.filePath)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}
