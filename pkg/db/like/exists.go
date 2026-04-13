package likedao

import "log"

func (l *LikeDao) IsCommentLikeExists(userID, commentID int64) (bool, error) {
	var err error

	count, err := l.q.Like.
		Select(l.q.Like.ID).
		Where(l.q.Like.UserID.Eq(userID), l.q.Like.CommentID.Eq(commentID)).
		Count()
	if err != nil {
		log.Printf("failed to check if comment like exists for userID %d and commentID %d: %v", userID, commentID, err)
		return false, err
	}

	return count > 0, nil
}
