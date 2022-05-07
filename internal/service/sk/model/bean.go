package model

import (
	"time"
)

type SKListResp struct {
	ID       int `json:"id"`
	ShopId   int `json:"shop_id"`
	ShopInfo struct {
		ShopName  string `json:"shop_name"`
		ShopDesc  string `json:"shop_desc"`
		ShopCount int    `json:"shop_count"`
	}
	StartTime  int `json:"start_time"`
	EndTime    int `json:"end_time"`
	Status     int `json:"status"`
	CreateTime int `json:"create_time"`
}

type SKAdd struct {
	ShopId    int       `form:"shop_id" v:"required#请输入正确的ID"`
	StartTime time.Time `form:"start_time" v:"required" time_format:"2006-01-02 15:04:05" time_location:"Asia/Shanghai"` // 秒杀活动的开始时间至少要在当前时间的30秒之后，也就是强制必须定时。
	EndTime   time.Time `form:"end_time" v:"required" time_format:"2006-01-02 15:04:05" time_location:"Asia/Shanghai"`   // 最晚截止时间，如果商品提前被秒杀完毕，则更新为当前时间
}

type SKStatus int

const (
	SKWait     = 1  // 待开始
	SKShopping = 2  // 秒杀中
	SKFinish   = 3  // 已完成（结束时间已到，或者商品被秒杀完毕）
	SKClose    = -1 // 下架（状态为1和2的可下架）
)

const SKDelayAddTag = "SKDelayAddTag"
