package service

import (
	"errors"
	"log"
	"slices"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/ACaiCat/tiktok-go/biz/model/interaction"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func (s *InteractionService) LikeVideo(req *interaction.LikeReq, userID int64) error {
	if req.VideoID == nil && req.CommentID == nil {
		return errno.ParamErr.WithMessage("视频ID或评论ID不能为空")
	}

	if req.VideoID != nil && req.CommentID != nil {
		return errno.ParamErr.WithMessage("视频ID和评论ID不能同时存在")
	}

	if req.VideoID != nil {
		return s.likeVideoByID(*req.VideoID, userID, req.ActionType)
	}

	return s.likeCommentByID(*req.CommentID, userID, req.ActionType)
}

func (s *InteractionService) likeVideoByID(videoIDStr string, userID int64, actionType interaction.LikeActionType) error {
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
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

	isLiked, err := s.userCache.IsVideoLiked(userID, videoID)
	if err != nil {
		likedVideos, err := s.likeDao.GetUserLikes(userID)
		if err != nil {
			return errno.ServiceErr
		}
		isLiked = slices.Contains(likedVideos, videoID)
		if cacheErr := s.userCache.SetLikeVideos(userID, likedVideos); cacheErr != nil {
			log.Println("failed to cache liked videos for userID", userID, ":", cacheErr)
		}
		if !errors.Is(err, redis.Nil) {
			log.Println("failed to check if video is liked for userID", userID, "and videoID", videoID, ":", err)
		}
	}

	switch actionType {
	case interaction.LikeActionType_ADD:
		if isLiked {
			return errno.LikeAlreadyExistErr
		}
		if err := s.likeDao.AddVideoLike(userID, videoID); err != nil {
			return errno.ServiceErr
		}
		if err := s.userCache.SetLikeVideo(userID, videoID); err != nil {
			log.Println("failed to cache like video for userID", userID, "and videoID", videoID, ":", err)
		}
		return nil
	case interaction.LikeActionType_DELETE:
		if !isLiked {
			return errno.LikeNotExistErr
		}
		if err := s.likeDao.DeleteVideoLike(userID, videoID); err != nil {
			return errno.ServiceErr
		}
		if err := s.userCache.SetUnlikeVideo(userID, videoID); err != nil {
			log.Println("failed to cache unlike video for userID", userID, "and videoID", videoID, ":", err)
		}
		return nil
	}

	return errno.NotSupportActionErr
}

func (s *InteractionService) likeCommentByID(commentIDStr string, userID int64, actionType interaction.LikeActionType) error {
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
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

	switch actionType {
	case interaction.LikeActionType_ADD:
		if isLiked {
			return errno.LikeAlreadyExistErr
		}
		if err := s.likeDao.AddCommentLike(userID, commentID); err != nil {
			return errno.ServiceErr
		}
		return nil
	case interaction.LikeActionType_DELETE:
		if !isLiked {
			return errno.LikeNotExistErr
		}
		if err := s.likeDao.DeleteCommentLike(userID, commentID); err != nil {
			return errno.ServiceErr
		}
		return nil
	}

	return errno.NotSupportActionErr
}
