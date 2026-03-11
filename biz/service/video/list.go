package service

import (
	"strconv"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/video"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func (s *VideoService) GetVideoList(req *video.ListReq) ([]*model.Video, int64, error) {
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = constants.DefaultVideoPageSize
	}

	pageNum := req.PageNum
	if pageNum < 0 {
		pageNum = 0
	}

	if pageSize > constants.MaxVideoPageSize {
		pageSize = constants.MaxVideoPageSize
	}

	userID, err := strconv.ParseInt(req.UserID, 10, 64)

	if err != nil {
		return nil, 0, errno.ParamErr.WithError(err)
	}

	videosDao, err := s.videoDao.GetVideosByUserID(userID, int(pageSize), int(pageNum))
	if err != nil {
		return nil, 0, errno.ServiceErr
	}

	videos, err := s.GetLikeAndCommentCount(videosDao)
	if err != nil {
		return nil, 0, errno.ServiceErr
	}

	total, err := s.videoDao.GetVideoCountByUserID(userID)
	if err != nil {
		return nil, 0, errno.ServiceErr
	}

	return videos, total, nil

}
