package usercache

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
)

func TestUserCache_SetJwchSession(t *testing.T) {
	type testCase struct {
		userID  int64
		jwchID  string
		cookie  string
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"set jwch session success": {
			userID:  1,
			jwchID:  "jwch-id-001",
			cookie:  "session=abc",
			wantErr: false,
		},
		"pipeline error returns error": {
			userID:  1,
			jwchID:  "jwch-id-001",
			cookie:  "session=abc",
			mockErr: assert.AnError,
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewUserCache(db)

			key := getJwchSessionKey(tc.userID)
			session := &jwchSession{ID: tc.jwchID, Cookie: tc.cookie}
			data, _ := json.Marshal(session)

			if tc.mockErr != nil {
				mock.ExpectSet(key, data, constants.JwchSessionCacheExpiration).SetErr(tc.mockErr)
			} else {
				mock.ExpectSet(key, data, constants.JwchSessionCacheExpiration).SetVal("OK")
			}

			err := cache.SetJwchSession(context.Background(), tc.userID, tc.jwchID, tc.cookie)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserCache_CleanJwchSession(t *testing.T) {
	type testCase struct {
		userID  int64
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"clean jwch session success": {
			userID:  1,
			wantErr: false,
		},
		"redis error returns error": {
			userID:  1,
			mockErr: assert.AnError,
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewUserCache(db)

			key := getJwchSessionKey(tc.userID)
			if tc.mockErr != nil {
				mock.ExpectDel(key).SetErr(tc.mockErr)
			} else {
				mock.ExpectDel(key).SetVal(1)
			}

			err := cache.CleanJwchSession(context.Background(), tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
