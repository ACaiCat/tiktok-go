package commentdao

import (
	"context"

	"github.com/pkg/errors"
)

func (c *CommentDao) IsCommentExists(ctx context.Context, commentID int64) (bool, error) {
	var err error

	count, err := c.q.Comment.WithContext(ctx).
		Select(c.q.Comment.ID).
		Where(c.q.Comment.ID.Eq(commentID)).
		Count()

	if err != nil {
		return false, errors.Wrapf(err, "IsCommentExists failed, commentID: %d", commentID)
	}

	return count > 0, nil
}
