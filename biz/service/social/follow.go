package service

import (
	"context"
	"slices"
	"strconv"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/biz/model/model"
	"github.com/ACaiCat/tiktok-go/biz/model/social"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
	"github.com/ACaiCat/tiktok-go/pkg/utils"
)

func (s *SocialService) FollowAction(req *social.FollowReq, followerID int64) error {
	followingID, err := strconv.ParseInt(req.ToUserID, 10, 64)
	if err != nil {
		return errno.ParamErr.WithError(err)
	}

	if followingID == followerID {
		return errno.FollowSelfErr
	}

	if req.ActionType != social.FollowActionType_FOLLOW && req.ActionType != social.FollowActionType_UNFOLLOW {
		return errno.NotSupportActionErr
	}

	exists, err := s.userDao.IsUserExists(s.ctx, followingID)
	if err != nil {
		return errors.WithMessagef(err, "service.FollowAction: check user exists failed, userID=%d", followingID)
	}
	if !exists {
		return errno.UserIsNotExistErr
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		followed, err := s.userCache.IsFollowed(s.ctx, followerID, followingID)
		if err != nil {
			if !errors.Is(err, redis.Nil) {
				hlog.CtxErrorf(s.ctx, "service.FollowAction: IsFollowed failed, userID=%d, followingID=%d, err=%v", followingID, followerID, err)
			}

			followingIDs, err := s.followerDao.WithTx(tx).GetFollowingIDs(context.Background(), followerID)
			if err != nil {
				return errors.WithMessagef(err, "service.FollowAction: GetFollowingIDs failed, userID=%d", followingID)
			}

			followed = slices.Contains(followingIDs, followingID)

			go func() {
				err := s.userCache.SetFollowings(context.Background(), followerID, followingIDs)
				if err != nil {
					hlog.Errorf("service.FollowAction: SetFollowings failed, userID=%d, err=%v", followingID, err)
				}
			}()
		}

		if req.ActionType == social.FollowActionType_FOLLOW {
			if followed {
				return errno.FollowAlreadyExistErr
			}

			if err := s.followerDao.WithTx(tx).AddFollow(s.ctx, followingID, followerID); err != nil {
				return errors.WithMessagef(err, "service.FollowAction: db.AddFollow failed, userID=%d, followerID=%d", followingID, followerID)
			}

			go func() {
				err := s.userCache.SetFollow(context.Background(), followerID, followingID)
				if err != nil {
					hlog.Errorf("service.FollowAction: SetFollow failed, userID=%d, err=%v", followingID, err)
				}
			}()
		} else {
			if !followed {
				return errno.FollowNotExistErr
			}

			if err := s.followerDao.WithTx(tx).DeleteFollow(s.ctx, followingID, followerID); err != nil {
				return errors.WithMessagef(err, "service.FollowAction: db.DeleteFollow failed, userID=%d, followerID=%d", followingID, followerID)
			}

			go func() {
				err := s.userCache.SetUnfollow(context.Background(), followerID, followingID)
				if err != nil {
					hlog.Errorf("service.FollowAction: SetUnfollow failed, userID=%d, err=%v", followingID, err)
				}
			}()

		}
		return nil
	})

	if err != nil {
		return errors.WithMessagef(err, "service.FollowAction: tx failed, userID=%d, followerID=%d", followingID, followerID)
	}

	return nil
}

func (s *SocialService) ListFollowing(req *social.ListFollowingReq) ([]*model.SocialUser, int, error) {
	userID, err := strconv.ParseInt(req.UserID, 10, 64)
	if err != nil {
		return nil, 0, errno.ParamErr.WithError(err)
	}

	pageSize, pageNum := utils.NormalizePage(req.PageSize, req.PageNum, constants.DefaultSocialUserPageSize, constants.MaxSocialUserPageSize)

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

		users, total, err = s.followerDao.WithTx(tx).GetFollowing(s.ctx, userID, pageSize, pageNum)
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

	pageSize, pageNum := utils.NormalizePage(req.PageSize, req.PageNum, constants.DefaultSocialUserPageSize, constants.MaxSocialUserPageSize)

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

		users, total, err = s.followerDao.WithTx(tx).GetFollower(s.ctx, userID, pageSize, pageNum)
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
	pageSize, pageNum := utils.NormalizePage(req.PageSize, req.PageNum, constants.DefaultSocialUserPageSize, constants.MaxSocialUserPageSize)

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

		users, total, err = s.followerDao.WithTx(tx).GetFriends(s.ctx, userID, pageSize, pageNum)
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
