package videoDao

import (
	"log"
	"strings"
	"time"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (v *VideoDao) SearchVideo(
	keywords []string,
	pageSize int, pageNum int,
	fromDate time.Time, toDate time.Time,
	username string,
) ([]*model.Video, error) {
	var err error

	var statement = v.q.Video.Where()
	if !fromDate.IsZero() {
		statement = statement.Where(v.q.Video.CreatedAt.Gt(fromDate))
	}

	if !toDate.IsZero() {
		statement = statement.Where(v.q.Video.CreatedAt.Lt(toDate))
	}

	if len(keywords) > 0 {
		condition := v.q.Video.Where()
		for _, keyword := range keywords {
			pattern := "%" + keyword + "%"
			condition = condition.Or(v.q.Video.Title.Like(pattern)).Or(v.q.Video.Description.Like(pattern))
		}
		statement = statement.Where(condition)
	}

	if strings.TrimSpace(username) != "" {
		statement = statement.
			Join(v.q.User, v.q.User.ID.EqCol(v.q.Video.UserID)).
			Where(v.q.User.Username.Like(username))
	}

	videos, err := statement.
		Select(v.q.Video.ALL, v.q.Like.ID.Count().As("like_count")).
		LeftJoin(v.q.Like, v.q.Like.VideoID.EqCol(v.q.Video.ID)).
		Group(v.q.Video.ID).
		Offset(pageSize * pageNum).Limit(pageSize).Find()

	if err != nil {
		log.Println("failed to search videos:", err)
		return nil, err
	}

	return videos, nil
}
