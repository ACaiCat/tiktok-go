package service

import (
	"log"
	"slices"
	"strconv"

	"github.com/ACaiCat/tiktok-go/pkg/db"
	"gorm.io/gorm"

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

	if actionType != interaction.LikeActionType_ADD && actionType != interaction.LikeActionType_DELETE {
		return errno.NotSupportActionErr
	}

	videoExists, err := s.videoDao.IsVideoExists(videoID)
	if err != nil {
		return errno.ServiceErr
	}
	if !videoExists {
		return errno.VideoNotExistErr
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		isLiked, err := s.userCache.IsVideoLiked(userID, videoID)
		if err != nil {
			likedVideos, err := s.likeDao.WithTx(tx).GetUserLikes(userID)
			if err != nil {
				return errno.ServiceErr
			}
			isLiked = slices.Contains(likedVideos, videoID)

			go func() {
				if err := s.userCache.SetLikeVideos(userID, likedVideos); err != nil {
					log.Println("failed to cache liked videos for userID", userID, ":", err)
				}
			}()

		}

		if actionType == interaction.LikeActionType_ADD {
			if isLiked {
				return errno.LikeAlreadyExistErr
			}

			if err := s.likeDao.WithTx(tx).AddVideoLike(userID, videoID); err != nil {
				return errno.ServiceErr
			}
			if err := s.videoDao.WithTx(tx).IncrLikeCount(videoID); err != nil {
				return errno.ServiceErr
			}

			if err := s.userCache.SetLikeVideo(userID, videoID); err != nil {
				log.Println("failed to cache like video for userID", userID, "and videoID", videoID, ":", err)
			}
		} else {
			if !isLiked {
				return errno.LikeNotExistErr
			}

			if err := s.likeDao.WithTx(tx).DeleteVideoLike(userID, videoID); err != nil {
				return errno.ServiceErr
			}
			if err := s.videoDao.WithTx(tx).DecrLikeCount(videoID); err != nil {
				return errno.ServiceErr
			}

			if err != nil {
				return errno.ServiceErr
			}

			if err := s.userCache.SetUnlikeVideo(userID, videoID); err != nil {
				log.Println("failed to cache unlike video for userID", userID, "and videoID", videoID, ":", err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *InteractionService) likeCommentByID(commentIDStr string, userID int64, actionType interaction.LikeActionType) error {
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		return errno.ParamErr.WithError(err)
	}

	if actionType != interaction.LikeActionType_ADD && actionType != interaction.LikeActionType_DELETE {
		return errno.NotSupportActionErr
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		commentExists, err := s.commentDao.WithTx(tx).IsCommentExists(commentID)
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

		if actionType == interaction.LikeActionType_ADD {
			if isLiked {
				return errno.LikeAlreadyExistErr
			}
			err := db.DB.Transaction(func(tx *gorm.DB) error {
				if err := s.likeDao.WithTx(tx).AddCommentLike(userID, commentID); err != nil {
					return errno.ServiceErr
				}
				if err := s.commentDao.WithTx(tx).IncrLikeCount(commentID); err != nil {
					return errno.ServiceErr
				}
				return nil
			})

			if err != nil {
				return errno.ServiceErr
			}
		} else {
			if !isLiked {
				return errno.LikeNotExistErr
			}

			if err := s.likeDao.WithTx(tx).DeleteCommentLike(userID, commentID); err != nil {
				return errno.ServiceErr
			}
			if err := s.commentDao.WithTx(tx).DecrLikeCount(commentID); err != nil {
				return errno.ServiceErr
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
