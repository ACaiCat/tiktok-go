package img

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
)

func makeJPEGBytes() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, nil)
	return buf.Bytes()
}

func makePNGBytes() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	img.Set(0, 0, color.RGBA{0, 255, 0, 255})
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func TestGetImageFormat(t *testing.T) {
	type testCase struct {
		data       []byte
		wantFormat string
		wantErr    bool
	}

	testCases := map[string]testCase{
		"detect jpeg format": {
			data:       makeJPEGBytes(),
			wantFormat: "jpeg",
			wantErr:    false,
		},
		"detect png format": {
			data:       makePNGBytes(),
			wantFormat: "png",
			wantErr:    false,
		},
		"invalid image data": {
			data:    []byte("not an image"),
			wantErr: true,
		},
		"empty data": {
			data:    []byte{},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			format, err := GetImageFormat(tc.data)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantFormat, format)
		})
	}
}

func TestConvertToJPEG(t *testing.T) {
	type testCase struct {
		data    []byte
		wantErr bool
	}

	testCases := map[string]testCase{
		"convert png to jpeg": {
			data:    makePNGBytes(),
			wantErr: false,
		},
		"convert jpeg to jpeg": {
			data:    makeJPEGBytes(),
			wantErr: false,
		},
		"invalid data": {
			data:    []byte("not an image"),
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			result, err := ConvertToJPEG(tc.data)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, result)
			// Verify output is valid JPEG
			format, err := GetImageFormat(result)
			assert.NoError(t, err)
			assert.Equal(t, "jpeg", format)
		})
	}
}
