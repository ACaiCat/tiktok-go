package commentdao

import "log"

func (c *CommentDao) IncrCommentCount(commentID int64) error {
	_, err := c.q.Comment.
		Where(c.q.Comment.ID.Eq(commentID)).
		UpdateColumn(c.q.Comment.CommentCount, c.q.Comment.CommentCount.Add(1))

	if err != nil {
		log.Printf("failed to increase comment count: %v", err)
		return err
	}

	return nil
}

func (c *CommentDao) DecrCommentCount(commentID int64) error {
	_, err := c.q.Comment.
		Where(c.q.Comment.ID.Eq(commentID)).
		UpdateColumn(c.q.Comment.CommentCount, c.q.Comment.CommentCount.Add(-1))

	if err != nil {
		log.Printf("failed to decrease comment count: %v", err)
		return err
	}

	return nil
}
