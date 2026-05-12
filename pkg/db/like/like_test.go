package likedao

import (
	"gorm.io/gorm"

	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func newTestDao() *LikeDao {
	return NewLikeDao(&gorm.DB{})
}

func mockLikeQueryChain() {
	dbtestutil.MockDOChain()
}
