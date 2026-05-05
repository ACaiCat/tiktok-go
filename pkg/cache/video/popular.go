package videocache

import (
	"context"
	"encoding/json"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func GetPopularVideoKey() string {
	return "popular_videos"
}

func (p *VideoCache) SetPopularVideos(ctx context.Context, videos []*model.Video) error {
	data, err := json.Marshal(videos)
	if err != nil {
		return err
	}

	return p.c.Set(ctx, GetPopularVideoKey(), data, constants.PopularVideoCacheExpiration).Err()
}

func (p *VideoCache) GetPopularVideos(ctx context.Context) ([]*model.Video, error) {
	var videos []*model.Video
	data, err := p.c.Get(ctx, GetPopularVideoKey()).Bytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &videos)
	if err != nil {
		return nil, err
	}
	return videos, nil
}
