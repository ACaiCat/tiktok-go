package videodao

import (
	"context"
	"log"
)

func (v *VideoDao) IncrLikeCount(ctx context.Context, videoID int64) error {
	_, err := v.q.Video.WithContext(ctx).
		Where(v.q.Video.ID.Eq(videoID)).
		UpdateColumn(v.q.Video.LikeCount, v.q.Video.LikeCount.Add(1))

	if err != nil {
		log.Printf("failed to increase like count: %v", err)
		return err
	}

	return nil
}

func (v *VideoDao) DecrLikeCount(ctx context.Context, videoID int64) error {
	_, err := v.q.Video.WithContext(ctx).
		Where(v.q.Video.ID.Eq(videoID)).
		UpdateColumn(v.q.Video.LikeCount, v.q.Video.LikeCount.Add(-1))

	if err != nil {
		log.Printf("failed to decrease like count: %v", err)
		return err
	}

	return nil
}
