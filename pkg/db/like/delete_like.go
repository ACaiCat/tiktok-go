package likeDao

import "log"

func (l *LikeDao) DeleteVideoLike(userID, videoID int64) error {
	var err error

	_, err = l.q.Like.Where(l.q.Like.UserID.Eq(userID), l.q.Like.VideoID.Eq(videoID)).Delete()
	if err != nil {
		log.Printf("failed to delete like: %v", err)
		return err
	}

	return nil
}

func (l *LikeDao) DeleteCommentLike(userID, commentID int64) error {
	var err error

	_, err = l.q.Like.Where(l.q.Like.UserID.Eq(userID), l.q.Like.CommentID.Eq(commentID)).Delete()
	if err != nil {
		log.Printf("failed to delete like: %v", err)
		return err
	}

	return nil
}
