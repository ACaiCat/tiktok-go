package service

import (
	"errors"
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/video"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
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

	popularVideos, err := s.cache.GetPopularVideos()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Println("failed to get popular videos from cache:", err)
		}

		popularVideos, err = s.videoDao.GetPopularVideos(constants.PopularVideoCacheCount, 0)
		if err != nil {
			return nil, errno.ServiceErr
		}

		err = s.cache.SetPopularVideos(popularVideos)
		if err != nil {
			log.Println("failed to cache popular videos:", err)
		}
	}

	if pageSize*pageNum > constants.PopularVideoCacheCount {
		videosDao, err := s.videoDao.GetPopularVideos(int(pageSize), int(pageNum))
		if err != nil {
			return nil, errno.ServiceErr
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
