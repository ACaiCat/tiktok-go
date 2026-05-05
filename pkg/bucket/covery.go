package bucket

import (
	"bytes"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func coverObject(videoID int64) string {
	return fmt.Sprintf("cover_%d", videoID)
}

func UploadCover(ctx context.Context, videoID int64, data []byte) error {
	_, err := Bucket.PutObject(ctx, constants.CoverBucketName, coverObject(videoID),
		bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: "image/jpeg"},
	)
	return err
}

func GetCoverURL(videoID int64) string {
	return objectURL(constants.CoverBucketName, coverObject(videoID))
}
