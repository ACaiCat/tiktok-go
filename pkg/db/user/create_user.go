package userdao

import (
	"log"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (u *UserDao) CreateUser(username string, password string) (int64, error) {
	user := &model.User{
		Username: username,
		Password: password,
	}

	err := u.q.User.Create(user)
	if err != nil {
		log.Println("failed to create user: ", err)
		return 0, err
	}

	return user.ID, nil
}
