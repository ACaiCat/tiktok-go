package videocache

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func (p *VideoCache) GetPopularVideos(ctx context.Context, pageSize int, pageNum int) ([]int64, error) {
	start := int64(pageSize * pageNum)
	stop := start + int64(pageSize) - 1

	videoIDsStr, err := p.c.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   getPopularVideoKey(),
		Start: start,
		Stop:  stop,
		Rev:   true,
	}).Result()
	if err != nil {
		return nil, errors.Wrap(err, "GetPopularVideos failed")
	}
	if len(videoIDsStr) == 0 {
		return nil, redis.Nil
	}

	videoIDs := make([]int64, 0, len(videoIDsStr))
	for _, videoIDStr := range videoIDsStr {
		videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "GetPopularVideos parse videoID failed, videoID=%q", videoIDStr)
		}
		videoIDs = append(videoIDs, videoID)
	}

	return videoIDs, nil
}
