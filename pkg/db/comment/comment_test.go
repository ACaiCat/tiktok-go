package commentdao

import (
	"gorm.io/gorm"

	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func newTestDao() *CommentDao {
	return NewCommentDao(&gorm.DB{})
}

func mockCommentQueryChain() {
	dbtestutil.MockDOChain()
}
