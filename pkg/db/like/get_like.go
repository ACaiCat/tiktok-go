package likedao

import (
	"context"

	"github.com/pkg/errors"
)

func (l *LikeDao) GetLikeCounts(ctx context.Context, videoIDs []int64) (map[int64]int64, error) {
	var err error

	type Result struct {
		VideoID int64 `gorm:"column:video_id"`
		Count   int64 `gorm:"column:count"`
	}

	var results []Result

	err = l.q.Like.WithContext(ctx).
		Select(l.q.Like.VideoID, l.q.Like.ID.Count().As("count")).
		Where(l.q.Like.VideoID.In(videoIDs...)).
		Group(l.q.Like.VideoID).
		Scan(&results)

	if err != nil {
		return nil, errors.Wrapf(err, "GetLikeCount failed, videoID: %v, count: %v", videoIDs, l.q.Like.VideoID)
	}

	likeMap := make(map[int64]int64)
	for _, r := range results {
		likeMap[r.VideoID] = r.Count
	}

	return likeMap, nil
}

func (l *LikeDao) GetUserLikes(ctx context.Context, userID int64) ([]int64, error) {
	var err error

	var videoIDs []int64

	err = l.q.Like.WithContext(ctx).
		Select(l.q.Like.VideoID).
		Where(l.q.Like.UserID.Eq(userID), l.q.Like.CommentID.IsNull()).
		Scan(&videoIDs)

	if err != nil {
		return nil, errors.Wrapf(err, "GetUserLikes failed, userID: %d", userID)
	}

	return videoIDs, nil
}

func (l *LikeDao) IsCommentLikeExists(ctx context.Context, userID, commentID int64) (bool, error) {
	var err error

	count, err := l.q.Like.WithContext(ctx).
		Select(l.q.Like.ID).
		Where(l.q.Like.UserID.Eq(userID), l.q.Like.CommentID.Eq(commentID)).
		Count()

	if err != nil {
		return false, errors.Wrapf(err, "IsCommentLikeExists failed, userID: %d, commentID: %d", userID, commentID)
	}

	return count > 0, nil
}
