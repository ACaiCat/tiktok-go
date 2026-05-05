package userdao

import (
	"context"
	"log"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (u *UserDao) CreateUser(ctx context.Context, username string, password string) (int64, error) {
	user := &model.User{
		Username: username,
		Password: password,
	}

	err := u.q.User.WithContext(ctx).Create(user)
	if err != nil {
		log.Println("failed to create user: ", err)
		return 0, err
	}

	return user.ID, nil
}
