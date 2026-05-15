package videocache

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

type userVideoListCache struct {
	VideoIDs []int64 `json:"video_ids"`
	Total    int64   `json:"total"`
}

func (p *VideoCache) GetUserVideoList(ctx context.Context, userID int64, pageSize int, pageNum int) ([]int64, int64, error) {
	version, err := p.getUserVideoListVersion(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	data, err := p.c.Get(ctx, getUserVideoListKey(userID, version, pageSize, pageNum)).Bytes()
	if err != nil {
		return nil, 0, errors.Wrapf(err, "GetUserVideoList failed, userID=%d", userID)
	}

	var list userVideoListCache
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, 0, errors.Wrapf(err, "GetUserVideoList json failed, userID=%d", userID)
	}

	return list.VideoIDs, list.Total, nil
}

func (p *VideoCache) getUserVideoListVersion(ctx context.Context, userID int64) (int64, error) {
	versionStr, err := p.c.Get(ctx, getUserVideoListVersionKey(userID)).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "GetUserVideoListVersion failed, userID=%d", userID)
	}

	version, err := strconv.ParseInt(versionStr, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "getUserVideoListVersion parse failed, userID=%d", userID)
	}
	return version, nil
}
