package videocache

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (p *VideoCache) SetPopularVideos(ctx context.Context, videos []*model.Video) error {
	key := getPopularVideoKey()
	pipe := p.c.Pipeline()
	pipe.Del(ctx, key)

	if len(videos) > 0 {
		members := make([]redis.Z, 0, len(videos))
		for _, video := range videos {
			members = append(members, redis.Z{
				Score:  float64(video.VisitCount),
				Member: video.ID,
			})
		}
		pipe.ZAdd(ctx, key, members...)
		pipe.Expire(ctx, key, constants.PopularVideoCacheExpiration)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return errors.Wrapf(err, "SetPopularVideosCache failed, %v", videos)
	}

	return nil
}

func (p *VideoCache) IncrPopularVideoVisitCount(ctx context.Context, videoID int64) error {
	if err := p.c.ZIncrBy(ctx, getPopularVideoKey(), 1, strconv.FormatInt(videoID, 10)).Err(); err != nil {
		return errors.Wrapf(err, "IncrPopularVideoVisitCount failed, videoID=%d", videoID)
	}
	return nil
}
