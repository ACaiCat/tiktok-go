package service

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/biz/model/interaction"
	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/ACaiCat/tiktok-go/pkg/utils"
)

func (s *InteractionService) CommentAction(req *interaction.CommentReq, userID int64) error {
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

		exist, err := s.videoDao.IsVideoExists(s.ctx, videoID)
		if err != nil {
			return errors.WithMessagef(err, "service.CommentAction: check video exists failed, videoID=%d, userID=%d", videoID, userID)
		}

		if !exist {
			return errno.VideoNotExistErr
		}

		err = db.DB.Transaction(func(tx *gorm.DB) error {
			if err = s.commentDao.WithTx(tx).AddVideoComment(s.ctx, userID, videoID, req.Content); err != nil {
				return err
			}
			if err = s.videoDao.WithTx(tx).IncrCommentCount(s.ctx, videoID); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			return errors.WithMessagef(err, "service.CommentAction: add video comment tx failed, videoID=%d, userID=%d", videoID, userID)
		}

		go func() {
			err := s.videoCache.IncrVideoCommentCount(context.Background(), videoID, 1)
			if err != nil {
				hlog.Errorf("service.CommentAction: cache.IncrVideoCommentCount failed, videoID=%d, err=%v", videoID, err)
			}
		}()

		return nil
	}
	commentID, err := strconv.ParseInt(*req.CommentID, 10, 64)

	if err != nil {
		return errno.ParamErr.WithError(err)
	}
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		comment, err := s.commentDao.WithTx(tx).GetCommentByID(s.ctx, commentID)
		if err != nil {
			return errors.WithMessagef(err, "service.CommentAction: db.GetCommentByID failed, commentID=%d, userID=%d", commentID, userID)
		}

		if comment == nil {
			return errno.CommentNotExistErr
		}

		if err = s.commentDao.WithTx(tx).AddCommentReply(s.ctx, userID, comment.VideoID, commentID, req.Content); err != nil {
			return err
		}
		if err = s.commentDao.WithTx(tx).IncrCommentCount(s.ctx, commentID); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return errors.WithMessagef(err, "service.CommentAction: reply comment tx failed, commentID=%d, userID=%d", commentID, userID)
	}

	return nil
}

func (s *InteractionService) DeleteComment(req *interaction.DeleteCommentReq, userID int64) error {
	var err error

	commentID, err := strconv.ParseInt(req.CommentID, 10, 64)

	if err != nil {
		return errno.ParamErr.WithError(err)
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		comment, err := s.commentDao.WithTx(tx).GetCommentByID(s.ctx, commentID)

		if err != nil {
			return errors.WithMessagef(err, "service.DeleteComment: db.GetCommentByID failed, commentID=%d, userID=%d", commentID, userID)
		}

		if comment == nil {
			return errno.CommentNotExistErr
		}

		if comment.UserID != userID {
			return errno.CommentNotBelongToUserErr
		}

		if comment.ParentID != nil {
			if err := s.commentDao.WithTx(tx).DecrCommentCount(s.ctx, *comment.ParentID); err != nil {
				return errors.WithMessagef(err, "service.DeleteComment: decr parent comment count failed, commentID=%d", commentID)
			}
		} else {
			if err := s.videoDao.WithTx(tx).DecrCommentCount(s.ctx, comment.VideoID); err != nil {
				return errors.WithMessagef(err, "service.DeleteComment: decr video comment count failed, videoID=%d", comment.VideoID)
			}

			go func() {
				err := s.videoCache.IncrVideoCommentCount(context.Background(), comment.VideoID, -1)
				if err != nil {
					hlog.Errorf("service.DeleteComment: cache.IncrVideoCommentCount failed, videoID=%d, err=%v", comment.VideoID, err)
				}
			}()

		}

		err = s.commentDao.WithTx(tx).DeleteComment(s.ctx, commentID)
		if err != nil {
			return errors.WithMessagef(err, "service.DeleteComment: db.DeleteComment failed, commentID=%d", commentID)
		}

		return nil
	})
	if err != nil {
		return errors.WithMessagef(err, "service.DeleteComment: tx failed, commentID=%d, userID=%d", commentID, userID)
	}

	return nil
}

func (s *InteractionService) ListComment(req *interaction.ListCommentReq) ([]*model.Comment, error) {
	var err error

	pageSize, pageNum := utils.NormalizePage(req.PageSize, req.PageNum, constants.DefaultCommentPageSize, constants.MaxCommentPageSize)

	if req.VideoID == nil && req.CommentID == nil {
		return nil, errno.ParamErr.WithMessage("视频ID或评论ID不能为空")
	}

	if req.VideoID != nil && req.CommentID != nil {
		return nil, errno.ParamErr.WithMessage("视频ID和评论ID不能同时存在")
	}

	if req.VideoID != nil {
		videoID, err := strconv.ParseInt(*req.VideoID, 10, 64)

		if err != nil {
			return nil, errno.ParamErr.WithError(err)
		}

		exist, err := s.videoDao.IsVideoExists(s.ctx, videoID)
		if err != nil {
			return nil, errors.WithMessagef(err, "service.ListComment: check video exists failed, videoID=%d", videoID)
		}

		if !exist {
			return nil, errno.VideoNotExistErr
		}

		comments, err := s.commentDao.GetCommentsByVideoID(s.ctx, videoID, pageSize, pageNum)
		if err != nil {
			return nil, errors.WithMessagef(err, "service.ListComment: db.GetCommentsByVideoID failed, videoID=%d", videoID)
		}

		return CommentsDaoToDto(comments), nil
	}

	commentID, err := strconv.ParseInt(*req.CommentID, 10, 64)

	if err != nil {
		return nil, errno.ParamErr.WithError(err)
	}

	exist, err := s.commentDao.IsCommentExists(s.ctx, commentID)
	if err != nil {
		return nil, errors.WithMessagef(err, "service.ListComment: check comment exists failed, commentID=%d", commentID)
	}

	if !exist {
		return nil, errno.CommentNotExistErr
	}

	comments, err := s.commentDao.GetCommentsByCommentID(s.ctx, commentID, pageSize, pageNum)
	if err != nil {
		return nil, errors.WithMessagef(err, "service.ListComment: db.GetCommentsByCommentID failed, commentID=%d", commentID)
	}

	return CommentsDaoToDto(comments), nil
}
