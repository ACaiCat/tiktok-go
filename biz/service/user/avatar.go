package service

import (
	"mime/multipart"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/bucket"
	"github.com/ACaiCat/tiktok-go/pkg/img"
	"github.com/ACaiCat/tiktok-go/pkg/utils"
)

func (s *UserService) UploadAvatar(fileHeader *multipart.FileHeader, userID int64) error {
	var err error

	data, err := utils.FileHeaderToBytes(fileHeader)
	if err != nil {
		return errors.WithMessagef(err, "service.UploadAvatar: read file failed, userID=%d", userID)
	}

	format, err := img.CheckAvatar(data)
	if err != nil {
		return err
	}

	if format == "png" {
		data, err = img.ConvertToJPEG(data)
		if err != nil {
			return errors.WithMessage(err, "service.UploadAvatar failed")
		}
	}

	err = bucket.UploadAvatar(s.ctx, userID, data)
	if err != nil {
		return errors.WithMessagef(err, "service.UploadAvatar: bucket.UploadAvatar failed, userID=%d", userID)
	}

	avatarURL := bucket.GetAvatarURL(userID)

	err = s.dao.UpdateUserAvatarURL(s.ctx, userID, avatarURL)
	if err != nil {
		return errors.WithMessagef(err, "service.UploadAvatar: db.UpdateUserAvatarURL failed, userID=%d", userID)
	}

	return nil
}
