package userdao

import (
	"context"
	"log"
)

func (u *UserDao) UpdateUserMFA(ctx context.Context, userID int64, secrete string) error {
	_, err := u.q.User.WithContext(ctx).
		Where(u.q.User.ID.Eq(userID)).
		Update(u.q.User.TotpSecret, secrete)
	if err != nil {
		log.Println("failed to update user mfa secret: ", err)
		return err
	}

	return nil
}

func (u *UserDao) UpdateUserAvatarURL(ctx context.Context, userID int64, avatarURL string) error {
	_, err := u.q.User.WithContext(ctx).
		Where(u.q.User.ID.Eq(userID)).
		Update(u.q.User.AvatarURL, avatarURL)
	if err != nil {
		log.Println("failed to update user avatar url: ", err)
		return err
	}
	return nil
}
