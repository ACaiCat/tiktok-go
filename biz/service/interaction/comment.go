package service

import (
	"strconv"

	"github.com/ACaiCat/tiktok-go/biz/model/interaction"
	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"gorm.io/gorm"
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
			return errno.ServiceErr
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
			return errno.ServiceErr
		}

		return nil
	}
	commentID, err := strconv.ParseInt(*req.CommentID, 10, 64)

	if err != nil {
		return errno.ParamErr.WithError(err)
	}
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		comment, err := s.commentDao.WithTx(tx).GetCommentByID(s.ctx, commentID)
		if err != nil {
			return errno.ServiceErr
		}

		if comment == nil {
			return errno.CommentNotExistErr
		}

		if err = s.commentDao.WithTx(tx).AddCommentReply(s.ctx, userID, comment.VideoID, commentID, req.Content); err != nil {
			return errno.ServiceErr
		}
		if err = s.commentDao.WithTx(tx).IncrCommentCount(s.ctx, commentID); err != nil {
			return errno.ServiceErr
		}
		return nil
	})

	if err != nil {
		return err
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
			return errno.ServiceErr
		}

		if comment == nil {
			return errno.CommentNotExistErr
		}

		if comment.UserID != userID {
			return errno.CommentNotBelongToUserErr
		}

		if comment.ParentID != nil {
			if err := s.commentDao.WithTx(tx).DecrCommentCount(s.ctx, *comment.ParentID); err != nil {
				return errno.ServiceErr
			}
		} else {
			if err := s.videoDao.WithTx(tx).DecrCommentCount(s.ctx, comment.VideoID); err != nil {
				return errno.ServiceErr
			}
		}

		err = s.commentDao.WithTx(tx).DeleteComment(s.ctx, commentID)
		if err != nil {
			return errno.ServiceErr
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *InteractionService) ListComment(req *interaction.ListCommentReq) ([]*model.Comment, error) {
	var err error

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = constants.DefaultCommentPageSize
	}

	pageNum := req.PageNum
	if pageNum < 0 {
		pageNum = 0
	}

	if pageSize > constants.MaxCommentPageSize {
		pageSize = constants.MaxCommentPageSize
	}

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
			return nil, errno.ServiceErr
		}

		if !exist {
			return nil, errno.VideoNotExistErr
		}

		comments, err := s.commentDao.GetCommentsByVideoID(s.ctx, videoID, int(pageSize), int(pageNum))
		if err != nil {
			return nil, errno.ServiceErr
		}

		return CommentsDaoToDto(comments), nil
	}

	commentID, err := strconv.ParseInt(*req.CommentID, 10, 64)

	if err != nil {
		return nil, errno.ParamErr.WithError(err)
	}

	exist, err := s.commentDao.IsCommentExists(s.ctx, commentID)
	if err != nil {
		return nil, errno.ServiceErr
	}

	if !exist {
		return nil, errno.CommentNotExistErr
	}

	comments, err := s.commentDao.GetCommentsByCommentID(s.ctx, commentID, int(pageSize), int(pageNum))
	if err != nil {
		return nil, errno.ServiceErr
	}

	return CommentsDaoToDto(comments), nil
}
