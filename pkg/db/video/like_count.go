package videodao

import (
	"context"

	"github.com/pkg/errors"
)

func (v *VideoDao) IncrLikeCount(ctx context.Context, videoID int64) error {
	_, err := v.q.Video.WithContext(ctx).
		Where(v.q.Video.ID.Eq(videoID)).
		UpdateColumn(v.q.Video.LikeCount, v.q.Video.LikeCount.Add(1))

	if err != nil {
		return errors.Wrapf(err, "IncrLikeCount failed, videoID: %d", videoID)
	}

	return nil
}

func (v *VideoDao) DecrLikeCount(ctx context.Context, videoID int64) error {
	_, err := v.q.Video.WithContext(ctx).
		Where(v.q.Video.ID.Eq(videoID)).
		UpdateColumn(v.q.Video.LikeCount, v.q.Video.LikeCount.Add(-1))

	if err != nil {
		return errors.Wrapf(err, "DecrLikeCount failed, videoID: %d", videoID)
	}

	return nil
}
