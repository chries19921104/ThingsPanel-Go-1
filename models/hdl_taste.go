package models

type HdlTaste struct {
	Id         string `json:"id" gorm:"primaryKey"`           // ID
	Name       string `json:"name,omitempty" gorm:"size:500"` // 口味名称
	TasteId    string `json:"taste_id,omitempty"`             // pos口味id
	CreateAt   int64  `json:"create_at,omitempty"`            // 创建时间
	UpdateTime int64  `json:"update_time,omitempty"`          // 更新时间
	Remark     string `json:"remark,omitempty" gorm:"size:255"`
}

func (HdlTaste) TableName() string {
	return "hdl_taste"
}
