package service

import (
	"strconv"

	"github.com/ACaiCat/tiktok-go/biz/model/video"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func (s *VideoService) VisitVideo(req *video.VisitVideoReq) error {
	var err error

	videoID, err := strconv.ParseInt(req.VideoID, 10, 64)

	if err != nil {
		return errno.ParamErr.WithError(err)
	}

	err = s.videoDao.IncrVisitCount(videoID)
	if err != nil {
		return errno.ServiceErr
	}

	return nil
}
