package models

type HdlRRecipeTaste struct {
	HdlRecipeId string `json:"hdl_recipe_id" gorm:"primaryKey"` // ID
	HdlTasteId  string `json:"hdl_taste_id" gorm:"primaryKey"`  // ID
}

//hdl_r_recipe_taste
func (HdlRRecipeTaste) TableName() string {
	return "hdl_r_recipe_taste"
}
