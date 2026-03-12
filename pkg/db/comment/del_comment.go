package commentDao

import "log"

func (c *CommentDao) DeleteComment(commentID int64) error {
	var err error

	_, err = c.q.Comment.Where(c.q.Comment.ID.Eq(commentID)).Delete()
	if err != nil {
		log.Printf("failed to delete video comment: %v", err)
		return err
	}

	return nil
}
