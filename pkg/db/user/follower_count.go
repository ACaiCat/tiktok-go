package userdao

import "log"

func (u *UserDao) IncrFollowerCount(userID int64) error {
	_, err := u.q.User.
		Where(u.q.User.ID.Eq(userID)).
		UpdateColumn(u.q.User.FollowerCount, u.q.User.FollowerCount.Add(1))

	if err != nil {
		log.Printf("failed to increase follower count: %v", err)
		return err
	}

	return nil
}

func (u *UserDao) DecrFollowerCount(userID int64) error {
	_, err := u.q.User.
		Where(u.q.User.ID.Eq(userID)).
		UpdateColumn(u.q.User.FollowerCount, u.q.User.FollowerCount.Add(-1))

	if err != nil {
		log.Printf("failed to decrease follower count: %v", err)
		return err
	}

	return nil
}
