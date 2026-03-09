package userDao

import (
	"github.com/ACaiCat/tiktok-go/pkg/db/query"
)

type UserDao struct {
	q *query.Query
}

func NewUserDao() *UserDao {
	return &UserDao{q: query.Q}
}
