package service

import (
	"strconv"

	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/biz/model/interaction"
	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/video"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
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

	var videosDao []*modelDao.Video
	var total int64

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		videosDao, err = s.videoDao.WithTx(tx).GetVideosByUserID(s.ctx, userID, int(pageSize), int(pageNum))
		if err != nil {
			return errno.ServiceErr
		}

		total, err = s.videoDao.WithTx(tx).GetVideoCountByUserID(s.ctx, userID)
		if err != nil {
			return errno.ServiceErr
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return VideosDaoToDto(videosDao), total, nil
}

func (s *VideoService) GetLikedVideos(req *interaction.ListLikeReq) ([]*model.Video, error) {
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
		return nil, errno.ParamErr.WithError(err)
	}

	likedVideos, err := s.videoDao.GetUserLikeList(s.ctx, userID, int(pageSize), int(pageNum))
	if err != nil {
		return nil, errno.ServiceErr
	}

	videos := VideosDaoToDto(likedVideos)

	return videos, nil
}
