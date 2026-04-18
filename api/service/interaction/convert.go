package service

import (
	"strconv"

	"github.com/ACaiCat/tiktok-go/api/model/model"
	commentDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func CommentDaoToDto(comment *commentDao.Comment) *model.Comment {
	return &model.Comment{
		ID:         strconv.FormatInt(comment.ID, 10),
		UserID:     strconv.FormatInt(comment.UserID, 10),
		VideoID:    strconv.FormatInt(comment.VideoID, 10),
		Content:    comment.Content,
		CreatedAt:  strconv.FormatInt(comment.CreatedAt.UnixMilli(), 10),
		LikeCount:  comment.LikeCount,
		ChildCount: comment.ChildCount,
	}
}

func CommentsDaoToDto(commentsDao []*commentDao.Comment) []*model.Comment {
	comments := make([]*model.Comment, 0, len(commentsDao))

	for _, comment := range commentsDao {
		comments = append(comments, CommentDaoToDto(comment))
	}

	return comments
}
