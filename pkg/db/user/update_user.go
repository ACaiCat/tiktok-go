package userdao

import (
	"context"

	"github.com/pkg/errors"
)

func (u *UserDao) UpdateUserMFA(ctx context.Context, userID int64, secrete string) error {
	_, err := u.q.User.WithContext(ctx).
		Where(u.q.User.ID.Eq(userID)).
		Update(u.q.User.TotpSecret, secrete)
	if err != nil {
		return errors.Wrapf(err, "UpdateUserMFA failed, userID: %d", userID)
	}

	return nil
}

func (u *UserDao) UpdateUserAvatarURL(ctx context.Context, userID int64, avatarURL string) error {
	_, err := u.q.User.WithContext(ctx).
		Where(u.q.User.ID.Eq(userID)).
		Update(u.q.User.AvatarURL, avatarURL)
	if err != nil {
		return errors.Wrapf(err, "UpdateUserAvatarURL failed, userID: %d", userID)
	}
	return nil
}

func (u *UserDao) UpdateUserJwch(ctx context.Context, userID int64, jwchID string, jwchPassword string) error {
	_, err := u.q.User.WithContext(ctx).
		Where(u.q.User.ID.Eq(userID)).
		Updates(map[string]any{
			"jwch_password": jwchPassword,
			"jwch_id":       jwchID,
		})
	if err != nil {
		return errors.Wrapf(err, "UpdateUserJwch failed, userID: %d", userID)
	}

	return nil
}
