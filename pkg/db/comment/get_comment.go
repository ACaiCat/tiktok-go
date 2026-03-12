package commentDao

import "log"

func (c *CommentDao) GetCommentCount(videoID int64) (int64, error) {
	var err error

	count, err := c.q.Comment.
		Where(c.q.Comment.VideoID.Eq(videoID)).
		Count()
	if err != nil {
		log.Printf("failed to get comment count for videoID %d: %v", videoID, err)
		return 0, err
	}

	return count, nil
}

func (c *CommentDao) GetCommentCounts(videoIDs []int64) (map[int64]int64, error) {
	var err error

	type Result struct {
		VideoID int64 `gorm:"column:video_id"`
		Count   int64 `gorm:"column:count"`
	}

	var results []Result

	err = c.q.Comment.
		Select(c.q.Comment.VideoID, c.q.Comment.ID.Count().As("count")).
		Where(c.q.Comment.VideoID.In(videoIDs...)).
		Group(c.q.Comment.VideoID).
		Scan(&results)

	if err != nil {
		log.Printf("failed to get comment counts for videoIDs %v: %v", videoIDs, err)
		return nil, err
	}
	commentMap := make(map[int64]int64)
	for _, r := range results {
		commentMap[r.VideoID] = r.Count
	}

	return commentMap, nil
}
