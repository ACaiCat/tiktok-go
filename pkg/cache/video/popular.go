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

func (p *VideoCache) SetPopularVideos(videos []*model.Video) error {
	data, err := json.Marshal(videos)
	if err != nil {
		return err
	}

	return p.c.Set(context.Background(), GetPopularVideoKey(), data, constants.PopularVideoCacheExpiration).Err()
}

func (p *VideoCache) GetPopularVideos() ([]*model.Video, error) {
	var videos []*model.Video
	data, err := p.c.Get(context.Background(), GetPopularVideoKey()).Bytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &videos)
	if err != nil {
		return nil, err
	}
	return videos, nil
}
