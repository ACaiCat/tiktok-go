package service

import (
	"strconv"

	"github.com/ACaiCat/tiktok-go/api/model/model"
	"github.com/ACaiCat/tiktok-go/api/model/social"
	"github.com/ACaiCat/tiktok-go/pkg/constants"
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

	exists, err := s.userDao.IsUserExists(userID)
	if err != nil {
		return errno.ServiceErr
	}
	if !exists {
		return errno.UserIsNotExistErr
	}

	followed, err := s.followerDao.IsExistFollow(userID, followerID)

	if err != nil {
		return errno.ServiceErr.WithError(err)
	}

	switch req.ActionType {
	case model.FollowActionType_FOLLOW:
		if followed {
			return errno.FollowAlreadyExistErr
		}

		err := s.followerDao.AddFollow(userID, followerID)
		if err != nil {
			return errno.ServiceErr
		}
		return nil
	case model.FollowActionType_UNFOLLOW:
		if !followed {
			return errno.FollowNotExistErr
		}

		err := s.followerDao.DeleteFollow(userID, followerID)
		if err != nil {
			return errno.ServiceErr
		}
		return nil
	}

	return errno.NotSupportActionErr
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

	exists, err := s.userDao.IsUserExists(userID)
	if err != nil {
		return nil, 0, errno.ServiceErr
	}
	if !exists {
		return nil, 0, errno.UserIsNotExistErr
	}

	users, total, err := s.followerDao.GetFollowing(userID, int(pageSize), int(pageNum))
	if err != nil {
		return nil, 0, errno.ServiceErr
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

	exists, err := s.userDao.IsUserExists(userID)
	if err != nil {
		return nil, 0, errno.ServiceErr
	}

	if !exists {
		return nil, 0, errno.UserIsNotExistErr
	}

	users, total, err := s.followerDao.GetFollower(userID, int(pageSize), int(pageNum))
	if err != nil {
		return nil, 0, errno.ServiceErr
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

	exists, err := s.userDao.IsUserExists(userID)
	if err != nil {
		return nil, 0, errno.ServiceErr
	}

	if !exists {
		return nil, 0, errno.UserIsNotExistErr
	}

	users, total, err := s.followerDao.GetFriends(userID, int(pageSize), int(pageNum))
	if err != nil {
		return nil, 0, errno.ServiceErr
	}

	return UsersToSocialUsers(users), total, nil
}
