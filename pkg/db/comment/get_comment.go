package commentdao

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
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

		log.Printf("failed to get comment by ID %d: %v", commentID, err)
		return nil, err
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
		log.Printf("failed to get comments by video ID %d: %v", videoID, err)
		return nil, err
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
		log.Printf("failed to get comments by comment ID %d: %v", commentID, err)
		return nil, err
	}

	return comments, nil
}
