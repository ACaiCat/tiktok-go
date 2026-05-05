package userdao

import (
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/pkg/db/query"
)

type UserDao struct {
	q *query.Query
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{q: query.Use(db)}
}

func (u *UserDao) WithTx(tx *gorm.DB) *UserDao {
	return &UserDao{q: query.Use(tx)}
}
