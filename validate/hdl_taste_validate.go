package valid

import "ThingsPanel-Go/models"

type AddHdlTasteValidate struct {
	Name    string `json:"name,omitempty" alias:"口味名称" valid:"Required;MaxSize(500)"`       // 口味名称
	TasteId string `json:"taste_id,omitempty" alias:"pos口味id" valid:"Required;MaxSize(50)"` // pos口味id
	Remark  string `json:"remark,omitempty" alias:"备注" valid:"MaxSize(500)"`
}
type EditHdlTasteValidate struct {
	Id      string `json:"id" alias:"ID" valid:"Required;MaxSize(36)"`             // ID
	Name    string `json:"name,omitempty" alias:"口味名称" valid:"MaxSize(50)"`        // 口味名称
	TasteId string `json:"taste_id,omitempty" alias:"pos口味id" valid:"MaxSize(50)"` // pos口味id
	Remark  string `json:"remark,omitempty" alias:"备注" valid:"MaxSize(500)"`
}
type HdlRecipeTasteValidate struct {
	Id            string                    `json:"id" alias:"ID" valid:"Required;MaxSize(36)"`             // ID
	Name          string                    `json:"name,omitempty" alias:"口味名称" valid:"MaxSize(50)"`        // 口味名称
	TasteId       string                    `json:"taste_id,omitempty" alias:"pos口味id" valid:"MaxSize(50)"` // pos口味id
	Remark        string                    `json:"remark,omitempty" alias:"备注" valid:"MaxSize(500)"`
	MaterialsList []RecipeMaterialsValidate `json:"materials_list,omitempty" alias:"口味物料列表"`
}

type HdlTastePaginationValidate struct {
	CurrentPage int    `json:"current_page"  alias:"当前页" valid:"Required;Min(1)"`
	PerPage     int    `json:"per_page"  alias:"每页页数" valid:"Required;Max(10000)"`
	Name        string `json:"name,omitempty" alias:"口味名称" valid:"MaxSize(500)"`
	Id          string `json:"id,omitempty" alias:"Id" valid:"MaxSize(36)"`
}
type RspHdlTastePaginationValidate struct {
	CurrentPage int               `json:"current_page"  alias:"当前页" valid:"Required;Min(1)"`
	PerPage     int               `json:"per_page"  alias:"每页页数" valid:"Required;Max(10000)"`
	Data        []models.HdlTaste `json:"data" alias:"返回数据"`
	Total       int64             `json:"total" alias:"总数" valid:"Max(10000)"`
}

type HdlTasteIdValidate struct {
	Id string `json:"id"  gorm:"primaryKey" valid:"Required;MaxSize(36)"`
}
