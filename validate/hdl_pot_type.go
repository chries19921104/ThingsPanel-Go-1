package valid

import "ThingsPanel-Go/models"

type AddHdlPotTypeValidate struct {
	Name         string `json:"name,omitempty" alias:"锅型名称" valid:"Required;MaxSize(500)"` // 锅型名称
	Image        string `json:"image,omitempty" alias:"图片" valid:"Required;MaxSize(500)"`  // 图片
	SoupStandard int64  `json:"soup_standard,omitempty" alias:"水位线标准" valid:"Required"`    //水位线标准
	Remark       string `json:"remark,omitempty" alias:"备注" valid:"MaxSize(500)"`
	PotTypeId    string `json:"pot_type_id,omitempty" alias:"锅型id" valid:"MaxSize(50)"` // 锅型id
}
type EditHdlPotTypeValidate struct {
	Id           string `json:"id" alias:"ID" valid:"Required;MaxSize(36)"`                // ID
	Name         string `json:"name,omitempty" alias:"锅型名称" valid:"Required;MaxSize(500)"` // 锅型名称
	Image        string `json:"image,omitempty" alias:"图片" valid:"Required;MaxSize(500)"`  // 图片
	SoupStandard int64  `json:"soup_standard,omitempty" alias:"水位线标准" valid:"Required"`    //水位线标准
	Remark       string `json:"remark" alias:"备注" valid:"MaxSize(500)"`
	PotTypeId    string `json:"pot_type_id,omitempty" alias:"锅型id" valid:"MaxSize(50)"` // 锅型id
}

type HdlPotTypePaginationValidate struct {
	CurrentPage int    `json:"current_page"  alias:"当前页" valid:"Required;Min(1)"`
	PerPage     int    `json:"per_page"  alias:"每页页数" valid:"Required;Max(10000)"`
	Name        string `json:"name,omitempty" alias:"锅型名称" valid:"Required;MaxSize(500)"`
	Id          string `json:"id,omitempty" alias:"Id" valid:"MaxSize(36)"`
}
type RspHdlPotTypePaginationValidate struct {
	CurrentPage int                 `json:"current_page"  alias:"当前页" valid:"Required;Min(1)"`
	PerPage     int                 `json:"per_page"  alias:"每页页数" valid:"Required;Max(10000)"`
	Data        []models.HdlPotType `json:"data" alias:"返回数据"`
	Total       int64               `json:"total" alias:"总数" valid:"Max(10000)"`
}

type HdlPotTypeIdValidate struct {
	Id string `json:"id"  gorm:"primaryKey" valid:"Required;MaxSize(36)"`
}
