package likedao

import (
	"context"

	"github.com/pkg/errors"
)

func (l *LikeDao) IsCommentLikeExists(ctx context.Context, userID, commentID int64) (bool, error) {
	var err error

	count, err := l.q.Like.WithContext(ctx).
		Select(l.q.Like.ID).
		Where(l.q.Like.UserID.Eq(userID), l.q.Like.CommentID.Eq(commentID)).
		Count()

	if err != nil {
		return false, errors.Wrapf(err, "IsCommentLikeExists failed, userID: %d, commentID: %d", userID, commentID)
	}

	return count > 0, nil
}
