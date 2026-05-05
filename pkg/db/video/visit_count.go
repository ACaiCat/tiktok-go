package videodao

import (
	"context"
	"log"
)

func (v *VideoDao) IncrVisitCount(ctx context.Context, videoID int64) error {
	_, err := v.q.Video.WithContext(ctx).
		Where(v.q.Video.ID.Eq(videoID)).
		UpdateColumn(v.q.Video.VisitCount, v.q.Video.VisitCount.Add(1))

	if err != nil {
		log.Printf("failed to increase visit count: %v", err)
		return err
	}

	return nil
}
