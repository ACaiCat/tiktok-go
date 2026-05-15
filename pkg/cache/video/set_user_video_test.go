package videocache

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func TestVideoCache_SetUserVideoList(t *testing.T) {
	type testCase struct {
		userID        int64
		pageSize      int
		pageNum       int
		total         int64
		videos        []*model.Video
		version       string
		versionErr    error
		setErr        error
		wantErr       bool
		wantErrString string
	}

	testCases := map[string]testCase{
		"set user video list success": {
			userID:   10,
			pageSize: 20,
			pageNum:  1,
			total:    9,
			videos:   []*model.Video{{ID: 2}, {ID: 1}},
			version:  "3",
		},
		"missing version uses zero version": {
			userID:     10,
			pageSize:   20,
			pageNum:    0,
			total:      1,
			videos:     []*model.Video{{ID: 1}},
			versionErr: redis.Nil,
		},
		"invalid version returns error": {
			userID:        10,
			pageSize:      20,
			pageNum:       1,
			total:         9,
			videos:        []*model.Video{{ID: 2}, {ID: 1}},
			version:       "bad",
			wantErr:       true,
			wantErrString: "getUserVideoListVersion parse failed",
		},
		"redis set error returns error": {
			userID:        10,
			pageSize:      20,
			pageNum:       1,
			total:         9,
			videos:        []*model.Video{{ID: 2}, {ID: 1}},
			version:       "3",
			setErr:        assert.AnError,
			wantErr:       true,
			wantErrString: "SetUserVideoList failed",
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
				data := expectedUserVideoListJSON(tc.videos, tc.total)
				if tc.setErr != nil {
					mock.ExpectSet(key, data, constants.UserVideoCacheExpiration).SetErr(tc.setErr)
				} else {
					mock.ExpectSet(key, data, constants.UserVideoCacheExpiration).SetVal("OK")
				}
			}

			err := cache.SetUserVideoList(context.Background(), tc.userID, tc.pageSize, tc.pageNum, tc.total, tc.videos)
			if tc.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.wantErrString)
				assert.NoError(t, mock.ExpectationsWereMet())
				return
			}

			require.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestVideoCache_ClearUserVideoList(t *testing.T) {
	testCases := map[string]struct {
		userID        int64
		mockErr       error
		wantErr       bool
		wantErrString string
	}{
		"clear user video list success": {
			userID: 10,
		},
		"redis error returns error": {
			userID:        10,
			mockErr:       assert.AnError,
			wantErr:       true,
			wantErrString: "ClearUserVideoList failed",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewVideoCache(db)

			key := getUserVideoListVersionKey(tc.userID)
			if tc.mockErr != nil {
				mock.ExpectIncr(key).SetErr(tc.mockErr)
			} else {
				mock.ExpectIncr(key).SetVal(1)
			}

			err := cache.ClearUserVideoList(context.Background(), tc.userID)
			if tc.wantErr {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.wantErrString)
				assert.NoError(t, mock.ExpectationsWereMet())
				return
			}

			require.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func expectedUserVideoListJSON(videos []*model.Video, total int64) []byte {
	ids := make([]int64, 0, len(videos))
	for _, video := range videos {
		ids = append(ids, video.ID)
	}
	data, _ := json.Marshal(userVideoListCache{VideoIDs: ids, Total: total})
	return data
}
