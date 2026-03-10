package service

import (
	"strconv"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func VideoDaoToDto(video *modelDao.Video, likeCount int64, commentCount int64) *model.Video {
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

func VideosDaoToDto(videos []*modelDao.Video, likeCounts map[int64]int64, commentCounts map[int64]int64) []*model.Video {
	videoDtos := make([]*model.Video, len(videos))
	for i, video := range videos {
		videoDtos[i] = VideoDaoToDto(video, likeCounts[video.ID], commentCounts[video.ID])
	}

	return videoDtos
}

func (s *VideoService) GetLikeAndCommentCount(videosDao []*modelDao.Video) ([]*model.Video, error) {
	videoIDs := make([]int64, len(videosDao))
	for i, v := range videosDao {
		videoIDs[i] = v.ID
	}

	idWithCommentCount, err := s.commentDao.GetCommentCounts(videoIDs)
	if err != nil {
		return nil, errno.ServiceErr
	}

	idWithLikeCount, err := s.likeDao.GetLikeCounts(videoIDs)
	if err != nil {
		return nil, errno.ServiceErr
	}

	videos := VideosDaoToDto(videosDao, idWithLikeCount, idWithCommentCount)

	return videos, nil
}
