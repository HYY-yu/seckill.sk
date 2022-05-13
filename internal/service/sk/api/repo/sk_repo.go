package repo

import (
	"context"

	"gorm.io/gorm"

	"github.com/HYY-yu/seckill.pkg/pkg/util"

	"github.com/HYY-yu/seckill.sk/internal/service/sk/model"
)

type SKRepo interface {
	Mgr(ctx context.Context, db *gorm.DB) *_SecKillMgr
}

type skRepo struct {
}

func NewSKRepo() SKRepo {
	return &skRepo{}
}

func (*skRepo) Mgr(ctx context.Context, db *gorm.DB) *_SecKillMgr {
	skMgr := SecKillMgr(ctx, db)
	return skMgr
}

// ------- 自定义方法 -------

func (obj *_SecKillMgr) ListSK(
	limit, offset int,
	filter map[string]interface{},
	sort string,
) (result []model.SecKill, err error) {
	err = obj.
		addWhere(filter[model.SecKillColumns.ID], util.IsNotZero, func(db *gorm.DB, i interface{}) *gorm.DB {
			return db.Where(model.SecKillColumns.ID+" = ?", i)
		}).
		addWhere(filter[model.SecKillColumns.Status], util.IsNotZero, func(db *gorm.DB, i interface{}) *gorm.DB {
			return db.Where(model.SecKillColumns.Status + " > 0")
		}).
		sort(sort, model.SecKillColumns.ID+" desc").
		Limit(limit).
		Offset(offset).
		Find(&result).Error
	return
}

func (obj *_SecKillMgr) CountSK(
	filter map[string]interface{},
) (count int64, err error) {
	err = obj.
		addWhere(filter[model.SecKillColumns.ID], util.IsNotZero, func(db *gorm.DB, i interface{}) *gorm.DB {
			return db.Where(model.SecKillColumns.ID+" = ?", i)
		}).
		addWhere(filter[model.SecKillColumns.Status], util.IsNotZero, func(db *gorm.DB, i interface{}) *gorm.DB {
			return db.Where(model.SecKillColumns.Status + " > 0")
		}).
		Count(&count).Error
	return
}
