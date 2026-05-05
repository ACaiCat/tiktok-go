package videodao

import (
	"context"
	"errors"
	"log"
	"math/rand/v2"
	"time"

	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (v *VideoDao) GetVideoByID(ctx context.Context, videoID int64) (*model.Video, error) {
	var err error

	video, err := v.q.Video.WithContext(ctx).
		Where(v.q.Video.ID.Eq(videoID)).
		First()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return video, nil
}

func (v *VideoDao) GetFeedByLatestTime(ctx context.Context, latestTime time.Time, limit int) ([]*model.Video, error) {
	var err error

	statement := v.q.Video.WithContext(ctx).Where()

	if !latestTime.IsZero() {
		statement = statement.Where(v.q.Video.CreatedAt.Gt(latestTime))
	}

	videos, err := statement.
		Limit(limit * constants.FetchVideoMultiple).
		Order(v.q.Video.CreatedAt.Desc()).
		Find()

	if err != nil {
		log.Printf("failed to get feed by latest time: %v", err)
		return nil, err
	}

	rand.Shuffle(len(videos), func(i, j int) {
		videos[i], videos[j] = videos[j], videos[i]
	})

	if len(videos) <= limit {
		return videos, nil
	}

	return videos[:limit], nil
}

func (v *VideoDao) GetVideosByUserID(ctx context.Context, userID int64, pageSize int, pageNum int) ([]*model.Video, error) {
	var err error

	videos, err := v.q.Video.WithContext(ctx).
		Where(v.q.Video.UserID.Eq(userID)).
		Offset(pageSize * pageNum).
		Limit(pageSize).
		Find()

	if err != nil {
		log.Printf("failed to get videos by user id: %v", err)
		return nil, err
	}
	return videos, nil
}

func (v *VideoDao) GetPopularVideos(ctx context.Context, pageSize int, pageNum int) ([]*model.Video, error) {
	var err error

	videos, err := v.q.Video.WithContext(ctx).
		Order(v.q.Video.VisitCount.Desc()).
		Offset(pageSize * pageNum).
		Limit(pageSize).
		Find()

	if err != nil {
		log.Printf("failed to get popular videos: %v", err)
		return nil, err
	}
	return videos, nil
}

func (v *VideoDao) GetVideoCountByUserID(ctx context.Context, userID int64) (int64, error) {
	var err error

	count, err := v.q.Video.WithContext(ctx).
		Where(v.q.Video.UserID.Eq(userID)).
		Count()

	if err != nil {
		log.Printf("failed to get video count by user id: %v", err)
		return 0, err
	}
	return count, nil
}

func (v *VideoDao) GetUserLikeList(ctx context.Context, userID int64, pageSize int, pageNum int) ([]*model.Video, error) {
	var err error

	var videoIDs []int64

	err = v.q.Like.WithContext(ctx).
		Select(v.q.Like.VideoID).
		Where(v.q.Like.UserID.Eq(userID)).
		Scan(&videoIDs)

	if err != nil {
		log.Printf("failed to get user likes for userID %d: %v", userID, err)
		return nil, err
	}

	videos, err := v.q.Video.WithContext(ctx).
		Offset(pageSize * pageNum).
		Where(v.q.Video.ID.In(videoIDs...)).
		Limit(pageSize).
		Find()

	if err != nil {
		log.Printf("failed to get user like list: %v\n", err)
		return nil, err
	}

	return videos, nil
}
