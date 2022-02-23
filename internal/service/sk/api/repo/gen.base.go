package repo

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var globalIsRelated bool = true // 全局预加载

// prepare for other
type _BaseMgr struct {
	*gorm.DB
	ctx       context.Context
	cancel    context.CancelFunc
	timeout   time.Duration
	isRelated bool
}

// SetTimeOut set timeout
func (obj *_BaseMgr) SetTimeOut(timeout time.Duration) {
	obj.ctx, obj.cancel = context.WithTimeout(obj.ctx, timeout)
	obj.timeout = timeout
}

// Cancel cancel context
func (obj *_BaseMgr) Cancel(c context.Context) {
	obj.cancel()
}

// GetDB get gorm.DB info
func (obj *_BaseMgr) GetDB() *gorm.DB {
	return obj.DB
}

// UpdateDB update gorm.DB info
func (obj *_BaseMgr) UpdateDB(db *gorm.DB) {
	obj.DB = db
}

// GetIsRelated Query foreign key Association.获取是否查询外键关联(gorm.Related)
func (obj *_BaseMgr) GetIsRelated() bool {
	return obj.isRelated
}

// SetIsRelated Query foreign key Association.设置是否查询外键关联(gorm.Related)
func (obj *_BaseMgr) SetIsRelated(b bool) {
	obj.isRelated = b
}

// New new gorm.新gorm,重置条件
func (obj *_BaseMgr) new() {
	obj.DB = obj.newDB()
}

// NewDB new gorm.新gorm
func (obj *_BaseMgr) newDB() *gorm.DB {
	return obj.DB.Session(&gorm.Session{NewDB: true, Context: obj.ctx})
}

type options struct {
	query map[string]interface{}
}

// Option overrides behavior of Connect.
type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

// OpenRelated 打开全局预加载
func OpenRelated() {
	globalIsRelated = true
}

// CloseRelated 关闭全局预加载
func CloseRelated() {
	globalIsRelated = true
}

// -------- sql where helper ----------

type CheckWhere func(v interface{}) bool
type DoWhere func(*gorm.DB, interface{}) *gorm.DB

// AddWhere
// CheckWhere 函数 如果返回true，则表明 DoWhere 的查询条件需要加到sql中去
func (obj *_BaseMgr) addWhere(v interface{}, c CheckWhere, d DoWhere) *_BaseMgr {
	if c(v) {
		obj.DB = d(obj.DB, v)
	}
	return obj
}

func (obj *_BaseMgr) sort(userSort, defaultSort string) *_BaseMgr {
	if len(userSort) > 0 {
		obj.DB = obj.DB.Order(userSort)
	} else {
		if len(defaultSort) > 0 {
			obj.DB = obj.DB.Order(defaultSort)
		}
	}
	return obj
}
