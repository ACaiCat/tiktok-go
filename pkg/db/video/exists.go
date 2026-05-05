package videodao

import (
	"context"
	"log"
)

func (v *VideoDao) IsVideoExists(ctx context.Context, videoID int64) (bool, error) {
	var err error

	count, err := v.q.Video.WithContext(ctx).
		Select(v.q.Video.ID).
		Where(v.q.Video.ID.Eq(videoID)).
		Limit(1).
		Count()

	if err != nil {
		log.Printf("failed to check if video exists for videoID %d: %v", videoID, err)
		return false, err
	}

	return count > 0, nil
}
