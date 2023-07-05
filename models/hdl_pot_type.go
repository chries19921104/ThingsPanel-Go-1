package models

type HdlPotType struct {
	Id           string `json:"id" gorm:"primaryKey"`             // ID
	Name         string `json:"name,omitempty" gorm:"size:500"`   // 锅型名称
	Image        string `json:"image,omitempty"  gorm:"size:500"` // 图片
	CreateAt     int64  `json:"create_at,omitempty"`              // 创建时间
	UpdateAt     int64  `json:"update_at,omitempty"`              // 更新时间
	SoupStandard int64  `json:"soup_standard,omitempty"`          // 加汤水位线标准
	PotTypeId    string `json:"pot_type_id,omitempty"`            // 锅型id
	Remark       string `json:"remark,omitempty" gorm:"size:255"`
}

func (HdlPotType) TableName() string {
	return "hdl_pot_type"
}
