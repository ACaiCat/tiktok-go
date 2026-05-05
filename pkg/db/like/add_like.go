package likedao

import (
	"context"
	"log"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (l *LikeDao) AddVideoLike(ctx context.Context, userID, videoID int64) error {
	var err error

	like := model.Like{
		UserID:  userID,
		VideoID: new(videoID),
	}

	err = l.q.Like.WithContext(ctx).Create(&like)
	if err != nil {
		log.Printf("failed to add like: %v", err)
		return err
	}

	return nil
}

func (l *LikeDao) AddCommentLike(ctx context.Context, userID, commentID int64) error {
	var err error

	like := model.Like{
		UserID:    userID,
		CommentID: new(commentID),
	}

	err = l.q.Like.WithContext(ctx).Create(&like)
	if err != nil {
		log.Printf("failed to add like: %v", err)
		return err
	}

	return nil
}
