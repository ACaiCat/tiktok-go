package service

import (
	"strconv"
	"strings"
	"time"

	"github.com/ACaiCat/tiktok-go/api/model/model"
	"github.com/ACaiCat/tiktok-go/api/model/video"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func (s *VideoService) SearchVideo(req *video.SearchReq) ([]*model.Video, error) {
	var err error

	var keywords []string

	for _, keyword := range strings.Split(req.Keywords, " ") {
		if keyword != "" {
			keywords = append(keywords, keyword)
		}
	}

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

	var toDate time.Time
	var fromDate time.Time

	if req.FromDate != nil {
		fromDateUnixMill, err := strconv.ParseInt(*req.FromDate, 10, 64)
		if err != nil {
			return nil, errno.ParamErr.WithError(err)
		}
		fromDate = time.UnixMilli(fromDateUnixMill)
	}

	if req.ToDate != nil {
		toDateUnixMill, err := strconv.ParseInt(*req.ToDate, 10, 64)
		if err != nil {
			return nil, errno.ParamErr.WithError(err)
		}
		toDate = time.UnixMilli(toDateUnixMill)
	}

	username := ""
	if req.Username != nil {
		username = *req.Username
	}

	videosDao, err := s.videoDao.SearchVideo(keywords, int(pageSize), int(pageNum), fromDate, toDate, username)
	if err != nil {
		return nil, errno.ServiceErr
	}

	videos := VideosDaoToDto(videosDao)

	return videos, nil
}
