package likedao

import (
	"context"

	"github.com/pkg/errors"
)

func (l *LikeDao) DeleteVideoLike(ctx context.Context, userID, videoID int64) error {
	var err error

	_, err = l.q.Like.WithContext(ctx).
		Where(l.q.Like.UserID.Eq(userID), l.q.Like.VideoID.Eq(videoID)).
		Delete()

	if err != nil {
		return errors.Wrapf(err, "DeleteVideoLike failed, userID: %d, videoID: %d", userID, videoID)
	}

	return nil
}

func (l *LikeDao) DeleteCommentLike(ctx context.Context, userID, commentID int64) error {
	var err error

	_, err = l.q.Like.WithContext(ctx).
		Where(l.q.Like.UserID.Eq(userID), l.q.Like.CommentID.Eq(commentID)).
		Delete()

	if err != nil {
		return errors.Wrapf(err, "DeleteCommentLike failed, userID: %d, commentID: %d", userID, commentID)
	}

	return nil
}
