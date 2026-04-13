package bucket

import (
	"context"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/ACaiCat/tiktok-go/config"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

var Bucket *minio.Client

func InitMinIO() {
	bucketConfig := config.AppConfig.Minio
	var err error
	Bucket, err = minio.New(bucketConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(bucketConfig.AccessKey, bucketConfig.SecretKey, ""),
		Secure: bucketConfig.UseSSL,
	})
	if err != nil {
		panic(err)
	}

	cancel, err := Bucket.HealthCheck(constants.HealthCheckTime)
	if err != nil {
		panic(err)
	}
	defer cancel()

	initBucket(constants.AvatarBucketName, constants.AvatarBucketPolicy)
	initBucket(constants.VideoBucketName, constants.VideoBucketPolicy)
	initBucket(constants.CoverBucketName, constants.CoverBucketPolicy)
}

func initBucket(name, policy string) {
	ctx := context.Background()
	exist, err := Bucket.BucketExists(ctx, name)
	if err != nil {
		panic(err)
	}
	if exist {
		return
	}
	if err = Bucket.MakeBucket(ctx, name, minio.MakeBucketOptions{}); err != nil {
		panic(err)
	}
	if err = Bucket.SetBucketPolicy(ctx, name, policy); err != nil {
		panic(err)
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
