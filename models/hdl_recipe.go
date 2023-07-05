package models

type HdlRecipe struct {
	Id               string `json:"id" gorm:"primaryKey"`                        // ID
	BottomPotId      string `json:"bottom_pot_id,omitempty" gorm:"size:500"`     // 锅底id
	BottomPot        string `json:"bottom_pot,omitempty"  gorm:"size:500"`       // 锅底名称
	BottomProperties string `json:"bottom_properties,omitempty"  gorm:"size:50"` // 锅底属性 辣 不辣
	HdlPotTypeId     string `json:"hdl_pot_type_id,omitempty"  gorm:"size:50"`   // 锅型id
	CreateAt         int64  `json:"create_at,omitempty"`                         // 创建时间
	UpdateAt         int64  `json:"update_at,omitempty"`                         // 更新时间
	TenantId         string `json:"tenant_id,omitempty"  gorm:"size:36"`         // 更新时间
	Remark           string `json:"remark,omitempty" gorm:"size:255"`
}

func (HdlRecipe) TableName() string {
	return "hdl_recipe"
}
