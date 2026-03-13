package service

import (
	"errors"
	"log"
	"strconv"

	"github.com/ACaiCat/tiktok-go/biz/model/interaction"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/redis/go-redis/v9"
)

func (s *InteractionService) LikeVideo(req *interaction.LikeReq, userID int64) error {
	var err error

	if req.VideoID == nil && req.CommentID == nil {
		return errno.ParamErr.WithMessage("视频ID或评论ID不能为空")
	}

	if req.VideoID != nil && req.CommentID != nil {
		return errno.ParamErr.WithMessage("视频ID和评论ID不能同时存在")
	}

	if req.VideoID != nil {

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
		isLiked := false

		isLiked, err = s.userCache.IsVideoLiked(userID, videoID)
		if err != nil {

			isLiked = false

			likedVideos, err := s.likeDao.GetUserLikes(userID)
			if err != nil {
				return errno.ServiceErr
			}

			for _, id := range likedVideos {
				if videoID == id {
					isLiked = true
					break
				}
			}

			err = s.userCache.SetLikeVideos(userID, likedVideos)
			if err != nil {
				log.Println("failed to cache liked videos for userID", userID, ":", err)
			}

			if !errors.Is(err, redis.Nil) {
				log.Println("failed to check if video is liked for userID", userID, "and videoID", videoID, ":", err)
			}
		}

		switch req.ActionType {
		case interaction.LikeActionType_ADD:
			if isLiked {
				return errno.LikeAlreadyExistErr
			}
			err := s.likeDao.AddVideoLike(userID, videoID)
			if err != nil {
				return errno.ServiceErr
			}
			err = s.userCache.SetLikeVideo(userID, videoID)
			if err != nil {
				log.Println("failed to cache like video for userID", userID, "and videoID", videoID, ":", err)
			}
			return nil
		case interaction.LikeActionType_DELETE:
			if !isLiked {
				return errno.LikeNotExistErr
			}
			err := s.likeDao.DeleteVideoLike(userID, videoID)
			if err != nil {
				return errno.ServiceErr
			}
			err = s.userCache.SetUnlikeVideo(userID, videoID)
			if err != nil {
				log.Println("failed to cache unlike video for userID", userID, "and videoID", videoID, ":", err)
			}

			return nil
		}

		return errno.NotSupportActionErr
	}

	commentID, err := strconv.ParseInt(*req.CommentID, 10, 64)

	if err != nil {
		return errno.ParamErr.WithError(err)
	}

	commentExists, err := s.commentDao.IsCommentExists(commentID)
	if err != nil {
		return errno.ServiceErr
	}

	if !commentExists {
		return errno.CommentNotExistErr
	}

	isLiked, err := s.likeDao.IsCommentLikeExists(userID, commentID)

	if err != nil {
		return errno.ServiceErr
	}

	switch req.ActionType {
	case interaction.LikeActionType_ADD:
		if isLiked {
			return errno.LikeAlreadyExistErr
		}
		err := s.likeDao.AddCommentLike(userID, commentID)
		if err != nil {
			return errno.ServiceErr
		}
		return nil
	case interaction.LikeActionType_DELETE:
		if !isLiked {
			return errno.LikeNotExistErr
		}
		err := s.likeDao.DeleteCommentLike(userID, commentID)
		if err != nil {
			return errno.ServiceErr
		}
		return nil
	}

	return errno.NotSupportActionErr
}
