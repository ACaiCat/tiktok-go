package userdao

import (
	"context"

	"github.com/bytedance/mockey"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func newTestDao() *UserDao {
	return NewUserDao(&gorm.DB{})
}

func mockUserQueryChain() {
	mockey.Mock((*gen.DO).UseDB).To(func(_ *gen.DO, _ *gorm.DB, _ ...gen.DOOption) {}).Build()
	mockey.Mock((*gen.DO).UseModel).To(func(_ *gen.DO, _ interface{}) {}).Build()
	mockey.Mock((*gen.DO).WithContext).To(func(do *gen.DO, _ context.Context) gen.Dao {
		return do
	}).Build()
	mockey.Mock((*gen.DO).Select).To(func(do *gen.DO, _ ...field.Expr) gen.Dao {
		return do
	}).Build()
	mockey.Mock((*gen.DO).Where).To(func(do *gen.DO, _ ...gen.Condition) gen.Dao {
		return do
	}).Build()
	mockey.Mock((*gen.DO).Limit).To(func(do *gen.DO, _ int) gen.Dao {
		return do
	}).Build()
}
