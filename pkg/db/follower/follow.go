package followerdao

import (
	"context"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
	"github.com/pkg/errors"
)

func (f *FollowerDao) AddFollow(ctx context.Context, userID int64, followerID int64) error {
	var err error

	follower := model.Follower{
		UserID:     userID,
		FollowerID: followerID,
	}

	err = f.q.Follower.WithContext(ctx).Create(&follower)
	if err != nil {
		return errors.Wrapf(err, "AddFollow failed, create follower failed")
	}

	return nil
}

func (f *FollowerDao) DeleteFollow(ctx context.Context, userID int64, followerID int64) error {
	var err error

	_, err = f.q.Follower.WithContext(ctx).
		Where(f.q.Follower.UserID.Eq(userID), f.q.Follower.FollowerID.Eq(followerID)).
		Delete()

	if err != nil {
		return errors.Wrapf(err, "DeleteFollow failed, delete follower failed")
	}

	return nil
}
