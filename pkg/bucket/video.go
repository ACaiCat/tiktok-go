package bucket

import (
	"bytes"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func videoObject(videoID int64) string {
	return fmt.Sprintf("video_%d", videoID)
}

func UploadVideo(ctx context.Context, videoID int64, data []byte) error {
	_, err := Bucket.PutObject(ctx, constants.VideoBucketName, videoObject(videoID),
		bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: "video/mp4"},
	)
	return err
}

func GetVideoURL(videoID int64) string {
	return objectURL(constants.VideoBucketName, videoObject(videoID))
}
