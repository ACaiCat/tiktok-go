package videodao

import "log"

func (v *VideoDao) IsVideoExists(videoID int64) (bool, error) {
	var err error

	count, err := v.q.Video.
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
