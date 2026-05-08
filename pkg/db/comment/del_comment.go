package commentdao

import (
	"context"

	"github.com/pkg/errors"
)

func (c *CommentDao) DeleteComment(ctx context.Context, commentID int64) error {
	var err error

	_, err = c.q.Comment.WithContext(ctx).
		Where(c.q.Comment.ID.Eq(commentID)).
		Delete()

	if err != nil {
		return errors.Wrapf(err, "DeleteComment failed, commentID: %d", commentID)
	}

	return nil
}

func (c *CommentDao) DeleteCommentReply(ctx context.Context, commentID int64) error {
	var err error

	_, err = c.q.Comment.WithContext(ctx).
		Where(c.q.Comment.ID.Eq(commentID)).
		Delete()

	if err != nil {
		return errors.Wrapf(err, "DeleteCommentReply failed, commentID: %d", commentID)
	}

	return nil
}
