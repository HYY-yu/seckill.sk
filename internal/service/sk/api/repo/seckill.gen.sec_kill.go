package repo

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/HYY-yu/seckill.sk/internal/service/sk/model"
)

type _SecKillMgr struct {
	*_BaseMgr
}

// SecKillMgr open func
func SecKillMgr(db *gorm.DB) *_SecKillMgr {
	if db == nil {
		panic(fmt.Errorf("SecKillMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &_SecKillMgr{_BaseMgr: &_BaseMgr{DB: db.Table("sec_kill"), isRelated: globalIsRelated, ctx: ctx, cancel: cancel, timeout: -1}}
}

// WithContext set context to db
func (obj *_SecKillMgr) WithContext(c context.Context) *_SecKillMgr {
	if c != nil {
		obj.ctx = c
	}
	return obj
}

// GetTableName get sql table name.获取数据库名字
func (obj *_SecKillMgr) GetTableName() string {
	return "sec_kill"
}

// Reset 重置gorm会话
func (obj *_SecKillMgr) Reset() *_SecKillMgr {
	obj.New()
	return obj
}

// Get 获取
func (obj *_SecKillMgr) Get() (result model.SecKill, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_SecKillMgr) Gets() (results []*model.SecKill, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Find(&results).Error

	return
}

////////////////////////////////// gorm replace /////////////////////////////////
func (obj *_SecKillMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Count(count)
}

//////////////////////////////////////////////////////////////////////////////////

//////////////////////////option case ////////////////////////////////////////////

// WithID id获取
func (obj *_SecKillMgr) WithID(id int) Option {
	return optionFunc(func(o *options) { o.query["id"] = id })
}

// WithShopID shop_id获取
func (obj *_SecKillMgr) WithShopID(shopID int) Option {
	return optionFunc(func(o *options) { o.query["shop_id"] = shopID })
}

// WithStartTime start_time获取
func (obj *_SecKillMgr) WithStartTime(startTime int) Option {
	return optionFunc(func(o *options) { o.query["start_time"] = startTime })
}

// WithEndTime end_time获取
func (obj *_SecKillMgr) WithEndTime(endTime int) Option {
	return optionFunc(func(o *options) { o.query["end_time"] = endTime })
}

// WithStatus status获取
func (obj *_SecKillMgr) WithStatus(status int8) Option {
	return optionFunc(func(o *options) { o.query["status"] = status })
}

// WithCreateTime create_time获取
func (obj *_SecKillMgr) WithCreateTime(createTime int) Option {
	return optionFunc(func(o *options) { o.query["create_time"] = createTime })
}

// GetByOption 功能选项模式获取
func (obj *_SecKillMgr) GetByOption(opts ...Option) (result model.SecKill, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Where(options.query).Find(&result).Error

	return
}

// GetByOptions 批量功能选项模式获取
func (obj *_SecKillMgr) GetByOptions(opts ...Option) (results []*model.SecKill, err error) {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Where(options.query).Find(&results).Error

	return
}

//////////////////////////enume case ////////////////////////////////////////////

// GetFromID 通过id获取内容
func (obj *_SecKillMgr) GetFromID(id int) (result model.SecKill, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// GetBatchFromID 批量查找
func (obj *_SecKillMgr) GetBatchFromID(ids []int) (results []*model.SecKill, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Where("`id` IN (?)", ids).Find(&results).Error

	return
}

// GetFromShopID 通过shop_id获取内容
func (obj *_SecKillMgr) GetFromShopID(shopID int) (results []*model.SecKill, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Where("`shop_id` = ?", shopID).Find(&results).Error

	return
}

// GetBatchFromShopID 批量查找
func (obj *_SecKillMgr) GetBatchFromShopID(shopIDs []int) (results []*model.SecKill, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Where("`shop_id` IN (?)", shopIDs).Find(&results).Error

	return
}

// GetFromStartTime 通过start_time获取内容
func (obj *_SecKillMgr) GetFromStartTime(startTime int) (results []*model.SecKill, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Where("`start_time` = ?", startTime).Find(&results).Error

	return
}

// GetBatchFromStartTime 批量查找
func (obj *_SecKillMgr) GetBatchFromStartTime(startTimes []int) (results []*model.SecKill, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Where("`start_time` IN (?)", startTimes).Find(&results).Error

	return
}

// GetFromEndTime 通过end_time获取内容
func (obj *_SecKillMgr) GetFromEndTime(endTime int) (results []*model.SecKill, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Where("`end_time` = ?", endTime).Find(&results).Error

	return
}

// GetBatchFromEndTime 批量查找
func (obj *_SecKillMgr) GetBatchFromEndTime(endTimes []int) (results []*model.SecKill, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Where("`end_time` IN (?)", endTimes).Find(&results).Error

	return
}

// GetFromStatus 通过status获取内容
func (obj *_SecKillMgr) GetFromStatus(status int8) (results []*model.SecKill, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Where("`status` = ?", status).Find(&results).Error

	return
}

// GetBatchFromStatus 批量查找
func (obj *_SecKillMgr) GetBatchFromStatus(statuss []int8) (results []*model.SecKill, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Where("`status` IN (?)", statuss).Find(&results).Error

	return
}

// GetFromCreateTime 通过create_time获取内容
func (obj *_SecKillMgr) GetFromCreateTime(createTime int) (results []*model.SecKill, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Where("`create_time` = ?", createTime).Find(&results).Error

	return
}

// GetBatchFromCreateTime 批量查找
func (obj *_SecKillMgr) GetBatchFromCreateTime(createTimes []int) (results []*model.SecKill, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Where("`create_time` IN (?)", createTimes).Find(&results).Error

	return
}

//////////////////////////primary index case ////////////////////////////////////////////

// FetchByPrimaryKey primary or index 获取唯一内容
func (obj *_SecKillMgr) FetchByPrimaryKey(id int) (result model.SecKill, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.SecKill{}).Where("`id` = ?", id).Find(&result).Error

	return
}