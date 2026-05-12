package videodao

import (
	"context"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/ACaiCat/tiktok-go/pkg/db/model"
	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func TestVideoDao_SearchVideo(t *testing.T) {
	type testCase struct {
		keywords []string
		pageSize int
		pageNum  int
		fromDate time.Time
		toDate   time.Time
		username string
		mockRet  []*model.Video
		mockErr  error
		wantErr  bool
	}

	videos := []*model.Video{{ID: 1, Title: "golang tutorial"}}

	testCases := map[string]testCase{
		"search success":   {keywords: []string{"golang"}, pageSize: 10, pageNum: 0, mockRet: videos},
		"no results":       {keywords: []string{"nope"}, pageSize: 10, pageNum: 0, mockRet: []*model.Video{}},
		"db error":         {keywords: []string{"golang"}, pageSize: 10, pageNum: 0, mockErr: assert.AnError, wantErr: true},
		"with date filter": {keywords: []string{"go"}, pageSize: 10, pageNum: 0, fromDate: time.Now().Add(-24 * time.Hour), toDate: time.Now(), mockRet: videos},
	}

	defer mockey.UnPatchAll()

	for name, tc := range testCases {
		mockey.PatchConvey(name, t, func() {
			mockVideoQueryChain()
			dao := newTestDao()
			dbtestutil.MockFind(tc.mockRet, tc.mockErr)

			vs, err := dao.SearchVideo(context.Background(), tc.keywords, tc.pageSize, tc.pageNum, tc.fromDate, tc.toDate, tc.username)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "SearchVideo failed")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.mockRet, vs)
			}
		})
	}
}
