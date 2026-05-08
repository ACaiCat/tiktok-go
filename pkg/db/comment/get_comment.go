package commentdao

import (
	"context"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (c *CommentDao) GetCommentByID(ctx context.Context, commentID int64) (*model.Comment, error) {
	var err error

	comment, err := c.q.Comment.WithContext(ctx).
		Where(c.q.Comment.ID.Eq(commentID)).
		First()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "GetCommentByID failed, commentID: %d", commentID)
	}

	return comment, nil
}

func (c *CommentDao) GetCommentsByVideoID(ctx context.Context, videoID int64, pageSize int, pageNum int) ([]*model.Comment, error) {
	var err error

	comments, err := c.q.Comment.WithContext(ctx).
		Where(c.q.Comment.VideoID.Eq(videoID), c.q.Comment.ParentID.IsNull()).
		Order(c.q.Comment.CreatedAt.Desc()).
		Offset(pageSize * pageNum).
		Limit(pageSize).
		Find()

	if err != nil {
		return nil, errors.Wrapf(err, "GetCommentsByVideoID failed, videoID: %d", videoID)
	}

	return comments, nil
}

func (c *CommentDao) GetCommentsByCommentID(ctx context.Context, commentID int64, pageSize int, pageNum int) ([]*model.Comment, error) {
	var err error

	comments, err := c.q.Comment.WithContext(ctx).
		Where(c.q.Comment.ParentID.Eq(commentID)).
		Order(c.q.Comment.CreatedAt.Desc()).
		Offset(pageSize * pageNum).
		Limit(pageSize).
		Find()

	if err != nil {
		return nil, errors.Wrapf(err, "GetCommentsByCommentID failed, commentID: %d", commentID)
	}

	return comments, nil
}
