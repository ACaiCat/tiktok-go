package userdao

import "log"

func (u *UserDao) IsUserExists(userID int64) (bool, error) {
	var err error

	count, err := u.q.User.
		Select(u.q.User.ID).
		Where(u.q.User.ID.Eq(userID)).
		Limit(1).
		Count()

	if err != nil {
		log.Println("failed to check if user exists for userID", userID, ":", err)
		return false, err
	}

	return count > 0, nil
}
