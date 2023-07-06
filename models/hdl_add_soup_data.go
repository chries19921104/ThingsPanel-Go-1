package models

import "time"

type HdlAddSoupData struct {
	Id                string `json:"id" gorm:"primaryKey"`
	ShopId            string `json:"shop_id" gorm:"size:50"`
	ShopName          string `json:"shop_name" gorm:"size:200"`    //店铺名称
	OrderSn           string `json:"order_sn" gorm:"size:200"`     //订单号
	BottomId          string `json:"bottom_id" gorm:"size:50"`     //锅底id(|分割)
	PotTypeId         string `json:"pot_type_id" gorm:"size:50"`   //锅型id
	BottomPot         string `json:"bottom_pot" gorm:"size:200"`   //锅底名称(|分割)
	TableNumber       string `json:"table_number" gorm:"size:200"` //桌号
	OrderTime         string `json:"order_time"`                   //下单时间
	SoupStartTime     string `json:"soup_start_time"`              //加汤开始时间
	SoupEndTime       string `json:"soup_end_time"`                //加汤结束时间
	FeedingStartTime  string `json:"feeding_start_time"`           //投料开始时间
	FeedingEndTime    string `json:"feeding_end_time"`             //投料结束时间
	TurningPotEndTime string `json:"turning_pot_end_time"`         //翻锅结束时间
	CreationTime      string `json:"creation_time"`                //订单创建时间
	CreateAt          int64  `json:"create_at"`                    //创建时间
}

type AddSoupDataValue struct {
	Id               string    `gorm:"column:id;NOT NULL" json:"id,omitempty"`
	ShopName         string    `gorm:"column:name;NOT NULL"`
	OrderSn          string    `gorm:"column:order_sn;NOT NULL"`
	BottomPot        string    `gorm:"column:bottom_pot;NOT NULL"`
	TableNumber      string    `gorm:"column:table_number"`
	OrderTime        time.Time `gorm:"column:order_time"`
	SoupStartTime    time.Time `gorm:"column:soup_start_time"`
	SoupEndTime      time.Time `gorm:"column:soup_end_time"`
	FeedingStartTime time.Time `gorm:"column:feeding_start_time"`
	FeedingEndTime   time.Time `gorm:"column:feeding_end_time"`
	TurningPotEnd    time.Time `gorm:"column:turning_pot_end_time"`
}

type ReturnAddSoupDataValue struct {
	Id               string `gorm:"column:id;NOT NULL" json:"id,omitempty"`
	ShopName         string `gorm:"column:name;NOT NULL"`
	OrderSn          string `gorm:"column:order_sn;NOT NULL"`
	BottomPot        string `gorm:"column:bottom_pot;NOT NULL"`
	TableNumber      string `gorm:"column:table_number"`
	OrderTime        int64  `gorm:"column:order_time"`
	SoupStartTime    int64  `gorm:"column:soup_start_time"`
	SoupEndTime      int64  `gorm:"column:soup_end_time"`
	FeedingStartTime int64  `gorm:"column:feeding_start_time"`
	FeedingEndTime   int64  `gorm:"column:feeding_end_time"`
	TurningPotEnd    int64  `gorm:"column:turning_pot_end_time"`
}

func (a *HdlAddSoupData) TableName() string {
	return "hdl_add_soup_data"
}

type SoupDataKVResult struct {
	AddSoupDataValue
	Name     string `json:"name"`
	PluginId string `json:"plugin_id"`
}
