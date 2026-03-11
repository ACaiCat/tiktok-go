package videoDao

import (
	"log"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
	"github.com/ACaiCat/tiktok-go/pkg/db/query"
)

func (v *VideoDao) PublishVideo(
	userID int64,
	title string,
	description string,
	uploadFn func(videoID int64) error,
	videoURLFn func(int64) string,
	coverURLFn func(int64) string,
) error {
	return v.q.Transaction(func(tx *query.Query) error {
		video := model.Video{
			UserID:      userID,
			Title:       title,
			Description: description,
			VisitCount:  0,
		}
		if err := tx.Video.Create(&video); err != nil {
			log.Printf("failed to publish video in tx: %v\n", err)
			return err
		}
		if err := uploadFn(video.ID); err != nil {
			log.Printf("failed to upload files in tx: %v\n", err)
			return err
		}
		_, err := tx.Video.Where(tx.Video.ID.Eq(video.ID)).
			Updates(map[string]interface{}{
				"video_url": videoURLFn(video.ID),
				"cover_url": coverURLFn(video.ID),
			})
		if err != nil {
			log.Printf("failed to update video url and cover url in tx: %v\n", err)
			return err
		}
		return nil
	})
}
