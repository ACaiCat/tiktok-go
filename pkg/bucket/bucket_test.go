package bucket

import (
	"os"
	"testing"

	"github.com/minio/minio-go/v7"

	"github.com/ACaiCat/tiktok-go/config"
)

func TestMain(m *testing.M) {
	// Use a zero-value client; PutObject is always mocked in upload tests
	Bucket = &minio.Client{}
	config.AppConfig.Minio.ExternalEndpoint = "minio.example.com"
	config.AppConfig.Minio.ExternalUseSSL = false
	os.Exit(m.Run())
}

func setSSL(useSSL bool) {
	config.AppConfig.Minio.ExternalUseSSL = useSSL
}
