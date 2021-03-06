package repo

import (
	"context"
	"fmt"
	"gorm.io/gorm"

	"github.com/HYY-yu/seckill.sk/internal/service/sk/model"
)

// Code generated by gormt. DO NOT EDIT.

type _UsersMgr struct {
	*_BaseMgr
}

// UsersMgr open func
func UsersMgr(ctx context.Context, db *gorm.DB) *_UsersMgr {
	if db == nil {
		panic(fmt.Errorf("UsersMgr need init by db"))
	}
	ctx, cancel := context.WithCancel(ctx)
	return &_UsersMgr{_BaseMgr: &_BaseMgr{DB: db.Table("users"), isRelated: globalIsRelated, ctx: ctx, cancel: cancel, timeout: -1}}
}

func (obj *_UsersMgr) WithSelects(idName string, selects ...string) *_UsersMgr {
	if len(selects) > 0 {
		if len(idName) > 0 {
			selects = append(selects, idName)
		}
		// 对Select进行去重
		selectMap := make(map[string]int, len(selects))
		for _, e := range selects {
			if _, ok := selectMap[e]; !ok {
				selectMap[e] = 1
			}
		}

		newSelects := make([]string, 0, len(selects))
		for k := range selectMap {
			if len(k) > 0 {
				newSelects = append(newSelects, k)
			}
		}
		obj.DB = obj.DB.Select(newSelects)
	}
	return obj
}

func (obj *_UsersMgr) WithOmit(omit ...string) *_UsersMgr {
	if len(omit) > 0 {
		obj.DB = obj.DB.Omit(omit...)
	}
	return obj
}

func (obj *_UsersMgr) WithOptions(opts ...Option) *_UsersMgr {
	options := options{
		query: make(map[string]interface{}, len(opts)),
	}
	for _, o := range opts {
		o.apply(&options)
	}
	obj.DB = obj.DB.Where(options.query)
	return obj
}

// GetTableName get sql table name.获取数据库名字
func (obj *_UsersMgr) GetTableName() string {
	return "users"
}

// Reset 重置gorm会话
func (obj *_UsersMgr) Reset() *_UsersMgr {
	obj.new()
	return obj
}

// Get 获取
func (obj *_UsersMgr) Get() (result model.Users, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Users{}).Find(&result).Error

	return
}

// Gets 获取批量结果
func (obj *_UsersMgr) Gets() (results []*model.Users, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Users{}).Find(&results).Error

	return
}

func (obj *_UsersMgr) Count(count *int64) (tx *gorm.DB) {
	return obj.DB.WithContext(obj.ctx).Model(model.Users{}).Count(count)
}

func (obj *_UsersMgr) HasRecord() (bool, error) {
	var count int64
	err := obj.DB.WithContext(obj.ctx).Model(model.Users{}).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

// WithID id获取
func (obj *_UsersMgr) WithID(id int) Option {
	return optionFunc(func(o *options) { o.query["id"] = id })
}

// WithUserName user_name获取
func (obj *_UsersMgr) WithUserName(userName string) Option {
	return optionFunc(func(o *options) { o.query["user_name"] = userName })
}

// GetFromID 通过id获取内容
func (obj *_UsersMgr) GetFromID(id int) (result model.Users, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Users{}).Where("`id` = ?", id).Find(&result).Error

	return
}

// GetBatchFromID 批量查找
func (obj *_UsersMgr) GetBatchFromID(ids []int) (results []*model.Users, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Users{}).Where("`id` IN (?)", ids).Find(&results).Error

	return
}

// GetFromUserName 通过user_name获取内容
func (obj *_UsersMgr) GetFromUserName(userName string) (result model.Users, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Users{}).Where("`user_name` = ?", userName).Find(&result).Error

	return
}

// GetBatchFromUserName 批量查找
func (obj *_UsersMgr) GetBatchFromUserName(userNames []string) (results []*model.Users, err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Users{}).Where("`user_name` IN (?)", userNames).Find(&results).Error

	return
}

func (obj *_UsersMgr) CreateUsers(bean *model.Users) (err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Users{}).Create(bean).Error

	return
}

func (obj *_UsersMgr) UpdateUsers(bean *model.Users) (err error) {
	err = obj.DB.WithContext(obj.ctx).Model(bean).Updates(bean).Error

	return
}

func (obj *_UsersMgr) DeleteUsers(bean *model.Users) (err error) {
	err = obj.DB.WithContext(obj.ctx).Model(model.Users{}).Delete(bean).Error

	return
}
