package userdao

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (u *UserDao) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := u.q.User.WithContext(ctx).
		Where(u.q.User.ID.Eq(id)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		log.Println("failed to get user by id:", err)
		return nil, err
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

		log.Println("failed to get user by username:", err)
		return nil, err
	}
	return user, nil
}
