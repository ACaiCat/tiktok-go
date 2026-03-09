package userDao

import (
	"log"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (u *UserDao) GetByID(id int64) (*model.User, error) {
	user, err := u.q.User.Where(u.q.User.ID.Eq(id)).First()
	if err != nil {
		log.Println("failed to get user by id:", err)
		return nil, err
	}
	return user, nil
}

func (u *UserDao) GetByUsername(username string) (*model.User, error) {
	user, err := u.q.User.Where(u.q.User.Username.Eq(username)).First()
	if err != nil {
		log.Println("failed to get user by username:", err)
		return nil, err
	}
	return user, nil
}
