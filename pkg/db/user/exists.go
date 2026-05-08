package userdao

import (
	"context"

	"github.com/pkg/errors"
)

func (u *UserDao) IsUserExists(ctx context.Context, userID int64) (bool, error) {
	var err error

	count, err := u.q.User.WithContext(ctx).
		Select(u.q.User.ID).
		Where(u.q.User.ID.Eq(userID)).
		Limit(1).
		Count()

	if err != nil {
		return false, errors.Wrapf(err, "IsUserExists failed, userID: %d", userID)
	}

	return count > 0, nil
}
