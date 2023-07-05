package models

type HdlRRecipeMaterials struct {
	HdlRecipeId    string `json:"hdl_recipe_id" gorm:"primaryKey"`    // ID
	HdlMaterialsId string `json:"hdl_materials_id" gorm:"primaryKey"` // ID
}

//hdl_r_recipe_materials
func (HdlRRecipeMaterials) TableName() string {
	return "hdl_r_recipe_materials"
}
