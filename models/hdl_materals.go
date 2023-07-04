package models

type HdlMaterals struct {
	Id        string `json:"id" gorm:"primaryKey"`           // ID
	Name      string `json:"name,omitempty" gorm:"size:500"` // 系统名称
	Dosage    int64  `json:"dosage,omitempty"`               // 主题
	Unit      string `json:"unit,omitempty" gorm:"size:255"` // 首页logo
	WaterLine int64  `json:"water_line,omitempty"`           // 缓冲logo
	Station   string `json:"station,omitempty"  gorm:"size:255"`
	Resource  string `json:"resource,omitempty" gorm:"size:255"`
	CreateAt  int64  `json:"create_at,omitempty"`
	Remark    string `json:"remark,omitempty" gorm:"size:255"`
}

func (HdlMaterals) TableName() string {
	return "hdl_materals"
}
