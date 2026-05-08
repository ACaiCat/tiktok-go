package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/video"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (s *VideoService) GetPopularVideos(req *video.PopularReq) ([]*model.Video, error) {
	var err error

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = constants.DefaultVideoPageSize
	}

	pageNum := req.PageNum
	if pageNum < 0 {
		pageNum = 1
	}

	if pageSize > constants.MaxVideoPageSize {
		pageSize = constants.MaxVideoPageSize
	}

	popularVideos, err := s.cache.GetPopularVideos(s.ctx)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			hlog.CtxErrorf(s.ctx, "service.GetPopularVideos cache read failed: %v", err)
		}

		popularVideos, err = s.videoDao.GetPopularVideos(s.ctx, constants.PopularVideoCacheCount, 0)
		if err != nil {
			return nil, errors.WithMessage(err, "service.GetPopularVideos: db.GetPopularVideos failed")
		}

		go func() {
			err = s.cache.SetPopularVideos(context.Background(), popularVideos)
			if err != nil {
				hlog.Errorf("service.GetPopularVideos cache write failed: %v", err)
			}
		}()
	}

	if pageSize*pageNum > constants.PopularVideoCacheCount {
		videosDao, err := s.videoDao.GetPopularVideos(s.ctx, int(pageSize), int(pageNum))
		if err != nil {
			return nil, errors.WithMessagef(err, "service.GetPopularVideos: db.GetPopularVideos failed, page=%d, pageSize=%d", pageNum, pageSize)
		}
		videos := VideosDaoToDto(videosDao)
		return videos, nil
	}

	targetVideos := make([]*modelDao.Video, 0)

	for i := int(pageSize) * int(pageNum); i < int(pageSize)*(int(pageNum)+1) && i < len(popularVideos); i++ {
		targetVideos = append(targetVideos, popularVideos[i])
	}

	videos := VideosDaoToDto(targetVideos)

	return videos, nil
}
