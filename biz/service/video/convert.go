package service

import (
	"strconv"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	model2 "github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func VideoDaoToDto(video *model2.Video, likeCount int64, commentCount int64) *model.Video {
	return &model.Video{
		ID:           strconv.FormatInt(video.ID, 10),
		UserID:       strconv.FormatInt(video.UserID, 10),
		VideoURL:     video.VideoURL,
		CoverURL:     video.CoverURL,
		Title:        video.Title,
		Description:  video.Description,
		VisitCount:   video.VisitCount,
		LikeCount:    likeCount,
		CommentCount: commentCount,
		CreatedAt:    strconv.FormatInt(video.CreatedAt.UnixMilli(), 10),
	}
}

func VideosDaoToDto(videos []*model2.Video, likeCounts map[int64]int64, commentCounts map[int64]int64) []*model.Video {
	videoDtos := make([]*model.Video, len(videos))
	for i, video := range videos {
		videoDtos[i] = VideoDaoToDto(video, likeCounts[video.ID], commentCounts[video.ID])
	}

	return videoDtos
}
