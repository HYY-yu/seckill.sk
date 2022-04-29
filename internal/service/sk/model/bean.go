package model

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
	ShopId    int `json:"shop_id" v:"required|min:1#请输入正确的ID"`
	StartTime int `json:"start_time" v:"required"`
	EndTime   int `json:"end_time" v:"required"`
}
