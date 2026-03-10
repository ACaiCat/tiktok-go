package service

import (
	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/video"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
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

	videosDao, err := s.videoDao.GetPopularVideos(int(req.PageSize), int(req.PageNum))
	if err != nil {
		return nil, err
	}

	videos, err := s.GetLikeAndCommentCount(videosDao)
	if err != nil {
		return nil, err
	}

	return videos, nil

}
