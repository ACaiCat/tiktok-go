package commentdao

import (
	"context"

	"github.com/pkg/errors"
)

func (c *CommentDao) IncrLikeCount(ctx context.Context, commentID int64) error {
	_, err := c.q.Comment.WithContext(ctx).
		Where(c.q.Comment.ID.Eq(commentID)).
		UpdateColumn(c.q.Comment.LikeCount, c.q.Comment.LikeCount.Add(1))

	if err != nil {
		return errors.Wrapf(err, "IncrLikeCount failed, commentID: %d", commentID)
	}

	return nil
}

func (c *CommentDao) DecrLikeCount(ctx context.Context, commentID int64) error {
	_, err := c.q.Comment.WithContext(ctx).
		Where(c.q.Comment.ID.Eq(commentID)).
		UpdateColumn(c.q.Comment.LikeCount, c.q.Comment.LikeCount.Add(-1))

	if err != nil {
		return errors.Wrapf(err, "DecrLikeCount failed, commentID: %d", commentID)
	}

	return nil
}
