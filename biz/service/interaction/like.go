package service

import (
	"context"
	"slices"
	"strconv"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/biz/model/interaction"
	"github.com/ACaiCat/tiktok-go/pkg/db"
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
		if err := s.likeVideo(*req.VideoID, userID, req.ActionType); err != nil {
			return errors.WithMessagef(err, "service.LikeVideo: likeVideoByID failed, videoID=%q, userID=%d", *req.VideoID, userID)
		}

		return nil
	}

	if err := s.likeComment(*req.CommentID, userID, req.ActionType); err != nil {
		return errors.WithMessagef(err, "service.LikeVideo: likeCommentByID failed, commentID=%q, userID=%d", *req.CommentID, userID)
	}

	return nil
}

func (s *InteractionService) likeVideo(videoIDStr string, userID int64, actionType interaction.LikeActionType) error {
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		return errno.ParamErr.WithError(err)
	}

	if actionType != interaction.LikeActionType_ADD && actionType != interaction.LikeActionType_DELETE {
		return errno.NotSupportActionErr
	}

	videoExists, err := s.videoDao.IsVideoExists(s.ctx, videoID)
	if err != nil {
		return errors.WithMessagef(err, "service.likeVideo: check video exists failed, videoID=%d, userID=%d", videoID, userID)
	}
	if !videoExists {
		return errno.VideoNotExistErr
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		isLiked, err := s.userCache.IsVideoLiked(s.ctx, userID, videoID)
		if err != nil {
			if !errors.Is(err, redis.Nil) {
				hlog.CtxErrorf(s.ctx, "service.likeVideo: IsVideoLiked failed, userID=%d, videoID=%d, err=%v", userID, videoID, err)
			}

			likedVideos, err := s.likeDao.WithTx(tx).GetUserLikes(s.ctx, userID)
			if err != nil {
				return errors.WithMessagef(err, "service.likeVideo: db.GetUserLikes failed, userID=%d, userID=%d", userID, userID)
			}
			isLiked = slices.Contains(likedVideos, videoID)

			go func() {
				if err := s.userCache.SetLikeVideos(context.Background(), userID, likedVideos); err != nil {
					hlog.Errorf("service.likeVideo: cache.SetLikeVideos failed: %v", err)
				}
			}()
		}

		if actionType == interaction.LikeActionType_ADD {
			if err := s.addVideoLikeTx(tx, userID, videoID, isLiked); err != nil {
				return errors.WithMessagef(err, "service.likeVideoByID: addVideoLikeTx failed, videoID=%d, userID=%d", videoID, userID)
			}
			return nil
		}

		if err := s.cancelVideoLikeTx(tx, userID, videoID, isLiked); err != nil {
			return errors.WithMessagef(err, "service.likeVideoByID: cancelVideoLikeTx failed, videoID=%d, userID=%d", videoID, userID)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *InteractionService) addVideoLikeTx(tx *gorm.DB, userID, videoID int64, isLiked bool) error {
	if isLiked {
		return errno.LikeAlreadyExistErr
	}

	if err := s.likeDao.WithTx(tx).AddVideoLike(s.ctx, userID, videoID); err != nil {
		return errors.WithMessagef(err, "service.addVideoLikeTx: db.AddVideoLike failed, videoID=%d, userID=%d", videoID, userID)
	}
	if err := s.videoDao.WithTx(tx).IncrLikeCount(s.ctx, videoID); err != nil {
		return errors.WithMessagef(err, "service.addVideoLikeTx: db.IncrLikeCount failed, videoID=%d", videoID)
	}

	go func() {
		if err := s.userCache.SetLikeVideo(s.ctx, userID, videoID); err != nil {
			hlog.CtxErrorf(s.ctx, "service.addVideoLikeTx cache update failed: %v", err)
		}
		err := s.videoCache.IncrVideoLikeCount(context.Background(), videoID, 1)
		if err != nil {
			hlog.Errorf("service.addVideoLikeTx: cache.IncrVideoLikeCount failed, videoID=%d, err=%v", videoID, err)
		}
	}()

	return nil
}

func (s *InteractionService) cancelVideoLikeTx(tx *gorm.DB, userID, videoID int64, isLiked bool) error {
	if !isLiked {
		return errno.LikeNotExistErr
	}

	if err := s.likeDao.WithTx(tx).DeleteVideoLike(s.ctx, userID, videoID); err != nil {
		return errors.WithMessagef(err, "service.cancelVideoLikeTx: db.DeleteVideoLike failed, videoID=%d, userID=%d", videoID, userID)
	}

	if err := s.videoDao.WithTx(tx).DecrLikeCount(s.ctx, videoID); err != nil {
		return errors.WithMessagef(err, "service.cancelVideoLikeTx: db.DecrLikeCount failed, videoID=%d", videoID)
	}

	go func() {
		if err := s.userCache.SetUnlikeVideo(s.ctx, userID, videoID); err != nil {
			hlog.Errorf("service.cancelVideoLikeTx: cache.SetUnlikeVideo failed, videoID=%d, err=%v", videoID, err)
		}

		err := s.videoCache.IncrVideoLikeCount(context.Background(), videoID, -1)
		if err != nil {
			hlog.Errorf("service.cancelVideoLikeTx: cache.IncrVideoLikeCount failed, videoID=%d, err=%v", videoID, err)
		}
	}()

	return nil
}

func (s *InteractionService) likeComment(commentIDStr string, userID int64, actionType interaction.LikeActionType) error {
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		return errno.ParamErr.WithError(err)
	}

	if actionType != interaction.LikeActionType_ADD && actionType != interaction.LikeActionType_DELETE {
		return errno.NotSupportActionErr
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		commentExists, err := s.commentDao.WithTx(tx).IsCommentExists(s.ctx, commentID)
		if err != nil {
			return errors.WithMessagef(err, "service.likeCommentByID: check comment exists failed, commentID=%d, userID=%d", commentID, userID)
		}
		if !commentExists {
			return errno.CommentNotExistErr
		}

		isLiked, err := s.likeDao.IsCommentLikeExists(s.ctx, userID, commentID)
		if err != nil {
			return errors.WithMessagef(err, "service.likeCommentByID: check like exists failed, commentID=%d, userID=%d", commentID, userID)
		}

		if actionType == interaction.LikeActionType_ADD {
			return s.addCommentLikeTx(tx, userID, commentID, isLiked)
		}

		return s.cancelCommentLikeTx(tx, userID, commentID, isLiked)
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *InteractionService) addCommentLikeTx(tx *gorm.DB, userID, commentID int64, isLiked bool) error {
	if isLiked {
		return errno.LikeAlreadyExistErr
	}

	if err := s.likeDao.WithTx(tx).AddCommentLike(s.ctx, userID, commentID); err != nil {
		return errors.WithMessagef(err, "service.addCommentLikeTx: db.AddCommentLike failed, commentID=%d, userID=%d", commentID, userID)
	}

	if err := s.commentDao.WithTx(tx).IncrLikeCount(s.ctx, commentID); err != nil {
		return errors.WithMessagef(err, "service.addCommentLikeTx: db.IncrLikeCount failed, commentID=%d", commentID)
	}

	return nil
}

func (s *InteractionService) cancelCommentLikeTx(tx *gorm.DB, userID, commentID int64, isLiked bool) error {
	if !isLiked {
		return errno.LikeNotExistErr
	}

	if err := s.likeDao.WithTx(tx).DeleteCommentLike(s.ctx, userID, commentID); err != nil {
		return errors.WithMessagef(err, "service.cancelCommentLikeTx: db.DeleteCommentLike failed, commentID=%d, userID=%d", commentID, userID)
	}

	if err := s.commentDao.WithTx(tx).DecrLikeCount(s.ctx, commentID); err != nil {
		return errors.WithMessagef(err, "service.cancelCommentLikeTx: db.DecrLikeCount failed, commentID=%d", commentID)
	}

	return nil
}
