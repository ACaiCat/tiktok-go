package commentdao

import (
	"context"

	"github.com/pkg/errors"
)

func (c *CommentDao) IncrCommentCount(ctx context.Context, commentID int64) error {
	_, err := c.q.Comment.WithContext(ctx).
		Where(c.q.Comment.ID.Eq(commentID)).
		UpdateColumn(c.q.Comment.CommentCount, c.q.Comment.CommentCount.Add(1))

	if err != nil {
		return errors.Wrapf(err, "IncrCommentCount failed, commentID: %d", commentID)
	}

	return nil
}

func (c *CommentDao) DecrCommentCount(ctx context.Context, commentID int64) error {
	_, err := c.q.Comment.WithContext(ctx).
		Where(c.q.Comment.ID.Eq(commentID)).
		UpdateColumn(c.q.Comment.CommentCount, c.q.Comment.CommentCount.Add(-1))

	if err != nil {
		return errors.Wrapf(err, "DecrCommentCount failed, commentID: %d, err: %v", commentID, err)
	}

	return nil
}
