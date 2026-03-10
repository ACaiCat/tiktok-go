package bucket

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/minio/minio-go/v7"
)

func avatarObject(userID int64) string {
	return fmt.Sprintf("avatar_%d", userID)
}

func UploadAvatar(userID int64, data []byte) error {
	ctx := context.Background()
	_, err := Bucket.PutObject(ctx, constants.AvatarBucketName, avatarObject(userID),
		bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: "image/jpeg"},
	)
	return err
}

func GetAvatarURL(userID int64) string {
	return objectURL(constants.AvatarBucketName, avatarObject(userID))
}
