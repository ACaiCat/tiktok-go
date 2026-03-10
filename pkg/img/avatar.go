package img

import (
	"fmt"
	"log"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func CheckAvatar(data []byte) (string, error) {
	if len(data) > constants.AvatarMaxSize {
		log.Printf("avatar size exceeds the maximum limit: %d bytes", len(data))
		return "", errno.AvatarTooLargeErr.WithMessage(fmt.Sprintf("图片超过最大限制，请上传小于%dMB的图片", constants.AvatarMaxSize/(1024*1024)))
	}

	format, err := GetImageFormat(data)
	if err != nil {
		log.Println("failed to get image format:", err)
		return "", errno.AvatarFormatErr
	}

	for _, allowedFormat := range constants.AvatarFormat {
		if format == allowedFormat {
			return format, nil
		}
	}

	return "", errno.AvatarFormatErr
}
