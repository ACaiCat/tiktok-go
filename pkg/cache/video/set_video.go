package videocache

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (p *VideoCache) SetVideos(ctx context.Context, videos []*model.Video) error {
	pipe := p.c.Pipeline()
	for _, video := range videos {
		pipe.HSet(ctx, getVideoKey(video.ID), map[string]any{
			"id":            video.ID,
			"user_id":       video.UserID,
			"video_url":     video.VideoURL,
			"cover_url":     video.CoverURL,
			"title":         video.Title,
			"description":   video.Description,
			"visit_count":   video.VisitCount,
			"like_count":    video.LikeCount,
			"comment_count": video.CommentCount,
			"created_at":    video.CreatedAt.UnixMilli(),
		})
		pipe.Expire(ctx, getVideoKey(video.ID), constants.VideoCacheExpiration)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return errors.Wrapf(err, "SetVideos failed, videos=%v", videos)
	}

	return nil
}

func (p *VideoCache) SetVideo(ctx context.Context, video *model.Video) error {
	pipe := p.c.Pipeline()
	pipe.HSet(ctx, getVideoKey(video.ID), map[string]any{
		"id":            video.ID,
		"user_id":       video.UserID,
		"video_url":     video.VideoURL,
		"cover_url":     video.CoverURL,
		"title":         video.Title,
		"description":   video.Description,
		"visit_count":   video.VisitCount,
		"like_count":    video.LikeCount,
		"comment_count": video.CommentCount,
		"created_at":    video.CreatedAt.UnixMilli(),
	})
	pipe.Expire(ctx, getVideoKey(video.ID), constants.VideoCacheExpiration)

	if _, err := pipe.Exec(ctx); err != nil {
		return errors.Wrapf(err, "SetVideo failed, video=%v", video)
	}

	return nil
}

func (p *VideoCache) IncrVideoVisitCount(ctx context.Context, videoID int64) error {
	if err := p.c.HIncrBy(ctx, getVideoKey(videoID), "visit_count", 1).Err(); err != nil {
		return errors.Wrapf(err, "IncrVideoVisitCount failed, video=%v", videoID)
	}

	return nil
}

func (p *VideoCache) IncrVideoLikeCount(ctx context.Context, videoID int64, count int64) error {
	if err := p.c.HIncrBy(ctx, getVideoKey(videoID), "like_count", count).Err(); err != nil {
		return errors.Wrapf(err, "IncrVideoLikeCount failed, video=%v", videoID)
	}

	return nil
}

func (p *VideoCache) IncrVideoCommentCount(ctx context.Context, videoID int64, count int64) error {
	if err := p.c.HIncrBy(ctx, getVideoKey(videoID), "comment_count", count).Err(); err != nil {
		return errors.Wrapf(err, "IncrVideoCommentCount failed, video=%v", videoID)
	}
	return nil
}

func (p *VideoCache) CleanVideoCache(ctx context.Context, videoID int64) error {
	if err := p.c.Del(ctx, getVideoKey(videoID)).Err(); err != nil {
		return errors.Wrapf(err, "CleanVideoCache failed, video=%v", videoID)
	}
	return nil
}
