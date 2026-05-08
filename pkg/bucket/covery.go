package bucket

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
)

func coverObject(videoID int64) string {
	return fmt.Sprintf("cover_%d", videoID)
}

func UploadCover(ctx context.Context, videoID int64, data []byte) error {
	_, err := Bucket.PutObject(ctx, constants.CoverBucketName, coverObject(videoID),
		bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: "image/jpeg"},
	)
	return errors.Wrapf(err, "UploadCover failed, bucket=%s, object=%d", constants.CoverBucketName, videoID)
}

func GetCoverURL(videoID int64) string {
	return objectURL(constants.CoverBucketName, coverObject(videoID))
}
