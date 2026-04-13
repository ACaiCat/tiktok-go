package ffmpeg

import (
	"bytes"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func GetVideoCover(filePath string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	err := ffmpeg_go.Input(filePath).
		Output("pipe:1", ffmpeg_go.KwArgs{"vframes": "1", "format": "mjpeg"}).
		WithOutput(buf).
		Run()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
