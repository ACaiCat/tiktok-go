package videodao

import (
	"context"

	"github.com/pkg/errors"
)

func (v *VideoDao) IncrVisitCount(ctx context.Context, videoID int64) error {
	_, err := v.q.Video.WithContext(ctx).
		Where(v.q.Video.ID.Eq(videoID)).
		UpdateColumn(v.q.Video.VisitCount, v.q.Video.VisitCount.Add(1))

	if err != nil {
		return errors.Wrapf(err, "IncrVisitCount failed, videoID: %d", videoID)
	}

	return nil
}
