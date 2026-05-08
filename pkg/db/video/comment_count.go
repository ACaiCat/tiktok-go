package videodao

import (
	"context"

	"github.com/pkg/errors"
)

func (v *VideoDao) IncrCommentCount(ctx context.Context, videoID int64) error {
	_, err := v.q.Video.WithContext(ctx).
		Where(v.q.Video.ID.Eq(videoID)).
		UpdateColumn(v.q.Video.CommentCount, v.q.Video.CommentCount.Add(1))

	if err != nil {
		return errors.Wrapf(err, "IncrCommentCount failed, videoID: %d", videoID)
	}

	return nil
}

func (v *VideoDao) DecrCommentCount(ctx context.Context, videoID int64) error {
	_, err := v.q.Video.WithContext(ctx).
		Where(v.q.Video.ID.Eq(videoID)).
		UpdateColumn(v.q.Video.CommentCount, v.q.Video.CommentCount.Add(-1))

	if err != nil {
		return errors.Wrapf(err, "DecrCommentCount failed, videoID: %d", videoID)
	}

	return nil
}
