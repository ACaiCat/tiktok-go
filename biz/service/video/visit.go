package service

import (
	"strconv"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/biz/model/video"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func (s *VideoService) VisitVideo(req *video.VisitVideoReq) error {
	var err error

	videoID, err := strconv.ParseInt(req.VideoID, 10, 64)

	if err != nil {
		return errno.ParamErr.WithError(err)
	}

	err = s.videoDao.IncrVisitCount(s.ctx, videoID)
	if err != nil {
		return errors.WithMessagef(err, "service.VisitVideo: db.IncrVisitCount failed, videoID=%d", videoID)
	}

	return nil
}
