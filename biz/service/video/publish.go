package service

import (
	"mime/multipart"
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/bucket"
	"github.com/ACaiCat/tiktok-go/pkg/ffmpeg"
	"github.com/ACaiCat/tiktok-go/pkg/utils"
)

func (s *VideoService) PublishVideo(userID int64, title string, description string, fileHeader *multipart.FileHeader) error {
	data, err := utils.FileHeaderToBytes(fileHeader)
	if err != nil {
		return errors.WithMessagef(err, "service.PublishVideo: read file failed, userID=%d, filename=%q", userID, fileHeader.Filename)
	}

	tempFile, err := os.CreateTemp("", uuid.New().String()+fileHeader.Filename)
	if err != nil {
		hlog.CtxErrorf(s.ctx, "failed to create temp file: %v", err)
	}

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(data); err != nil {
		return errors.WithMessagef(err, "service.PublishVideo: write temp file failed, userID=%d", userID)
	}

	transcodeData, err := ffmpeg.TranscodeVideo(tempFile.Name())
	if err != nil {
		return errors.WithMessagef(err, "service.PublishVideo: ffmpeg.TranscodeVideo failed, userID=%d", userID)
	}

	coverData, err := ffmpeg.GetVideoCover(tempFile.Name())
	if err != nil {
		hlog.CtxErrorf(s.ctx, "failed to get video cover: %v", err)
	}

	err = s.videoDao.PublishVideo(
		s.ctx, userID, title, description,
		func(videoID int64) error {
			if err := bucket.UploadVideo(s.ctx, videoID, transcodeData); err != nil {
				return errors.WithMessagef(err, "service.PublishVideo: bucket.UploadVideo failed, videoID=%d, userID=%d", videoID, userID)
			}
			if err := bucket.UploadCover(s.ctx, videoID, coverData); err != nil {
				return errors.WithMessagef(err, "service.PublishVideo: bucket.UploadCover failed, videoID=%d, userID=%d", videoID, userID)
			}
			return nil
		},
		bucket.GetVideoURL,
		bucket.GetCoverURL,
	)
	if err != nil {
		return errors.WithMessagef(err, "service.PublishVideo: db.PublishVideo failed, userID=%d, title=%q", userID, title)
	}
	go func() {
		if err := s.videoCache.ClearUserVideoList(s.ctx, userID); err != nil {
			hlog.Errorf("service.PublishVideo: cache.ClearUserVideoList failed, userID=%d, err=%v", userID, err)
		}
	}()

	return nil
}
