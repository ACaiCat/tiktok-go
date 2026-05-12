package followerdao

import (
	"gorm.io/gorm"

	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func newTestDao() *FollowerDao {
	return NewFollowerDao(&gorm.DB{})
}

func mockFollowerQueryChain() {
	dbtestutil.MockDOChain()
}
