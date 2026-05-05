package service

import (
	"log"
	"mime/multipart"

	"github.com/ACaiCat/tiktok-go/pkg/bucket"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/ACaiCat/tiktok-go/pkg/img"
	"github.com/ACaiCat/tiktok-go/pkg/utils"
)

func (s *UserService) UploadAvatar(fileHeader *multipart.FileHeader, userID int64) error {
	var err error

	data, err := utils.FileHeaderToBytes(fileHeader)
	if err != nil {
		log.Printf("failed to read file header: %v\n", err)
		return errno.ServiceErr
	}

	format, err := img.CheckAvatar(data)
	if err != nil {
		return err
	}

	if format == "png" {
		data, err = img.ConvertToJPEG(data)
		if err != nil {
			log.Printf("failed to convert PNG to JPEG: %v\n", err)
			return errno.ServiceErr
		}
	}

	err = bucket.UploadAvatar(s.ctx, userID, data)
	if err != nil {
		log.Printf("failed to upload avatar: %v\n", err)
		return errno.ServiceErr
	}

	avatarURL := bucket.GetAvatarURL(userID)

	err = s.dao.UpdateUserAvatarURL(s.ctx, userID, avatarURL)
	if err != nil {
		return errno.ServiceErr
	}

	return nil
}
