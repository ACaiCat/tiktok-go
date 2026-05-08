package videodao

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (v *VideoDao) SearchVideo(
	ctx context.Context,
	keywords []string,
	pageSize int, pageNum int,
	fromDate time.Time, toDate time.Time,
	username string,
) ([]*model.Video, error) {
	var err error

	var statement = v.q.Video.WithContext(ctx).Where()
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
		Offset(pageSize * pageNum).Limit(pageSize).Find()

	if err != nil {
		return nil, errors.Wrapf(err, "SearchVideo failed, keywords: %s, fromDate: %s, username: %s",
			strings.Join(keywords, ","), toDate.Format(time.RFC3339), username)
	}

	return videos, nil
}
