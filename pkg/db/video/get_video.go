package videodao

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/pkg/errors"
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
		return nil, errors.Wrapf(err, "GetVideoByID failed, videoID: %d", videoID)
	}
	return video, nil
}

func (v *VideoDao) GetVideosByIDs(ctx context.Context, videoIDs []int64) ([]*model.Video, error) {
	if len(videoIDs) == 0 {
		return []*model.Video{}, nil
	}

	videos, err := v.q.Video.WithContext(ctx).
		Where(v.q.Video.ID.In(videoIDs...)).
		Find()

	if err != nil {
		return nil, errors.Wrapf(err, "GetVideosByIDs failed, videoIDs: %v", videoIDs)
	}

	videoByID := make(map[int64]*model.Video, len(videos))
	for _, video := range videos {
		videoByID[video.ID] = video
	}

	orderedVideos := make([]*model.Video, 0, len(videos))
	for _, videoID := range videoIDs {
		if video, ok := videoByID[videoID]; ok {
			orderedVideos = append(orderedVideos, video)
		}
	}

	return orderedVideos, nil
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
		return nil, errors.Wrapf(err, "GetFeedByLatestTime failed, latestTime=%s, limit=%d", latestTime, limit)
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
		return nil, errors.Wrapf(err, "GetVideosByUserID failed, userID: %d", userID)
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
		return nil, errors.Wrapf(err, "GetPopularVideos failed")
	}
	return videos, nil
}

func (v *VideoDao) GetVideoCountByUserID(ctx context.Context, userID int64) (int64, error) {
	var err error

	count, err := v.q.Video.WithContext(ctx).
		Where(v.q.Video.UserID.Eq(userID)).
		Count()

	if err != nil {
		return 0, errors.Wrapf(err, "GetVideoCountByUserID failed, userID: %d", userID)
	}
	return count, nil
}

func (v *VideoDao) GetUserLikeList(ctx context.Context, userID int64, pageSize int, pageNum int) ([]*model.Video, error) {
	var err error

	var videoIDs []int64

	err = v.q.Like.WithContext(ctx).
		Select(v.q.Like.VideoID).
		Where(v.q.Like.UserID.Eq(userID), v.q.Like.VideoID.IsNotNull()).
		Scan(&videoIDs)

	if err != nil {
		return nil, errors.Wrapf(err, "GetUserLikeList failed, userID: %d", userID)
	}

	videos, err := v.q.Video.WithContext(ctx).
		Offset(pageSize * pageNum).
		Where(v.q.Video.ID.In(videoIDs...)).
		Limit(pageSize).
		Find()

	if err != nil {
		return nil, errors.Wrapf(err, "GetUserLikeList failed, userID: %d", userID)
	}

	return videos, nil
}
