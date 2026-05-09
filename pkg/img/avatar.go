package img

import (
	"fmt"
	"slices"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

const BToMB = 1024 * 1024

func CheckAvatar(data []byte) (string, error) {
	if len(data) > constants.AvatarMaxSize {
		return "", errno.AvatarTooLargeErr.WithMessage(fmt.Sprintf("图片超过最大限制，请上传小于%dMB的图片", constants.AvatarMaxSize/BToMB))
	}

	format, err := GetImageFormat(data)
	if err != nil {
		return "", errno.AvatarFormatErr
	}

	if slices.Contains(constants.AvatarFormat, format) {
		return format, nil
	}

	return "", errno.AvatarFormatErr
}
