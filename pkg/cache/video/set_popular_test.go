package videocache

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

func TestVideoCache_SetPopularVideos(t *testing.T) {
	type testCase struct {
		videos        []*model.Video
		mockErr       error
		wantErr       bool
		wantErrString string
	}

	testCases := map[string]testCase{
		"set popular videos success": {
			videos: testVideos,
		},
		"redis error returns error": {
			videos:        testVideos,
			mockErr:       assert.AnError,
			wantErr:       true,
			wantErrString: "SetPopularVideosCache failed",
		},
		"empty video list success": {
			videos: []*model.Video{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewVideoCache(db)

			key := getPopularVideoKey()
			if tc.mockErr != nil {
				mock.ExpectDel(key).SetErr(tc.mockErr)
			} else {
				mock.ExpectDel(key).SetVal(1)
				if len(tc.videos) > 0 {
					members := make([]redis.Z, 0, len(tc.videos))
					for _, video := range tc.videos {
						members = append(members, redis.Z{Score: float64(video.VisitCount), Member: video.ID})
					}
					mock.ExpectZAdd(key, members...).SetVal(int64(len(members)))
					mock.ExpectExpire(key, constants.PopularVideoCacheExpiration).SetVal(true)
				}
			}

			err := cache.SetPopularVideos(context.Background(), tc.videos)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.wantErrString)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestVideoCache_IncrPopularVideoVisitCount(t *testing.T) {
	testCases := map[string]struct {
		videoID       int64
		mockErr       error
		wantErr       bool
		wantErrString string
	}{
		"incr popular video visit count success": {
			videoID: 1,
		},
		"redis error returns error": {
			videoID:       1,
			mockErr:       assert.AnError,
			wantErr:       true,
			wantErrString: "IncrPopularVideoVisitCount failed",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewVideoCache(db)

			if tc.mockErr != nil {
				mock.ExpectZIncrBy(getPopularVideoKey(), 1, "1").SetErr(tc.mockErr)
			} else {
				mock.ExpectZIncrBy(getPopularVideoKey(), 1, "1").SetVal(1)
			}

			err := cache.IncrPopularVideoVisitCount(context.Background(), tc.videoID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.wantErrString)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
