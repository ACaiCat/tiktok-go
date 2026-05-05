package followerdao

import (
	"context"
	"log"
)

const mutualFollowCount = 2

func (f *FollowerDao) IsExistFollow(ctx context.Context, userID int64, followerID int64) (bool, error) {
	count, err := f.q.Follower.WithContext(ctx).
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

func (f *FollowerDao) IsExistFriend(ctx context.Context, userID int64, friendID int64) (bool, error) {
	count, err := f.q.Follower.WithContext(ctx).
		Where(
			f.q.Follower.UserID.Eq(userID), f.q.Follower.FollowerID.Eq(friendID),
		).Or(
		f.q.Follower.UserID.Eq(friendID), f.q.Follower.FollowerID.Eq(userID),
	).Count()

	if err != nil {
		log.Println("failed to check exist friend:", err)
		return false, err
	}

	return count == mutualFollowCount, nil
}
