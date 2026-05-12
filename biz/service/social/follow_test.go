package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/ACaiCat/tiktok-go/biz/model/social"
	followerDao "github.com/ACaiCat/tiktok-go/pkg/db/follower"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	userDao "github.com/ACaiCat/tiktok-go/pkg/db/user"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func TestSocialService_FollowAction(t *testing.T) {
	type testCase struct {
		req            *social.FollowReq
		mockUserExists bool
		mockFollowed   bool
		expectError    string
	}

	testCases := map[string]testCase{
		"invalid user id": {
			req:         &social.FollowReq{ToUserID: "bad", ActionType: social.FollowActionType_FOLLOW},
			expectError: "invalid syntax",
		},
		"follow self": {
			req:         &social.FollowReq{ToUserID: "1", ActionType: social.FollowActionType_FOLLOW},
			expectError: errno.FollowSelfErr.ErrMsg,
		},
		"user not exists": {
			req:            &social.FollowReq{ToUserID: "2", ActionType: social.FollowActionType_FOLLOW},
			mockUserExists: false,
			expectError:    errno.UserIsNotExistErr.ErrMsg,
		},
		"follow success": {
			req:            &social.FollowReq{ToUserID: "2", ActionType: social.FollowActionType_FOLLOW},
			mockUserExists: true,
		},
		"unfollow missing relation": {
			req:            &social.FollowReq{ToUserID: "2", ActionType: social.FollowActionType_UNFOLLOW},
			mockUserExists: true,
			expectError:    errno.FollowNotExistErr.ErrMsg,
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*gorm.DB).Transaction).To(func(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
				return fc(&gorm.DB{})
			}).Build()
			mockey.Mock((*userDao.UserDao).WithTx).To(func(tx *gorm.DB) *userDao.UserDao {
				return &userDao.UserDao{}
			}).Build()
			mockey.Mock((*followerDao.FollowerDao).WithTx).To(func(tx *gorm.DB) *followerDao.FollowerDao {
				return &followerDao.FollowerDao{}
			}).Build()
			mockey.Mock((*userDao.UserDao).IsUserExists).To(func(ctx context.Context, userID int64) (bool, error) {
				return tc.mockUserExists, nil
			}).Build()
			mockey.Mock((*followerDao.FollowerDao).IsExistFollow).To(func(ctx context.Context, userID int64, followerID int64) (bool, error) {
				return tc.mockFollowed, nil
			}).Build()
			mockey.Mock((*followerDao.FollowerDao).AddFollow).Return(nil).Build()
			mockey.Mock((*followerDao.FollowerDao).DeleteFollow).Return(nil).Build()

			mockey.Mock(NewSocialService).To(func(_ context.Context) *SocialService {
				return &SocialService{}
			}).Build()

			err := NewSocialService(context.Background()).FollowAction(tc.req, 1)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestSocialService_ListFollowing(t *testing.T) {
	type testCase struct {
		req            *social.ListFollowingReq
		mockUserExists bool
		mockUsers      []*modelDao.User
		mockTotal      int
		expectError    string
	}

	testCases := map[string]testCase{
		"invalid user id": {req: &social.ListFollowingReq{UserID: "bad"}, expectError: "invalid syntax"},
		"user not exists": {req: &social.ListFollowingReq{UserID: "1"}, expectError: errno.UserIsNotExistErr.ErrMsg},
		"success":         {req: &social.ListFollowingReq{UserID: "1"}, mockUserExists: true, mockUsers: []*modelDao.User{}, mockTotal: 0},
	}

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*gorm.DB).Transaction).To(func(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
				return fc(&gorm.DB{})
			}).Build()
			mockey.Mock((*userDao.UserDao).WithTx).To(func(tx *gorm.DB) *userDao.UserDao { return &userDao.UserDao{} }).Build()
			mockey.Mock((*followerDao.FollowerDao).WithTx).To(func(tx *gorm.DB) *followerDao.FollowerDao { return &followerDao.FollowerDao{} }).Build()
			mockey.Mock((*userDao.UserDao).IsUserExists).To(func(ctx context.Context, userID int64) (bool, error) { return tc.mockUserExists, nil }).Build()
			mockey.Mock((*followerDao.FollowerDao).GetFollowing).To(func(ctx context.Context, userID int64, pageSize int, pageNum int) ([]*modelDao.User, int, error) {
				return tc.mockUsers, tc.mockTotal, nil
			}).Build()
			mockey.Mock(NewSocialService).To(func(_ context.Context) *SocialService { return &SocialService{} }).Build()
			result, total, err := NewSocialService(context.Background()).ListFollowing(tc.req)
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, UsersToSocialUsers(tc.mockUsers), result)
			assert.Equal(t, tc.mockTotal, total)
		})
	}
}

func TestSocialService_ListFollower(t *testing.T) {
	type testCase struct {
		req            *social.ListFollowerReq
		mockUserExists bool
		mockUsers      []*modelDao.User
		mockTotal      int
		expectError    string
	}

	testCases := map[string]testCase{
		"invalid user id": {req: &social.ListFollowerReq{UserID: "bad"}, expectError: "invalid syntax"},
		"user not exists": {req: &social.ListFollowerReq{UserID: "1"}, expectError: errno.UserIsNotExistErr.ErrMsg},
		"success":         {req: &social.ListFollowerReq{UserID: "1"}, mockUserExists: true, mockUsers: []*modelDao.User{}, mockTotal: 0},
	}

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*gorm.DB).Transaction).To(func(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
				return fc(&gorm.DB{})
			}).Build()
			mockey.Mock((*userDao.UserDao).WithTx).To(func(tx *gorm.DB) *userDao.UserDao { return &userDao.UserDao{} }).Build()
			mockey.Mock((*followerDao.FollowerDao).WithTx).To(func(tx *gorm.DB) *followerDao.FollowerDao { return &followerDao.FollowerDao{} }).Build()
			mockey.Mock((*userDao.UserDao).IsUserExists).To(func(ctx context.Context, userID int64) (bool, error) { return tc.mockUserExists, nil }).Build()
			mockey.Mock((*followerDao.FollowerDao).GetFollower).To(func(ctx context.Context, userID int64, pageSize int, pageNum int) ([]*modelDao.User, int, error) {
				return tc.mockUsers, tc.mockTotal, nil
			}).Build()
			mockey.Mock(NewSocialService).To(func(_ context.Context) *SocialService { return &SocialService{} }).Build()
			result, total, err := NewSocialService(context.Background()).ListFollower(tc.req)
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, UsersToSocialUsers(tc.mockUsers), result)
			assert.Equal(t, tc.mockTotal, total)
		})
	}
}

func TestSocialService_ListFriend(t *testing.T) {
	type testCase struct {
		req            *social.ListFriendReq
		mockUserExists bool
		mockUsers      []*modelDao.User
		mockTotal      int
		expectError    string
	}

	testCases := map[string]testCase{
		"user not exists": {req: &social.ListFriendReq{}, expectError: errno.UserIsNotExistErr.ErrMsg},
		"success":         {req: &social.ListFriendReq{}, mockUserExists: true, mockUsers: []*modelDao.User{}, mockTotal: 0},
	}

	defer mockey.UnPatchAll()
	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*gorm.DB).Transaction).To(func(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
				return fc(&gorm.DB{})
			}).Build()
			mockey.Mock((*userDao.UserDao).WithTx).To(func(tx *gorm.DB) *userDao.UserDao { return &userDao.UserDao{} }).Build()
			mockey.Mock((*followerDao.FollowerDao).WithTx).To(func(tx *gorm.DB) *followerDao.FollowerDao { return &followerDao.FollowerDao{} }).Build()
			mockey.Mock((*userDao.UserDao).IsUserExists).To(func(ctx context.Context, userID int64) (bool, error) { return tc.mockUserExists, nil }).Build()
			mockey.Mock((*followerDao.FollowerDao).GetFriends).To(func(ctx context.Context, userID int64, pageSize int, pageNum int) ([]*modelDao.User, int, error) {
				return tc.mockUsers, tc.mockTotal, nil
			}).Build()
			mockey.Mock(NewSocialService).To(func(_ context.Context) *SocialService { return &SocialService{} }).Build()
			result, total, err := NewSocialService(context.Background()).ListFriend(tc.req, 1)
			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, UsersToSocialUsers(tc.mockUsers), result)
			assert.Equal(t, tc.mockTotal, total)
		})
	}
}
