package followerDao

import "log"

func (f *FollowerDao) IsExistFollow(userID int64, followerID int64) (bool, error) {
	count, err := f.q.Follower.
		Select(f.q.Follower.ID).
		Where(f.q.Follower.UserID.Eq(userID), f.q.Follower.FollowerID.Eq(followerID)).
		Limit(1).
		Count()

	if err != nil {
		log.Println("failed to check exist follow:", err)
		return false, err
	}
	return count > 0, nil
}
