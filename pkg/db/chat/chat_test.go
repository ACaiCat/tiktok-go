package chatdao

import (
	"gorm.io/gorm"

	dbtestutil "github.com/ACaiCat/tiktok-go/pkg/db/testutil"
)

func newTestDao() *ChatDao {
	return NewChatDao(&gorm.DB{})
}

func mockChatQueryChain() {
	dbtestutil.MockDOChain()
}
