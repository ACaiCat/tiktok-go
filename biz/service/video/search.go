package service

import (
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/video"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/ACaiCat/tiktok-go/pkg/utils"
)

func (s *VideoService) SearchVideo(req *video.SearchReq) ([]*model.Video, error) {
	var err error

	var keywords []string

	for _, keyword := range strings.Split(req.Keywords, " ") {
		if keyword != "" {
			keywords = append(keywords, keyword)
		}
	}

	pageSize, pageNum := utils.NormalizePage(req.PageSize, req.PageNum, constants.DefaultVideoPageSize, constants.MaxVideoPageSize)

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

	videosDao, err := s.videoDao.SearchVideo(s.ctx, keywords, pageSize, pageNum, fromDate, toDate, username)
	if err != nil {
		return nil, errors.WithMessagef(err, "service.SearchVideo: db.SearchVideo failed, keywords=%q", keywords)
	}

	videos := VideosDaoToDto(videosDao)

	return videos, nil
}
