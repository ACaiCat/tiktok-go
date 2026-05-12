package videodao

import (
	"gorm.io/gorm"

	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func newTestDao() *VideoDao {
	return NewVideoDao(&gorm.DB{})
}

func mockVideoQueryChain() {
	dbtestutil.MockDOChain()
}
