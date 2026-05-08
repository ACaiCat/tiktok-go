package img

import (
	"bytes"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"

	"github.com/pkg/errors"
)

func GetImageFormat(data []byte) (string, error) {
	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", errors.Wrapf(err, "GetImageFormat failed")
	}
	return format, nil
}

func ConvertToJPEG(data []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrapf(err, "ConvertToJPEG failed")
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		return nil, errors.Wrap(err, "ConvertToJPEG failed")
	}

	return buf.Bytes(), nil
}
