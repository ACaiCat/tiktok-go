package img

import (
	"testing"

	. "github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func TestCheckAvatar(t *testing.T) {
	type testCase struct {
		data    []byte
		wantErr error
		wantFmt string
	}

	bigData := make([]byte, 6*1024*1024) // > 5MB

	testCases := map[string]testCase{
		"valid jpeg avatar": {
			data:    makeJPEGBytes(),
			wantErr: nil,
			wantFmt: "jpeg",
		},
		"valid png avatar": {
			data:    makePNGBytes(),
			wantErr: nil,
			wantFmt: "png",
		},
		"file too large": {
			data:    bigData,
			wantErr: errno.AvatarTooLargeErr.WithMessage("图片超过最大限制，请上传小于5MB的图片"),
		},
		"invalid format": {
			data:    []byte("not an image"),
			wantErr: errno.AvatarFormatErr,
		},
	}

	for name, tc := range testCases {
		PatchConvey(name, t, func() {
			format, err := CheckAvatar(tc.data)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantFmt, format)
		})
	}
}
