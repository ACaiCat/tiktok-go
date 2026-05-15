package usercache

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestUserCache_GetJwchSession(t *testing.T) {
	type testCase struct {
		userID  int64
		jwchID  string
		cookie  string
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"get jwch session success": {
			userID:  1,
			jwchID:  "jwch-id-001",
			cookie:  "session=abc",
			wantErr: false,
		},
		"key not found returns error": {
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
				mock.ExpectGet(key).SetErr(tc.mockErr)
			} else {
				session := &jwchSession{ID: tc.jwchID, Cookie: tc.cookie}
				data, _ := json.Marshal(session)
				mock.ExpectGet(key).SetVal(string(data))
			}

			id, cookie, err := cache.GetJwchSession(context.Background(), tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.jwchID, id)
			assert.Equal(t, tc.cookie, cookie)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
