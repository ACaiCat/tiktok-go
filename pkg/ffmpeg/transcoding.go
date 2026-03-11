package ffmpeg

import (
	"bytes"

	"github.com/u2takey/ffmpeg-go"
)

func TranscodeVideo(filePath string) ([]byte, error) {
	var err error

	buf := bytes.NewBuffer(nil)

	err = ffmpeg_go.Input(filePath).
		Output("pipe:1", ffmpeg_go.KwArgs{"format": "mp4", "vcodec": "libx264", "movflags": "frag_keyframe+empty_moov"}).
		WithOutput(buf).
		Run()

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
