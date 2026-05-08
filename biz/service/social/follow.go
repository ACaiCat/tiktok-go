package service

import (
	"strconv"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/social"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func (s *SocialService) FollowAction(req *social.FollowReq, followerID int64) error {
	userID, err := strconv.ParseInt(req.ToUserID, 10, 64)
	if err != nil {
		return errno.ParamErr.WithError(err)
	}

	if userID == followerID {
		return errno.FollowSelfErr
	}

	if req.ActionType != social.FollowActionType_FOLLOW && req.ActionType != social.FollowActionType_UNFOLLOW {
		return errno.NotSupportActionErr
	}

	exists, err := s.userDao.IsUserExists(s.ctx, userID)
	if err != nil {
		return errors.WithMessagef(err, "service.FollowAction: check user exists failed, userID=%d", userID)
	}
	if !exists {
		return errno.UserIsNotExistErr
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		followed, err := s.followerDao.WithTx(tx).IsExistFollow(s.ctx, userID, followerID)

		if err != nil {
			return errors.WithMessagef(err, "service.FollowAction: check follow exists failed, userID=%d, followerID=%d", userID, followerID)
		}

		if req.ActionType == social.FollowActionType_FOLLOW {
			if followed {
				return errno.FollowAlreadyExistErr
			}

			if err := s.followerDao.WithTx(tx).AddFollow(s.ctx, userID, followerID); err != nil {
				return errors.WithMessagef(err, "service.FollowAction: db.AddFollow failed, userID=%d, followerID=%d", userID, followerID)
			}
		} else {
			if !followed {
				return errno.FollowNotExistErr
			}

			if err := s.followerDao.WithTx(tx).DeleteFollow(s.ctx, userID, followerID); err != nil {
				return errors.WithMessagef(err, "service.FollowAction: db.DeleteFollow failed, userID=%d, followerID=%d", userID, followerID)
			}
		}
		return nil
	})

	if err != nil {
		return errors.WithMessagef(err, "service.FollowAction: tx failed, userID=%d, followerID=%d", userID, followerID)
	}

	return nil
}

func (s *SocialService) ListFollowing(req *social.ListFollowingReq) ([]*model.SocialUser, int, error) {
	userID, err := strconv.ParseInt(req.UserID, 10, 64)
	if err != nil {
		return nil, 0, errno.ParamErr.WithError(err)
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = constants.DefaultSocialUserPageSize
	}

	pageNum := req.PageNum
	if pageNum < 0 {
		pageNum = 0
	}

	if pageSize > constants.MaxSocialUserPageSize {
		pageSize = constants.MaxSocialUserPageSize
	}

	var users []*modelDao.User
	var total int

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		exists, err := s.userDao.WithTx(tx).IsUserExists(s.ctx, userID)
		if err != nil {
			return errors.WithMessagef(err, "service.ListFollowing: check user exists failed, userID=%d", userID)
		}
		if !exists {
			return errno.UserIsNotExistErr
		}

		users, total, err = s.followerDao.WithTx(tx).GetFollowing(s.ctx, userID, int(pageSize), int(pageNum))
		if err != nil {
			return errors.WithMessagef(err, "service.ListFollowing: db.GetFollowing failed, userID=%d", userID)
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return UsersToSocialUsers(users), total, nil
}

func (s *SocialService) ListFollower(req *social.ListFollowerReq) ([]*model.SocialUser, int, error) {
	userID, err := strconv.ParseInt(req.UserID, 10, 64)
	if err != nil {
		return nil, 0, errno.ParamErr.WithError(err)
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = constants.DefaultSocialUserPageSize
	}

	pageNum := req.PageNum
	if pageNum < 0 {
		pageNum = 0
	}

	if pageSize > constants.MaxSocialUserPageSize {
		pageSize = constants.MaxSocialUserPageSize
	}

	var users []*modelDao.User
	var total int

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		exists, err := s.userDao.WithTx(tx).IsUserExists(s.ctx, userID)
		if err != nil {
			return errors.WithMessagef(err, "service.ListFollower: check user exists failed, userID=%d", userID)
		}

		if !exists {
			return errno.UserIsNotExistErr
		}

		users, total, err = s.followerDao.WithTx(tx).GetFollower(s.ctx, userID, int(pageSize), int(pageNum))
		if err != nil {
			return errors.WithMessagef(err, "service.ListFollower: db.GetFollower failed, userID=%d", userID)
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return UsersToSocialUsers(users), total, nil
}

func (s *SocialService) ListFriend(req *social.ListFriendReq, userID int64) ([]*model.SocialUser, int, error) {
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = constants.DefaultSocialUserPageSize
	}

	pageNum := req.PageNum
	if pageNum < 0 {
		pageNum = 0
	}

	if pageSize > constants.MaxSocialUserPageSize {
		pageSize = constants.MaxSocialUserPageSize
	}

	var users []*modelDao.User
	var total int

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		exists, err := s.userDao.WithTx(tx).IsUserExists(s.ctx, userID)
		if err != nil {
			return errors.WithMessagef(err, "service.ListFriend: check user exists failed, userID=%d", userID)
		}

		if !exists {
			return errno.UserIsNotExistErr
		}

		users, total, err = s.followerDao.WithTx(tx).GetFriends(s.ctx, userID, int(pageSize), int(pageNum))
		if err != nil {
			return errors.WithMessagef(err, "service.ListFriend: db.GetFriends failed, userID=%d", userID)
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return UsersToSocialUsers(users), total, nil
}
