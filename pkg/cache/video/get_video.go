package videocache

import (
	"context"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (p *VideoCache) GetVideo(ctx context.Context, videoID int64) (*model.Video, error) {
	data, err := p.c.HGetAll(ctx, getVideoKey(videoID)).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "GetVideo failed, videoID=%d", videoID)
	}

	video, err := buildVideoFromHash(data)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, err
		}
		return nil, errors.Wrapf(err, "GetVideo decode failed, videoID=%d", videoID)
	}

	return video, nil
}

func (p *VideoCache) GetVideos(ctx context.Context, videoIDs []int64) ([]*model.Video, error) {
	if len(videoIDs) == 0 {
		return []*model.Video{}, nil
	}

	pipe := p.c.Pipeline()
	cmd := make([]*redis.MapStringStringCmd, len(videoIDs))
	for i, videoID := range videoIDs {
		cmd[i] = pipe.HGetAll(ctx, getVideoKey(videoID))
	}

	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, "GetVideos failed, videoIDs=%v", videoIDs)
	}

	videos := make([]*model.Video, len(videoIDs))
	for i, cmd := range cmd {
		data, err := cmd.Result()
		if err != nil {
			return nil, errors.Wrapf(err, "GetVideos read failed, videoID=%d", videoIDs[i])
		}

		video, err := buildVideoFromHash(data)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			return nil, errors.Wrapf(err, "GetVideos decode failed, videoID=%d", videoIDs[i])
		}

		videos[i] = video
	}

	return videos, nil
}

func buildVideoFromHash(data map[string]string) (*model.Video, error) {
	if len(data) == 0 {
		return nil, redis.Nil
	}

	id, err := parseInt64Field(data, "id")
	if err != nil {
		return nil, err
	}

	userID, err := parseInt64Field(data, "user_id")
	if err != nil {
		return nil, err
	}

	visitCount, err := parseInt64Field(data, "visit_count")
	if err != nil {
		return nil, err
	}

	likeCount, err := parseInt64Field(data, "like_count")
	if err != nil {
		return nil, err
	}

	commentCount, err := parseInt64Field(data, "comment_count")
	if err != nil {
		return nil, err
	}

	createdAtUnixMilli, err := parseInt64Field(data, "created_at")
	if err != nil {
		return nil, err
	}

	return &model.Video{
		ID:           id,
		UserID:       userID,
		VideoURL:     data["video_url"],
		CoverURL:     data["cover_url"],
		Title:        data["title"],
		Description:  data["description"],
		VisitCount:   visitCount,
		LikeCount:    likeCount,
		CommentCount: commentCount,
		CreatedAt:    time.UnixMilli(createdAtUnixMilli),
	}, nil
}

func parseInt64Field(data map[string]string, field string) (int64, error) {
	value, ok := data[field]
	if !ok {
		return 0, errors.Errorf("field %s missing", field)
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "field %s invalid", field)
	}

	return parsed, nil
}
