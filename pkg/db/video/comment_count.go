package videodao

import "log"

func (v *VideoDao) IncrCommentCount(videoID int64) error {
	_, err := v.q.Video.
		Where(v.q.Video.ID.Eq(videoID)).
		UpdateColumn(v.q.Video.CommentCount, v.q.Video.CommentCount.Add(1))

	if err != nil {
		log.Printf("failed to increase comment count: %v", err)
		return err
	}

	return nil
}

func (v *VideoDao) DecrCommentCount(videoID int64) error {
	_, err := v.q.Video.
		Where(v.q.Video.ID.Eq(videoID)).
		UpdateColumn(v.q.Video.CommentCount, v.q.Video.CommentCount.Add(-1))

	if err != nil {
		log.Printf("failed to decrease comment count: %v", err)
		return err
	}

	return nil
}
