package dbtestutil

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/bytedance/mockey"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/ACaiCat/tiktok-go/pkg/db/query"
)

func MockDOChain() {
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
	mockey.Mock((*gen.DO).Or).To(func(do *gen.DO, _ ...gen.Condition) gen.Dao {
		return do
	}).Build()
	mockey.Mock((*gen.DO).Order).To(func(do *gen.DO, _ ...field.Expr) gen.Dao {
		return do
	}).Build()
	mockey.Mock((*gen.DO).Offset).To(func(do *gen.DO, _ int) gen.Dao {
		return do
	}).Build()
	mockey.Mock((*gen.DO).Limit).To(func(do *gen.DO, _ int) gen.Dao {
		return do
	}).Build()
	mockey.Mock((*gen.DO).Group).To(func(do *gen.DO, _ ...field.Expr) gen.Dao {
		return do
	}).Build()
	mockey.Mock((*gen.DO).Join).To(func(do *gen.DO, _ schema.Tabler, _ ...field.Expr) gen.Dao {
		return do
	}).Build()
}

func MockCreate(err error) {
	mockey.Mock((*gen.DO).Create).Return(err).Build()
}

func MockCreateWithHook(hook func(value interface{}), err error) {
	mockey.Mock((*gen.DO).Create).To(func(_ *gen.DO, value interface{}) error {
		if hook != nil {
			hook(value)
		}
		return err
	}).Build()
}

func MockFirst(ret interface{}, err error) {
	mockey.Mock((*gen.DO).First).Return(ret, err).Build()
}

func MockFind(ret interface{}, err error) {
	mockey.Mock((*gen.DO).Find).Return(ret, err).Build()
}

func MockCount(count int64, err error) {
	mockey.Mock((*gen.DO).Count).Return(count, err).Build()
}

func MockScan(fill func(dest interface{}), err error) {
	mockey.Mock((*gen.DO).Scan).To(func(_ *gen.DO, dest interface{}) error {
		if err != nil {
			return err
		}
		if fill != nil {
			fill(dest)
		}
		return nil
	}).Build()
}

func FillValue(dest interface{}, value interface{}) {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return
	}
	valueOfValue := reflect.ValueOf(value)
	if !valueOfValue.IsValid() {
		destValue.Elem().Set(reflect.Zero(destValue.Elem().Type()))
		return
	}
	if valueOfValue.Type().AssignableTo(destValue.Elem().Type()) {
		destValue.Elem().Set(valueOfValue)
	}
}

func MockUpdate(err error) {
	mockey.Mock((*gen.DO).Update).Return(gen.ResultInfo{}, err).Build()
}

func MockUpdates(err error) {
	mockey.Mock((*gen.DO).Updates).Return(gen.ResultInfo{}, err).Build()
}

func MockUpdateColumn(err error) {
	mockey.Mock((*gen.DO).UpdateColumn).Return(gen.ResultInfo{}, err).Build()
}

func MockDelete(err error) {
	mockey.Mock((*gen.DO).Delete).Return(gen.ResultInfo{}, err).Build()
}

func MockTransaction(err error) {
	mockey.Mock((*query.Query).Transaction).To(func(q *query.Query, fc func(tx *query.Query) error, _ ...*sql.TxOptions) error {
		if err != nil {
			return err
		}
		return fc(q)
	}).Build()
}
