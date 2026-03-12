package commentDao

import (
	"log"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (c *CommentDao) AddVideoComment(userID int64, videoID int64, content string) error {
	var err error

	comment := model.Comment{
		UserID:  userID,
		VideoID: videoID,
		Content: content,
	}

	err = c.q.Comment.Create(&comment)
	if err != nil {
		log.Printf("failed to add video comment: %v", err)
		return err
	}

	return nil
}
