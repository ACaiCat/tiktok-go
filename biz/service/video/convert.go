package service

import (
	"strconv"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func VideoDaoToDto(video *modelDao.Video) *model.Video {
	return &model.Video{
		ID:           strconv.FormatInt(video.ID, 10),
		UserID:       strconv.FormatInt(video.UserID, 10),
		VideoURL:     video.VideoURL,
		CoverURL:     video.CoverURL,
		Title:        video.Title,
		Description:  video.Description,
		VisitCount:   video.VisitCount,
		LikeCount:    video.LikeCount,
		CommentCount: video.CommentCount,
		CreatedAt:    strconv.FormatInt(video.CreatedAt.UnixMilli(), 10),
	}
}

func VideosDaoToDto(videos []*modelDao.Video) []*model.Video {
	videoDtos := make([]*model.Video, len(videos))
	for i, video := range videos {
		videoDtos[i] = VideoDaoToDto(video)
	}

	return videoDtos
}
