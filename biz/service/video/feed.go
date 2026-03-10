package service

import (
	"strconv"
	"time"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/video"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func (s *VideoService) GetFeed(req *video.FeedReq) ([]*model.Video, error) {
	var err error

	latestTime := time.UnixMilli(0)
	if req.LatestTime != nil && *req.LatestTime != "" {
		unixMilliStamp, err := strconv.ParseInt(*req.LatestTime, 10, 64)
		if err != nil {
			return nil, errno.ParamErr.WithError(err)
		}

		latestTime = time.UnixMilli(unixMilliStamp)
	}

	videosDao, err := s.videoDao.GetFeedByLatestTime(latestTime, constants.FeedCount)
	if err != nil {
		return nil, errno.ServiceErr
	}

	videoIDs := make([]int64, len(videosDao))
	for i, v := range videosDao {
		videoIDs[i] = v.ID
	}

	idWithCommentCount, err := s.commentDao.GetCommentCounts(videoIDs)
	if err != nil {
		return nil, errno.ServiceErr
	}

	idWithLikeCount, err := s.likeDao.GetLikeCounts(videoIDs)
	if err != nil {
		return nil, errno.ServiceErr
	}

	videos := VideosDaoToDto(videosDao, idWithLikeCount, idWithCommentCount)

	return videos, nil

}
