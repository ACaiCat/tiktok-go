package likeDao

import "log"

func (c *LikeDao) GetLikeCount(videoID int64) (int64, error) {
	var err error

	count, err := c.q.Like.Where(c.q.Like.VideoID.Eq(videoID)).Count()
	if err != nil {
		log.Printf("failed to get like count for videoID %d: %v", videoID, err)
		return 0, err
	}

	return count, nil
}

func (c *LikeDao) GetLikeCounts(videoIDs []int64) (map[int64]int64, error) {
	var err error

	var result map[int64]int64

	err = c.q.Like.Select(c.q.Like.VideoID, c.q.Like.ID.Count().As("count")).
		Where(c.q.Like.VideoID.
			In(videoIDs...)).
		Group(c.q.Like.VideoID).
		Scan(&result)

	if err != nil {
		log.Printf("failed to get like counts for videoIDs %v: %v", videoIDs, err)
		return nil, err
	}

	return result, nil
}
