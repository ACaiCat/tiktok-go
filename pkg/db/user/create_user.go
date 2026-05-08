package userdao

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (u *UserDao) CreateUser(ctx context.Context, username string, password string) (int64, error) {
	user := &model.User{
		Username: username,
		Password: password,
	}

	err := u.q.User.WithContext(ctx).Create(user)
	if err != nil {
		return 0, errors.Wrapf(err, "CreateUser failed, user: %s", username)
	}

	return user.ID, nil
}
