package likedao

import (
	"context"
	"log"
)

func (l *LikeDao) IsCommentLikeExists(ctx context.Context, userID, commentID int64) (bool, error) {
	var err error

	count, err := l.q.Like.WithContext(ctx).
		Select(l.q.Like.ID).
		Where(l.q.Like.UserID.Eq(userID), l.q.Like.CommentID.Eq(commentID)).
		Count()

	if err != nil {
		log.Printf("failed to check if comment like exists for userID %d and commentID %d: %v", userID, commentID, err)
		return false, err
	}

	return count > 0, nil
}
