package ffmpeg

import (
	"bytes"

	"github.com/pkg/errors"
	ffmpeggo "github.com/u2takey/ffmpeg-go"
)

func GetVideoCover(filePath string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	err := ffmpeggo.Input(filePath).
		Output("pipe:1", ffmpeggo.KwArgs{"vframes": "1", "format": "mjpeg"}).
		WithOutput(buf).
		Run()
	if err != nil {
		return nil, errors.Wrapf(err, "GetVideoCover failed, filePath: %s", filePath)
	}

	return buf.Bytes(), nil
}
