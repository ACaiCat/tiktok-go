package commentdao

import (
	"context"
	"log"
)

func (c *CommentDao) DeleteComment(ctx context.Context, commentID int64) error {
	var err error

	_, err = c.q.Comment.WithContext(ctx).
		Where(c.q.Comment.ID.Eq(commentID)).
		Delete()

	if err != nil {
		log.Printf("failed to delete video comment: %v", err)
		return err
	}

	return nil
}

func (c *CommentDao) DeleteCommentReply(ctx context.Context, commentID int64) error {
	var err error

	_, err = c.q.Comment.WithContext(ctx).
		Where(c.q.Comment.ID.Eq(commentID)).
		Delete()

	if err != nil {
		log.Printf("failed to delete comment reply: %v", err)
		return err
	}

	return nil
}
