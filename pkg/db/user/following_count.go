package userdao

import "log"

func (u *UserDao) IncrFollowingCount(userID int64) error {
	_, err := u.q.User.
		Where(u.q.User.ID.Eq(userID)).
		UpdateColumn(u.q.User.FollowingCount, u.q.User.FollowingCount.Add(1))

	if err != nil {
		log.Printf("failed to increase following count: %v", err)
		return err
	}

	return nil
}

func (u *UserDao) DecrFollowingCount(userID int64) error {
	_, err := u.q.User.
		Where(u.q.User.ID.Eq(userID)).
		UpdateColumn(u.q.User.FollowingCount, u.q.User.FollowingCount.Add(-1))

	if err != nil {
		log.Printf("failed to decrease following count: %v", err)
		return err
	}

	return nil
}
