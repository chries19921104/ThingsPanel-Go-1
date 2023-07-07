package valid

import "ThingsPanel-Go/models"

type AddHdlRecipeValidate struct {
	BottomPotId      string `json:"bottom_pot_id,omitempty" alias:"hdl锅底id" valid:"MaxSize(500)"` // 锅底id
	BottomPot        string `json:"bottom_pot,omitempty" alias:"锅底名称" valid:"MaxSize(500)"`       // 锅底名称
	BottomProperties string `json:"bottom_properties,omitempty" alias:"锅底属性" valid:"MaxSize(50)"` // 锅底属性 辣 不辣
	HdlPotTypeId     string `json:"hdl_pot_type_id,omitempty" alias:"锅型id" valid:"MaxSize(50)"`   // 锅型id
	Remark           string `json:"remark,omitempty" alias:"备注" valid:"MaxSize(500)"`
}
type AddEntireHdlRecipeValidate struct {
	Id               string                    `json:"id" alias:"ID" valid:"MaxSize(36)"`                            // ID
	BottomPotId      string                    `json:"bottom_pot_id,omitempty" alias:"hdl锅底id" valid:"MaxSize(500)"` // 锅底id
	BottomPot        string                    `json:"bottom_pot,omitempty" alias:"锅底名称" valid:"MaxSize(500)"`       // 锅底名称
	BottomProperties string                    `json:"bottom_properties,omitempty" alias:"锅底属性" valid:"MaxSize(50)"` // 锅底属性 辣 不辣
	HdlPotTypeId     string                    `json:"hdl_pot_type_id,omitempty" alias:"锅型id" valid:"MaxSize(50)"`   // 锅型id
	Remark           string                    `json:"remark,omitempty" alias:"备注" valid:"MaxSize(500)"`
	MaterialsList    []RecipeMaterialsValidate `json:"materials_list,omitempty" alias:"物料列表"` // 物料列表
	TasteList        []HdlRecipeTasteValidate  `json:"taste_list,omitempty" alias:"口味列表"`     // 口味列表
}
type EditHdlRecipeValidate struct {
	Id               string `json:"id" alias:"ID" valid:"Required;MaxSize(36)"`                   // ID
	BottomPotId      string `json:"bottom_pot_id,omitempty" alias:"hdl锅底id" valid:"MaxSize(500)"` // 锅底id
	BottomPot        string `json:"bottom_pot,omitempty" alias:"锅底名称" valid:"MaxSize(500)"`       // 锅底名称
	BottomProperties string `json:"bottom_properties,omitempty" alias:"锅底属性" valid:"MaxSize(50)"` // 锅底属性 辣 不辣
	HdlPotTypeId     string `json:"hdl_pot_type_id,omitempty" alias:"锅型id" valid:"MaxSize(50)"`   // 锅型id
	Remark           string `json:"remark,omitempty" alias:"备注" valid:"MaxSize(500)"`
}

type HdlRecipePaginationValidate struct {
	CurrentPage int    `json:"current_page"  alias:"当前页" valid:"Required;Min(1)"`
	PerPage     int    `json:"per_page"  alias:"每页页数" valid:"Required;Max(10000)"`
	BottomPot   string `json:"bottom_pot,omitempty" alias:"口味名称" valid:"MaxSize(500)"`
	Id          string `json:"id,omitempty" alias:"Id" valid:"MaxSize(36)"`
	OrderBy     string `json:"order_by,omitempty" alias:"排序" valid:"MaxSize(50)"`
	SortBy      string `json:"sort_by,omitempty" alias:"排序" valid:"MaxSize(50)"`
}
type RspHdlRecipePaginationValidate struct {
	CurrentPage int                `json:"current_page"  alias:"当前页" valid:"Required;Min(1)"`
	PerPage     int                `json:"per_page"  alias:"每页页数" valid:"Required;Max(10000)"`
	Data        []models.HdlRecipe `json:"data" alias:"返回数据"`
	Total       int64              `json:"total" alias:"总数" valid:"Max(10000)"`
}
type RspHdlEntireRecipePaginationValidate struct {
	CurrentPage int                      `json:"current_page"  alias:"当前页" valid:"Required;Min(1)"`
	PerPage     int                      `json:"per_page"  alias:"每页页数" valid:"Required;Max(10000)"`
	Data        []map[string]interface{} `json:"data" alias:"返回数据"`
	Total       int64                    `json:"total" alias:"总数" valid:"Max(10000)"`
}
type HdlRecipeIdValidate struct {
	Id string `json:"id"  gorm:"primaryKey" valid:"Required;MaxSize(36)"`
}

// 下发配方
type SendHdlRecipeValidate struct {
	DeviceId string `json:"device_id"  valid:"Required;MaxSize(36)"`
}
