package videodao

import (
	"context"
	"log"
)

func (v *VideoDao) IncrCommentCount(ctx context.Context, videoID int64) error {
	_, err := v.q.Video.WithContext(ctx).
		Where(v.q.Video.ID.Eq(videoID)).
		UpdateColumn(v.q.Video.CommentCount, v.q.Video.CommentCount.Add(1))

	if err != nil {
		log.Printf("failed to increase comment count: %v", err)
		return err
	}

	return nil
}

func (v *VideoDao) DecrCommentCount(ctx context.Context, videoID int64) error {
	_, err := v.q.Video.WithContext(ctx).
		Where(v.q.Video.ID.Eq(videoID)).
		UpdateColumn(v.q.Video.CommentCount, v.q.Video.CommentCount.Add(-1))

	if err != nil {
		log.Printf("failed to decrease comment count: %v", err)
		return err
	}

	return nil
}
