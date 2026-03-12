package service

import (
	"strconv"

	"github.com/ACaiCat/tiktok-go/biz/model/interaction"
	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func (s *InteractionService) CommentAction(req *interaction.CommentReq, userID int64) error {
	var err error

	if req.VideoID == nil {
		return errno.ParamErr.WithMessage("视频ID不能为空")
	}

	videoID, err := strconv.ParseInt(*req.VideoID, 10, 64)

	if err != nil {
		return errno.ParamErr.WithError(err)
	}

	exist, err := s.videoDao.IsVideoExists(videoID)
	if err != nil {
		return errno.ServiceErr
	}

	if !exist {
		return errno.VideoNotExistErr
	}

	err = s.commentDao.AddVideoComment(userID, videoID, req.Content)
	if err != nil {
		return errno.ServiceErr
	}

	return nil
}

func (s *InteractionService) DeleteComment(req *interaction.DeleteCommentReq, userID int64) error {
	var err error

	commentID, err := strconv.ParseInt(req.CommentID, 10, 64)

	if err != nil {
		return errno.ParamErr.WithError(err)
	}

	comment, err := s.commentDao.GetCommentByID(commentID)

	if err != nil {
		return errno.ServiceErr
	}

	if comment == nil {
		return errno.CommentNotExistErr
	}

	if comment.UserID != userID {
		return errno.CommentNotBelongToUserErr
	}

	err = s.commentDao.DeleteComment(commentID)

	if err != nil {
		return errno.ServiceErr
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

	if req.VideoID == nil {
		return nil, errno.ParamErr.WithMessage("视频ID不能为空")
	}

	videoID, err := strconv.ParseInt(*req.VideoID, 10, 64)

	if err != nil {
		return nil, errno.ParamErr.WithError(err)
	}

	exist, err := s.videoDao.IsVideoExists(videoID)
	if err != nil {
		return nil, errno.ServiceErr
	}

	if !exist {
		return nil, errno.VideoNotExistErr
	}

	comments, err := s.commentDao.GetCommentsByVideoID(videoID, int(pageSize), int(pageNum))
	if err != nil {
		return nil, errno.ServiceErr
	}

	return CommentsDaoToDto(comments), nil

}
