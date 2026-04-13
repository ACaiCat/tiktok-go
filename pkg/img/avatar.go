package img

import (
	"fmt"
	"log"
	"slices"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

const BToMB = 1024 * 1024

func CheckAvatar(data []byte) (string, error) {
	if len(data) > constants.AvatarMaxSize {
		log.Printf("avatar size exceeds the maximum limit: %d bytes", len(data))
		return "", errno.AvatarTooLargeErr.WithMessage(fmt.Sprintf("图片超过最大限制，请上传小于%dMB的图片", constants.AvatarMaxSize/BToMB))
	}

	format, err := GetImageFormat(data)
	if err != nil {
		log.Println("failed to get image format:", err)
		return "", errno.AvatarFormatErr
	}

	if slices.Contains(constants.AvatarFormat, format) {
		return format, nil
	}

	return "", errno.AvatarFormatErr
}
