package commentdao

import (
	"context"
	"log"
)

func (c *CommentDao) IsCommentExists(ctx context.Context, commentID int64) (bool, error) {
	var err error

	count, err := c.q.Comment.WithContext(ctx).
		Select(c.q.Comment.ID).
		Where(c.q.Comment.ID.Eq(commentID)).
		Count()

	if err != nil {
		log.Printf("failed to check if comment exists for commentID %d: %v", commentID, err)
		return false, err
	}

	return count > 0, nil
}
