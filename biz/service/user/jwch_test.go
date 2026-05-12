package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/west2-online/jwch"

	"github.com/ACaiCat/tiktok-go/biz/model/user"
	userCache "github.com/ACaiCat/tiktok-go/pkg/cache/user"
	"github.com/ACaiCat/tiktok-go/pkg/crypt"
	modelDao "github.com/ACaiCat/tiktok-go/pkg/db/model"
	userDao "github.com/ACaiCat/tiktok-go/pkg/db/user"
	"github.com/ACaiCat/tiktok-go/pkg/errno"
)

func TestUserService_BindJwch(t *testing.T) {
	type testCase struct {
		req           *user.BindJwchReq
		mockLoginErr  error
		mockCryptErr  error
		mockUpdateErr error
		expectError   string
	}

	testCases := map[string]testCase{
		"success": {
			req: &user.BindJwchReq{JwchID: "233", JwchPassword: "pwd"},
		},
		"jwch login error": {
			req:          &user.BindJwchReq{JwchID: "233", JwchPassword: "pwd"},
			mockLoginErr: assert.AnError,
			expectError:  errno.JwchLoginFailedErr.ErrMsg,
		},
		"encrypt error": {
			req:          &user.BindJwchReq{JwchID: "233", JwchPassword: "pwd"},
			mockCryptErr: assert.AnError,
			expectError:  assert.AnError.Error(),
		},
		"update error": {
			req:           &user.BindJwchReq{JwchID: "233", JwchPassword: "pwd"},
			mockUpdateErr: assert.AnError,
			expectError:   assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*jwch.Student).Login).To(func(_ *jwch.Student) error {
				return tc.mockLoginErr
			}).Build()
			mockey.Mock(crypt.Encrypt).To(func(password string) (string, error) {
				return "encrypted", tc.mockCryptErr
			}).Build()
			mockey.Mock((*userDao.UserDao).UpdateUserJwch).To(
				func(_ *userDao.UserDao, ctx context.Context, userID int64, jwchID string, password string) error {
					return tc.mockUpdateErr
				}).Build()
			mockey.Mock((*userCache.UserCache).CleanJwchSession).To(func(ctx context.Context, userID int64) error {
				return nil
			}).Build()

			mockey.Mock(NewUserService).To(func(_ context.Context) *UserService {
				return &UserService{cache: userCache.NewUserCache(redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"}))}
			}).Build()

			err := NewUserService(t.Context()).BindJwch(tc.req, 1)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestUserService_GetJwchIdentifierAndCookies(t *testing.T) {
	type testCase struct {
		mockCacheID       string
		mockCacheCookie   string
		mockCacheErr      error
		mockSessionErr    error
		mockUser          *modelDao.User
		mockGetErr        error
		mockDecryptErr    error
		mockLoginErr      error
		mockIdentifierErr error
		expectError       string
	}

	testCases := map[string]testCase{
		"success from cache": {
			mockCacheID:     "id",
			mockCacheCookie: "a=b",
		},
		"success from db": {
			mockCacheErr: redis.Nil,
			mockUser: &modelDao.User{
				ID:           1,
				Username:     "u",
				JwchID:       new("233"),
				JwchPassword: new("cipher"),
			},
		},
		"db get error": {
			mockCacheErr: redis.Nil,
			mockGetErr:   assert.AnError,
			expectError:  assert.AnError.Error(),
		},
		"user not found": {
			mockCacheErr: redis.Nil,
			expectError:  errno.UserIsNotExistErr.ErrMsg,
		},
		"not bind": {
			mockCacheErr: redis.Nil,
			mockUser:     &modelDao.User{ID: 1, Username: "u"},
			expectError:  errno.JwchNotBindErr.ErrMsg,
		},
		"decrypt error": {
			mockCacheErr:   redis.Nil,
			mockUser:       &modelDao.User{ID: 1, Username: "u", JwchID: new("233"), JwchPassword: new("cipher")},
			mockDecryptErr: assert.AnError,
			expectError:    assert.AnError.Error(),
		},
		"login error": {
			mockCacheErr: redis.Nil,
			mockUser:     &modelDao.User{ID: 1, Username: "u", JwchID: new("233"), JwchPassword: new("cipher")},
			mockLoginErr: assert.AnError,
			expectError:  errno.JwchLoginFailedErr.ErrMsg,
		},
		"identifier error": {
			mockCacheErr:      redis.Nil,
			mockUser:          &modelDao.User{ID: 1, Username: "u", JwchID: new("233"), JwchPassword: new("cipher")},
			mockIdentifierErr: assert.AnError,
			expectError:       assert.AnError.Error(),
		},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockey.Mock((*userCache.UserCache).GetJwchSession).To(
				func(_ *userCache.UserCache, ctx context.Context, userID int64) (string, string, error) {
					return tc.mockCacheID, tc.mockCacheCookie, tc.mockCacheErr
				}).Build()
			mockey.Mock((*jwch.Student).WithLoginData).To(func(s *jwch.Student, identifier string, cookies []*http.Cookie) *jwch.Student {
				return s
			}).Build()
			mockey.Mock((*jwch.Student).CheckSession).To(func(_ *jwch.Student) error {
				return tc.mockSessionErr
			}).Build()
			mockey.Mock((*userDao.UserDao).GetByID).To(
				func(_ *userDao.UserDao, ctx context.Context, userID int64) (*modelDao.User, error) {
					return tc.mockUser, tc.mockGetErr
				}).Build()
			mockey.Mock(crypt.Decrypt).To(func(cipher string) (string, error) {
				return "pwd", tc.mockDecryptErr
			}).Build()
			mockey.Mock((*jwch.Student).WithUser).To(func(s *jwch.Student, id, password string) *jwch.Student {
				return s
			}).Build()
			mockey.Mock((*jwch.Student).Login).To(func(_ *jwch.Student) error {
				return tc.mockLoginErr
			}).Build()
			mockey.Mock((*jwch.Student).GetIdentifierAndCookies).To(func() (string, []*http.Cookie, error) {
				return "identifier", []*http.Cookie{{Name: "a", Value: "b"}}, tc.mockIdentifierErr
			}).Build()
			mockey.Mock((*userCache.UserCache).SetJwchSession).To(func(ctx context.Context, userID int64, jwchID string, cookie string) error {
				return nil
			}).Build()

			mockey.Mock(NewUserService).To(func(_ context.Context) *UserService {
				return &UserService{cache: userCache.NewUserCache(redis.NewClient(&redis.Options{Addr: "127.0.0.1:0"}))}
			}).Build()

			identifier, cookies, err := NewUserService(t.Context()).GetJwchIdentifierAndCookies(1)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, identifier)
			assert.NotEmpty(t, cookies)
		})
	}
}
