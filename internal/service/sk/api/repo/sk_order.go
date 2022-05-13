package repo

import (
	"context"

	"github.com/HYY-yu/seckill.pkg/pkg/util"
	"gorm.io/gorm"

	"github.com/HYY-yu/seckill.sk/internal/service/sk/model"
)

type OrderRepo interface {
	Mgr(ctx context.Context, db *gorm.DB) *_OrderMgr
}

type orderRepo struct {
}

func NewOrderRepo() OrderRepo {
	return &orderRepo{}
}

func (*orderRepo) Mgr(ctx context.Context, db *gorm.DB) *_OrderMgr {
	mgr := OrderMgr(ctx, db)
	return mgr
}

// ------- 自定义方法 -------

func (obj *_OrderMgr) ListOrder(
	limit, offset int,
	filter map[string]interface{},
	sort string,
) (result []model.Order, err error) {
	err = obj.
		addWhere(filter[model.OrderColumns.ID], util.IsNotZero, func(db *gorm.DB, i interface{}) *gorm.DB {
			return db.Where(model.OrderColumns.ID+" = ?", i)
		}).
		sort(sort, model.OrderColumns.ID+" desc").
		Limit(limit).
		Offset(offset).
		Find(&result).Error
	return
}

func (obj *_OrderMgr) CountOrder(
	filter map[string]interface{},
) (count int64, err error) {
	err = obj.
		addWhere(filter[model.OrderColumns.ID], util.IsNotZero, func(db *gorm.DB, i interface{}) *gorm.DB {
			return db.Where(model.OrderColumns.ID+" = ?", i)
		}).
		Count(&count).Error
	return
}
