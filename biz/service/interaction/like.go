package service

import (
	"strconv"

	"github.com/ACaiCat/tiktok-go/biz/model/interaction"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func (s *InteractionService) LikeVideo(req *interaction.LikeReq, userID int64) error {
	var err error

	if req.VideoID == nil {
		return errno.ParamErr.WithMessage("视频ID不能为空")
	}

	videoID, err := strconv.ParseInt(*req.VideoID, 10, 64)

	if err != nil {
		return errno.ParamErr.WithError(err)
	}

	videoExists, err := s.videoDao.IsVideoExists(videoID)
	if err != nil {
		return errno.ServiceErr
	}

	if !videoExists {
		return errno.VideoNotExistErr
	}

	likedVideos, err := s.likeDao.GetUserLikes(userID)
	if err != nil {
		return errno.ServiceErr
	}

	switch req.ActionType {
	case interaction.LikeActionType_ADD:
		for _, id := range likedVideos {
			if videoID == id {
				return errno.LikeAlreadyExistErr
			}
		}

		err := s.likeDao.AddVideoLike(userID, videoID)
		if err != nil {
			return errno.ServiceErr
		}

		return nil
	case interaction.LikeActionType_DELETE:
		for _, id := range likedVideos {
			if videoID == id {
				err := s.likeDao.DeleteVideoLike(userID, videoID)
				if err != nil {
					return errno.ServiceErr
				}
				return nil
			}
		}

		return errno.LikeNotExistErr
	}

	return errno.NotSupportActionErr
}
