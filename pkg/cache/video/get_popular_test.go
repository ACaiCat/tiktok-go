package videocache

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

var testVideos = []*model.Video{
	{ID: 1, Title: "video1", VideoURL: "https://example.com/1"},
	{ID: 2, Title: "video2", VideoURL: "https://example.com/2"},
}

func TestVideoCache_GetPopularVideos(t *testing.T) {
	type testCase struct {
		pageSize      int
		pageNum       int
		storedIDs     []string
		wantIDs       []int64
		mockErr       error
		wantErr       bool
		wantErrString string
	}

	testCases := map[string]testCase{
		"get popular videos success": {
			pageSize:  2,
			pageNum:   1,
			storedIDs: []string{"2", "1"},
			wantIDs:   []int64{2, 1},
		},
		"cache miss returns error": {
			pageSize:      1,
			pageNum:       0,
			mockErr:       assert.AnError,
			wantErr:       true,
			wantErrString: "GetPopularVideos failed",
		},
		"empty zset returns error": {
			pageSize:      10,
			pageNum:       0,
			storedIDs:     []string{},
			wantErr:       true,
			wantErrString: redis.Nil.Error(),
		},
		"invalid video id returns error": {
			pageSize:      1,
			pageNum:       0,
			storedIDs:     []string{"bad"},
			wantErr:       true,
			wantErrString: "GetPopularVideos parse videoID failed",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewVideoCache(db)

			key := getPopularVideoKey()
			start := int64(tc.pageSize * tc.pageNum)
			stop := start + int64(tc.pageSize) - 1
			args := redis.ZRangeArgs{
				Key:   key,
				Start: start,
				Stop:  stop,
				Rev:   true,
			}
			if tc.mockErr != nil {
				mock.ExpectZRangeArgs(args).SetErr(tc.mockErr)
			} else {
				mock.ExpectZRangeArgs(args).SetVal(tc.storedIDs)
			}

			videoIDs, err := cache.GetPopularVideos(context.Background(), tc.pageSize, tc.pageNum)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.wantErrString)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantIDs, videoIDs)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
