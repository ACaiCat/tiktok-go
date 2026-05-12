package videocache

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/constants"
	"github.com/ACaiCat/tiktok-go/pkg/db/model"
)

var testVideos = []*model.Video{
	{ID: 1, Title: "video1", VideoURL: "https://example.com/1"},
	{ID: 2, Title: "video2", VideoURL: "https://example.com/2"},
}

func TestVideoCache_SetPopularVideos(t *testing.T) {
	type testCase struct {
		videos  []*model.Video
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"set popular videos success": {
			videos:  testVideos,
			wantErr: false,
		},
		"redis error returns error": {
			videos:  testVideos,
			mockErr: assert.AnError,
			wantErr: true,
		},
		"empty video list success": {
			videos:  []*model.Video{},
			wantErr: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewVideoCache(db)

			data, _ := json.Marshal(tc.videos)
			key := getPopularVideoKey()
			if tc.mockErr != nil {
				mock.ExpectSet(key, data, constants.PopularVideoCacheExpiration).SetErr(tc.mockErr)
			} else {
				mock.ExpectSet(key, data, constants.PopularVideoCacheExpiration).SetVal("OK")
			}

			err := cache.SetPopularVideos(context.Background(), tc.videos)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestVideoCache_GetPopularVideos(t *testing.T) {
	type testCase struct {
		stored  []*model.Video
		mockErr error
		wantErr bool
	}

	testCases := map[string]testCase{
		"get popular videos success": {
			stored:  testVideos,
			wantErr: false,
		},
		"cache miss returns error": {
			mockErr: assert.AnError,
			wantErr: true,
		},
		"empty list returns empty slice": {
			stored:  []*model.Video{},
			wantErr: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			cache := NewVideoCache(db)

			key := getPopularVideoKey()
			if tc.mockErr != nil {
				mock.ExpectGet(key).SetErr(tc.mockErr)
			} else {
				data, _ := json.Marshal(tc.stored)
				mock.ExpectGet(key).SetVal(string(data))
			}

			videos, err := cache.GetPopularVideos(context.Background())
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, videos, len(tc.stored))
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
