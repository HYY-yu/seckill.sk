package model

// SecKill 秒杀表
type SecKill struct {
	ID         int  `gorm:"primaryKey;column:id;type:int(11);not null"`
	ShopID     int  `gorm:"column:shop_id;type:int(11);not null"`
	StartTime  int  `gorm:"column:start_time;type:int(11);not null"`
	EndTime    int  `gorm:"column:end_time;type:int(11);not null;default:0"`
	Status     int8 `gorm:"column:status;type:tinyint(4);not null"`
	CreateTime int  `gorm:"column:create_time;type:int(11);not null"`
}

// SecKillColumns get sql column name.获取数据库列名
var SecKillColumns = struct {
	ID         string
	ShopID     string
	StartTime  string
	EndTime    string
	Status     string
	CreateTime string
}{
	ID:         "id",
	ShopID:     "shop_id",
	StartTime:  "start_time",
	EndTime:    "end_time",
	Status:     "status",
	CreateTime: "create_time",
}
