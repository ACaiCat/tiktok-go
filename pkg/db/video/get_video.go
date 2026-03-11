package videoDao

import (
	"errors"
	"log"
	"math/rand/v2"
	"time"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
	"gorm.io/gorm"
)

func (v *VideoDao) GetVideoByID(videoID int64) (*model.Video, error) {
	var err error
	video, err := v.q.Video.
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

func (v *VideoDao) GetFeedByLatestTime(latestTime time.Time, limit int) ([]*model.Video, error) {
	var err error

	statement := v.q.Video.Where()

	if !latestTime.IsZero() {
		statement = statement.Where(v.q.Video.CreatedAt.Gt(latestTime))
	}

	videos, err := statement.
		Limit(limit * 3).
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

func (v *VideoDao) GetVideosByUserID(userID int64, pageSize int, pageNum int) ([]*model.Video, error) {
	var err error
	videos, err := v.q.Video.
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

func (v *VideoDao) GetPopularVideos(pageSize int, pageNum int) ([]*model.Video, error) {
	var err error
	videos, err := v.q.Video.
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

func (v *VideoDao) GetVideoCountByUserID(userID int64) (int64, error) {
	var err error
	count, err := v.q.Video.
		Where(v.q.Video.UserID.Eq(userID)).
		Count()
	if err != nil {
		log.Printf("failed to get video count by user id: %v", err)
		return 0, err
	}
	return count, nil
}
