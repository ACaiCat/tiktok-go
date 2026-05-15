package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"

	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (s *VideoService) getVideosByIDs(videoIDs []int64) ([]*modelDao.Video, error) {
	if len(videoIDs) == 0 {
		return []*modelDao.Video{}, nil
	}

	videos, err := s.videoCache.GetVideos(s.ctx, videoIDs)
	if err != nil {
		videos, err = s.videoDao.GetVideosByIDs(s.ctx, videoIDs)
		if err != nil {
			return nil, errors.WithMessagef(err, "service.getVideosByIDs: db.GetVideosByIDs failed, videoIDs=%v", videoIDs)
		}
		if len(videos) > 0 {
			if err := s.videoCache.SetVideos(context.Background(), videos); err != nil {
				hlog.Errorf("service.getVideosByIDs cache detail write failed: %v", err)
			}
		}
		return videos, nil
	}

	missingIDs := make([]int64, 0)
	for i, video := range videos {
		if video == nil {
			missingIDs = append(missingIDs, videoIDs[i])
		}
	}
	if len(missingIDs) > 0 {
		missingVideos, err := s.videoDao.GetVideosByIDs(s.ctx, missingIDs)
		if err != nil {
			return nil, errors.WithMessagef(err, "service.getVideosByIDs: db.GetVideosByIDs fallback failed, videoIDs=%v", missingIDs)
		}

		missingVideoByID := make(map[int64]*modelDao.Video, len(missingVideos))
		for _, missingVideo := range missingVideos {
			missingVideoByID[missingVideo.ID] = missingVideo
		}

		for i, videoID := range videoIDs {
			if videos[i] == nil {
				if missingVideo, ok := missingVideoByID[videoID]; ok {
					videos[i] = missingVideo
				}
			}
		}

		if len(missingVideos) > 0 {
			if err := s.videoCache.SetVideos(context.Background(), missingVideos); err != nil {
				hlog.Errorf("service.getVideosByIDs cache detail write failed: %v", err)
			}
		}
	}

	for i := 0; i < len(videos); {
		if videos[i] != nil {
			i++
			continue
		}
		videos = append(videos[:i], videos[i+1:]...)
	}

	return videos, nil
}
