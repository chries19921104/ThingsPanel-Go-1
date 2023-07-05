package models

type HdlRTasteMaterials struct {
	HdlMaterialsId string `json:"hdl_materials_id" gorm:"primaryKey"` // ID
	HdlTasteId     string `json:"hdl_taste_id" gorm:"primaryKey"`     // ID
}

//hdl_r_taste_materials
func (HdlRTasteMaterials) TableName() string {
	return "hdl_r_taste_materials"
}
