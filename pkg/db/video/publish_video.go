package videodao

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
	"github.com/ACaiCat/tiktok-go/pkg/db/query"
)

func (v *VideoDao) PublishVideo(
	ctx context.Context,
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
		if err := tx.Video.WithContext(ctx).Create(&video); err != nil {
			return errors.Wrapf(err, "PublishVideo failed, videoID: %d", video.ID)
		}
		if err := uploadFn(video.ID); err != nil {
			return errors.Wrapf(err, "PublishVideo failed, videoID: %d", video.ID)
		}
		_, err := tx.Video.WithContext(ctx).Where(tx.Video.ID.Eq(video.ID)).
			Updates(map[string]any{
				"video_url": videoURLFn(video.ID),
				"cover_url": coverURLFn(video.ID),
			})
		if err != nil {
			return errors.Wrapf(err, "PublishVideo failed, videoID: %d", video.ID)
		}
		return nil
	})
}
