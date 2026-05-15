package videocache

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVideoCache_GetUserVideoList(t *testing.T) {
	type testCase struct {
		userID        int64
		pageSize      int
		pageNum       int
		version       string
		versionErr    error
		stored        string
		listErr       error
		wantIDs       []int64
		wantTotal     int64
		wantErr       bool
		wantErrString string
	}

	testCases := map[string]testCase{
		"get user video list success": {
			userID:    10,
			pageSize:  20,
			pageNum:   1,
			version:   "3",
			stored:    `{"video_ids":[2,1],"total":9}`,
			wantIDs:   []int64{2, 1},
			wantTotal: 9,
		},
		"missing version uses zero version": {
			userID:     10,
			pageSize:   20,
			pageNum:    0,
			versionErr: redis.Nil,
			stored:     `{"video_ids":[1],"total":1}`,
			wantIDs:    []int64{1},
			wantTotal:  1,
		},
		"list cache miss returns error": {
			userID:        10,
			pageSize:      20,
			pageNum:       1,
			version:       "3",
			listErr:       redis.Nil,
			wantErr:       true,
			wantErrString: "GetUserVideoList failed",
		},
		"invalid version returns error": {
			userID:        10,
			pageSize:      20,
			pageNum:       1,
			version:       "bad",
			wantErr:       true,
			wantErrString: "getUserVideoListVersion parse failed",
		},
		"invalid json returns error": {
			userID:        10,
			pageSize:      20,
			pageNum:       1,
			version:       "3",
			stored:        `{"video_ids":`,
			wantErr:       true,
			wantErrString: "GetUserVideoList json failed",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewVideoCache(db)

			versionKey := getUserVideoListVersionKey(tc.userID)
			if tc.versionErr != nil {
				mock.ExpectGet(versionKey).SetErr(tc.versionErr)
			} else {
				mock.ExpectGet(versionKey).SetVal(tc.version)
			}

			if tc.wantErrString != "getUserVideoListVersion parse failed" {
				version := int64(0)
				if tc.version == "3" {
					version = 3
				}
				key := getUserVideoListKey(tc.userID, version, tc.pageSize, tc.pageNum)
				if tc.listErr != nil {
					mock.ExpectGet(key).SetErr(tc.listErr)
				} else {
					mock.ExpectGet(key).SetVal(tc.stored)
				}
			}

			gotIDs, gotTotal, err := cache.GetUserVideoList(context.Background(), tc.userID, tc.pageSize, tc.pageNum)
			if tc.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.wantErrString)
				assert.NoError(t, mock.ExpectationsWereMet())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantIDs, gotIDs)
			assert.Equal(t, tc.wantTotal, gotTotal)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
