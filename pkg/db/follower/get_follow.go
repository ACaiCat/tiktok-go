package followerDao

import (
	"log"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func (f *FollowerDao) GetFollower(userID int64, pageSize int, pageNum int) ([]*model.User, int, error) {
	var err error

	var userIDs []int64

	err = f.q.Follower.
		Select(f.q.Follower.FollowerID).
		Where(f.q.Follower.UserID.Eq(userID)).
		Scan(&userIDs)

	if err != nil {
		log.Println("failed to get follower IDs for userID", userID, ":", err)
		return nil, 0, err
	}

	users, err := f.q.User.
		Where(f.q.User.ID.In(userIDs...)).
		Offset(pageSize * pageNum).
		Limit(pageSize).
		Find()

	if err != nil {
		log.Println("failed to get followers for userID", userID, ":", err)
		return nil, 0, err
	}

	return users, len(userIDs), nil
}

func (f *FollowerDao) GetFollowing(userID int64, pageSize int, pageNum int) ([]*model.User, int, error) {
	var err error

	var userIDs []int64

	err = f.q.Follower.
		Select(f.q.Follower.UserID).
		Where(f.q.Follower.FollowerID.Eq(userID)).
		Scan(&userIDs)

	if err != nil {
		log.Println("failed to get following userIDs for userID", userID, ":", err)
		return nil, 0, err
	}

	users, err := f.q.User.
		Where(f.q.User.ID.In(userIDs...)).
		Offset(pageSize * pageNum).
		Limit(pageSize).
		Find()

	if err != nil {
		log.Println("failed to get followings for userID", userID, ":", err)
		return nil, 0, err
	}

	return users, len(userIDs), nil
}

func (f *FollowerDao) GetFriends(userID int64, pageSize int, pageNum int) ([]*model.User, int, error) {
	var err error

	var userIDs []int64

	err = f.q.Follower.
		Select(f.q.Follower.UserID).
		Where(f.q.Follower.FollowerID.Eq(userID)).
		Scan(&userIDs)

	if err != nil {
		log.Println("failed to get following userIDs for userID", userID, ":", err)
		return nil, 0, err
	}

	err = f.q.Follower.
		Select(f.q.Follower.FollowerID).
		Where(f.q.Follower.UserID.Eq(userID), f.q.Follower.FollowerID.In(userIDs...)).
		Scan(&userIDs)

	if err != nil {
		log.Println("failed to get friend userIDs for userID", userID, ":", err)
		return nil, 0, err
	}

	users, err := f.q.User.
		Where(f.q.User.ID.In(userIDs...)).
		Offset(pageSize * pageNum).
		Limit(pageSize).
		Find()

	if err != nil {
		log.Println("failed to get friends for userID", userID, ":", err)
		return nil, 0, err
	}

	return users, len(userIDs), nil
}
