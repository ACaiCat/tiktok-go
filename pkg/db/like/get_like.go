package likedao

import "log"

func (l *LikeDao) GetLikeCounts(videoIDs []int64) (map[int64]int64, error) {
	var err error

	type Result struct {
		VideoID int64 `gorm:"column:video_id"`
		Count   int64 `gorm:"column:count"`
	}

	var results []Result

	err = l.q.Like.
		Select(l.q.Like.VideoID, l.q.Like.ID.Count().As("count")).
		Where(l.q.Like.VideoID.In(videoIDs...)).
		Group(l.q.Like.VideoID).
		Scan(&results)

	if err != nil {
		log.Printf("failed to get like counts for videoIDs %v: %v", videoIDs, err)
		return nil, err
	}

	likeMap := make(map[int64]int64)
	for _, r := range results {
		likeMap[r.VideoID] = r.Count
	}

	return likeMap, nil
}

func (l *LikeDao) GetUserLikes(userID int64) ([]int64, error) {
	var err error

	var videoIDs []int64

	err = l.q.Like.
		Select(l.q.Like.VideoID).
		Where(l.q.Like.UserID.Eq(userID)).
		Scan(&videoIDs)

	if err != nil {
		log.Printf("failed to get user likes for userID %d: %v", userID, err)
		return nil, err
	}

	return videoIDs, nil
}
