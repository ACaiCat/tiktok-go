package videodao

import (
	"context"

	"github.com/pkg/errors"
)

func (v *VideoDao) IsVideoExists(ctx context.Context, videoID int64) (bool, error) {
	var err error

	count, err := v.q.Video.WithContext(ctx).
		Select(v.q.Video.ID).
		Where(v.q.Video.ID.Eq(videoID)).
		Limit(1).
		Count()

	if err != nil {
		return false, errors.Wrapf(err, "IsVideoExists failed, videoID: %d", videoID)
	}

	return count > 0, nil
}
