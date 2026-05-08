package userdao

import (
	"context"

	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"

	"github.com/pkg/errors"
)

func (u *UserDao) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := u.q.User.WithContext(ctx).
		Where(u.q.User.ID.Eq(id)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "GetUserByID failed, userID: %d", id)
	}
	return user, nil
}

func (u *UserDao) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := u.q.User.WithContext(ctx).
		Where(u.q.User.Username.Eq(username)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, errors.Wrapf(err, "GetUserByUsername failed, user: %s", username)
	}
	return user, nil
}
