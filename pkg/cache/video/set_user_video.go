package videocache

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (p *VideoCache) SetUserVideoList(ctx context.Context, userID int64, pageSize int, pageNum int, total int64, videos []*model.Video) error {
	version, err := p.getUserVideoListVersion(ctx, userID)
	if err != nil {
		return err
	}

	videoIDs := make([]int64, 0, len(videos))
	for _, video := range videos {
		videoIDs = append(videoIDs, video.ID)
	}

	data, err := json.Marshal(userVideoListCache{VideoIDs: videoIDs, Total: total})
	if err != nil {
		return errors.Wrapf(err, "SetUserVideoList json failed, userID=%d", userID)
	}

	key := getUserVideoListKey(userID, version, pageSize, pageNum)
	if err := p.c.Set(ctx, key, data, constants.UserVideoCacheExpiration).Err(); err != nil {
		return errors.Wrapf(err, "SetUserVideoList failed, userID=%d", userID)
	}

	return nil
}

func (p *VideoCache) ClearUserVideoList(ctx context.Context, userID int64) error {
	if err := p.c.Incr(ctx, getUserVideoListVersionKey(userID)).Err(); err != nil {
		return errors.Wrapf(err, "ClearUserVideoList failed, userID=%d", userID)
	}
	return nil
}
