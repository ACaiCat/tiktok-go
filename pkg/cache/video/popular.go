package videocache

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func getPopularVideoKey() string {
	return "popular_videos"
}

func (p *VideoCache) SetPopularVideos(ctx context.Context, videos []*model.Video) error {
	data, err := json.Marshal(videos)
	if err != nil {
		return errors.Wrapf(err, "SetPopularVideosJson failed, %v", videos)
	}

	if err := p.c.Set(ctx, getPopularVideoKey(), data, constants.PopularVideoCacheExpiration).Err(); err != nil {
		return errors.Wrapf(err, "SetPopularVideosCache failed, %v", videos)
	}

	return nil
}

func (p *VideoCache) GetPopularVideos(ctx context.Context) ([]*model.Video, error) {
	var videos []*model.Video
	data, err := p.c.Get(ctx, getPopularVideoKey()).Bytes()
	if err != nil {
		return nil, errors.Wrapf(err, "GetPopularVideos failed, %v", videos)
	}
	err = json.Unmarshal(data, &videos)
	if err != nil {
		return nil, errors.Wrapf(err, "GetPopularVideosJson failed, %v", data)
	}
	return videos, nil
}
