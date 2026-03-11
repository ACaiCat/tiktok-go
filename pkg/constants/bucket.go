package constants

import "fmt"

const (
	AvatarBucketName = "avatar"
	VideoBucketName  = "video"
	CoverBucketName  = "cover"
)

var (
	AvatarBucketPolicy = fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Action": ["s3:GetObject"],
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Resource": ["arn:aws:s3:::%s/*"],
				"Sid": ""
			}
		]
	}`, AvatarBucketName)

	VideoBucketPolicy = fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Action": ["s3:GetObject"],
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Resource": ["arn:aws:s3:::%s/*"],
				"Sid": ""
			}
		]
	}`, VideoBucketName)

	CoverBucketPolicy = fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Action": ["s3:GetObject"],
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Resource": ["arn:aws:s3:::%s/*"],
				"Sid": ""
			}
		]
	}`, CoverBucketName)
)
