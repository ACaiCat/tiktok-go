package followerdao

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
	"github.com/ACaiCat/tiktok-go/pkg/db/query"
)

const mutualFollowCount = 2

func (f *FollowerDao) GetFollower(ctx context.Context, userID int64, pageSize int, pageNum int) ([]*model.User, int, error) {
	var err error

	var userIDs []int64
	var users []*model.User

	err = f.q.Transaction(func(tx *query.Query) error {
		err = tx.Follower.WithContext(ctx).
			Select(f.q.Follower.FollowerID).
			Where(f.q.Follower.UserID.Eq(userID)).
			Scan(&userIDs)

		if err != nil {
			return errors.Wrapf(err, "GetFollower failed, userID: %d", userID)
		}

		users, err = tx.User.WithContext(ctx).
			Where(f.q.User.ID.In(userIDs...)).
			Offset(pageSize * pageNum).
			Limit(pageSize).
			Find()

		if err != nil {
			return errors.Wrapf(err, "GetFollower failed, userID: %d", userID)
		}
		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return users, len(userIDs), nil
}

func (f *FollowerDao) GetFollowing(ctx context.Context, userID int64, pageSize int, pageNum int) ([]*model.User, int, error) {
	var err error

	var userIDs []int64
	var users []*model.User

	err = f.q.Transaction(func(tx *query.Query) error {
		err = tx.Follower.WithContext(ctx).
			Select(f.q.Follower.UserID).
			Where(f.q.Follower.FollowerID.Eq(userID)).
			Scan(&userIDs)

		if err != nil {
			return errors.Wrapf(err, "GetFollowing failed, userID: %d", userID)
		}

		users, err = tx.User.WithContext(ctx).
			Where(f.q.User.ID.In(userIDs...)).
			Offset(pageSize * pageNum).
			Limit(pageSize).
			Find()

		if err != nil {
			return errors.Wrapf(err, "GetFollowing failed, userID: %d", userID)
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return users, len(userIDs), nil
}

func (f *FollowerDao) GetFollowingIDs(ctx context.Context, userID int64) ([]int64, error) {
	var err error

	var userIDs []int64

	err = f.q.Follower.WithContext(ctx).
		Select(f.q.Follower.UserID).
		Where(f.q.Follower.FollowerID.Eq(userID)).
		Scan(&userIDs)

	if err != nil {
		return nil, errors.Wrapf(err, "GetFollowingIDs failed, userID: %d", userID)
	}

	return userIDs, nil
}

func (f *FollowerDao) GetFriends(ctx context.Context, userID int64, pageSize int, pageNum int) ([]*model.User, int, error) {
	var err error

	var followerIDs []int64
	var friendIDs []int64
	var users []*model.User

	err = f.q.Transaction(func(tx *query.Query) error {
		err = tx.Follower.WithContext(ctx).
			Select(f.q.Follower.UserID).
			Where(f.q.Follower.FollowerID.Eq(userID)).
			Scan(&followerIDs)

		if err != nil {
			return errors.Wrapf(err, "GetFriend failed, userID: %d", userID)
		}

		err = tx.Follower.WithContext(ctx).
			Select(f.q.Follower.FollowerID).
			Where(f.q.Follower.UserID.Eq(userID), f.q.Follower.FollowerID.In(followerIDs...)).
			Scan(&friendIDs)

		if err != nil {
			return errors.Wrapf(err, "GetFriend failed, userID: %d", userID)
		}

		users, err = tx.User.WithContext(ctx).
			Where(f.q.User.ID.In(friendIDs...)).
			Offset(pageSize * pageNum).
			Limit(pageSize).
			Find()

		if err != nil {
			return errors.Wrapf(err, "GetFriend failed, userID: %d", userID)
		}
		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return users, len(followerIDs), nil
}

func (f *FollowerDao) IsExistFollow(ctx context.Context, userID int64, followerID int64) (bool, error) {
	count, err := f.q.Follower.WithContext(ctx).
		Select(f.q.Follower.ID).
		Where(f.q.Follower.UserID.Eq(userID), f.q.Follower.FollowerID.Eq(followerID)).
		Limit(1).
		Count()

	if err != nil {
		return false, errors.Wrapf(err, "IsExistFollow failed, count: %d, userID: %d", count, userID)
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
		return false, errors.Wrapf(err, "IsExistFriend failed, count: %d, userID: %d", count, userID)
	}

	return count == mutualFollowCount, nil
}
