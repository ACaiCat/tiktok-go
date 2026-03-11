package service

import (
	"log"
	"mime/multipart"
	"os"

	"github.com/ACaiCat/tiktok-go/pkg/bucket"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/ACaiCat/tiktok-go/pkg/ffmpeg"
	"github.com/ACaiCat/tiktok-go/pkg/utils"
	"github.com/google/uuid"
)

func (s *VideoService) PublishVideo(userID int64, title string, description string, fileHeader *multipart.FileHeader) error {
	data, err := utils.FileHeaderToBytes(fileHeader)
	if err != nil {
		log.Printf("failed to read file header: %v\n", err)
		return errno.ServiceErr
	}

	tempFile, err := os.CreateTemp("", uuid.New().String()+fileHeader.Filename)
	if err != nil {
		log.Printf("failed to create temp file: %v\n", err)
		return errno.ServiceErr
	}

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(data); err != nil {
		log.Printf("failed to write to temp file: %v\n", err)
		return errno.ServiceErr
	}

	transcodeData, err := ffmpeg.TranscodeVideo(tempFile.Name())
	if err != nil {
		log.Printf("failed to transcode video: %v\n", err)
		return errno.ServiceErr
	}

	coverData, err := ffmpeg.GetVideoCover(tempFile.Name())
	if err != nil {
		log.Printf("failed to get video cover: %v\n", err)
		return errno.ServiceErr
	}

	err = s.videoDao.PublishVideo(
		userID, title, description,
		func(videoID int64) error {
			if err := bucket.UploadVideo(videoID, transcodeData); err != nil {
				log.Printf("failed to upload video: %v\n", err)
				return err
			}
			if err := bucket.UploadCover(videoID, coverData); err != nil {
				log.Printf("failed to upload cover: %v\n", err)
				return err
			}
			return nil
		},
		bucket.GetVideoURL,
		bucket.GetCoverURL,
	)
	if err != nil {
		log.Printf("failed to publish video with url: %v\n", err)
		return errno.ServiceErr
	}
	return nil
}
