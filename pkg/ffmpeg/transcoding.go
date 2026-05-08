package ffmpeg

import (
	"bytes"

	"github.com/pkg/errors"
	ffmpeggo "github.com/u2takey/ffmpeg-go"
)

func TranscodeVideo(filePath string) ([]byte, error) {
	var err error

	buf := bytes.NewBuffer(nil)

	err = ffmpeggo.Input(filePath).
		Output("pipe:1", ffmpeggo.KwArgs{"format": "mp4", "vcodec": "libx264", "movflags": "frag_keyframe+empty_moov"}).
		WithOutput(buf).
		Run()

	if err != nil {
		return nil, errors.Wrapf(err, "TranscodeVideo failed, filePath: %s", filePath)
	}

	return buf.Bytes(), nil
}
