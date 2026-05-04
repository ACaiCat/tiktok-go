package commentdao

import "log"

func (c *CommentDao) IncrLikeCount(commentID int64) error {
	_, err := c.q.Comment.
		Where(c.q.Comment.ID.Eq(commentID)).
		UpdateColumn(c.q.Comment.LikeCount, c.q.Comment.LikeCount.Add(1))

	if err != nil {
		log.Printf("failed to increase like count: %v", err)
		return err
	}

	return nil
}

func (c *CommentDao) DecrLikeCount(commentID int64) error {
	_, err := c.q.Comment.
		Where(c.q.Comment.ID.Eq(commentID)).
		UpdateColumn(c.q.Comment.LikeCount, c.q.Comment.LikeCount.Add(-1))

	if err != nil {
		log.Printf("failed to decrease like count: %v", err)
		return err
	}

	return nil
}
