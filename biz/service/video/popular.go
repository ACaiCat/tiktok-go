package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/video"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/utils"
)

func (s *VideoService) GetPopularVideos(req *video.PopularReq) ([]*model.Video, error) {
	pageSize, pageNum := utils.NormalizePage(req.PageSize, req.PageNum, constants.DefaultVideoPageSize, constants.MaxVideoPageSize)

	start := pageSize * pageNum
	if start+pageSize > constants.PopularVideoCacheCount {
		videosDao, err := s.videoDao.GetPopularVideos(s.ctx, pageSize, pageNum)
		if err != nil {
			return nil, errors.WithMessagef(err, "service.GetPopularVideos: db.GetPopularVideos failed, page=%d, pageSize=%d", pageNum, pageSize)
		}
		return VideosDaoToDto(videosDao), nil
	}

	popularVideoIDs, err := s.videoCache.GetPopularVideos(s.ctx, pageSize, pageNum)
	if err == nil {
		targetVideos, err := s.getVideosByIDs(popularVideoIDs)
		if err != nil {
			return nil, err
		}
		return VideosDaoToDto(targetVideos), nil
	}

	if !errors.Is(err, redis.Nil) {
		hlog.CtxErrorf(s.ctx, "service.GetPopularVideos cache read failed: %v", err)
	}

	videosDao, err := s.videoDao.GetPopularVideos(s.ctx, pageSize, pageNum)
	if err != nil {
		return nil, errors.WithMessagef(err, "service.GetPopularVideos: db.GetPopularVideos failed, page=%d, pageSize=%d", pageNum, pageSize)
	}

	go func() {
		popularVideos, err := s.videoDao.GetPopularVideos(context.Background(), constants.PopularVideoCacheCount, 0)
		if err != nil {
			hlog.Errorf("service.GetPopularVideos refresh cache db failed: %v", err)
		} else if err := s.videoCache.SetPopularVideos(context.Background(), popularVideos); err != nil {
			hlog.Errorf("service.GetPopularVideos cache write failed: %v", err)
		}
		if err := s.videoCache.SetVideos(context.Background(), videosDao); err != nil {
			hlog.Errorf("service.GetPopularVideos cache detail write failed: %v", err)
		}
	}()

	return VideosDaoToDto(videosDao), nil
}
