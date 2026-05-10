package bucket

import (
	"context"
	"net/url"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/ACaiCat/tiktok-go/config"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

var Bucket *minio.Client

func InitMinIO(ctx context.Context) {
	bucketConfig := config.AppConfig.Minio
	var err error
	Bucket, err = minio.New(bucketConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(bucketConfig.AccessKey, bucketConfig.SecretKey, ""),
		Secure: bucketConfig.UseSSL,
	})
	if err != nil {
		hlog.Fatal(err)
	}

	cancel, err := Bucket.HealthCheck(constants.HealthCheckTime)
	if err != nil {
		hlog.Fatal(err)
	}
	defer cancel()

	initBucket(ctx, constants.AvatarBucketName, constants.AvatarBucketPolicy)
	initBucket(ctx, constants.VideoBucketName, constants.VideoBucketPolicy)
	initBucket(ctx, constants.CoverBucketName, constants.CoverBucketPolicy)
}

func initBucket(ctx context.Context, name, policy string) {
	exist, err := Bucket.BucketExists(ctx, name)
	if err != nil {
		hlog.Fatal(err)
	}
	if exist {
		return
	}
	if err = Bucket.MakeBucket(ctx, name, minio.MakeBucketOptions{}); err != nil {
		hlog.Fatal(err)
	}
	if err = Bucket.SetBucketPolicy(ctx, name, policy); err != nil {
		hlog.Fatal(err)
	}
}

func objectURL(bucketName, objectName string) string {
	scheme := "http"
	if config.AppConfig.Minio.ExternalUseSSL {
		scheme = "https"
	}
	u := url.URL{
		Scheme: scheme,
		Host:   config.AppConfig.Minio.ExternalEndpoint,
		Path:   bucketName + "/" + objectName,
	}
	return u.String()
}
