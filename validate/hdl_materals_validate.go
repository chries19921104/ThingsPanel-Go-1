package valid

import "ThingsPanel-Go/models"

type AddHdlMateralsValidate struct {
	Name      string `json:"name,omitempty" alias:"物料名称" valid:"MaxSize(500)"`             // 物料名称
	Dosage    int64  `json:"dosage,omitempty" alias:"用量"`                                  // 用量
	Unit      string `json:"unit,omitempty" alias:"单位" valid:"MaxSize(50)"`                // 单位
	WaterLine int64  `json:"water_line,omitempty" alias:"加汤水位标准"`                          // 加汤水位标准
	Station   string `json:"station,omitempty" alias:"工位" valid:"Required;MaxSize(50)"`    //工位 （鲜料工位、传锅工位、所有工位、不显示）
	Resource  string `json:"resource,omitempty" alias:"物料来源" valid:"Required;MaxSize(50)"` //物料来源 material-锅底 taste-口味
	CreateAt  int64  `json:"create_at,omitempty"  alias:"创建时间" `
	Remark    string `json:"remark,omitempty" alias:"备注" valid:"Required;MaxSize(500)"`
}
type EditHdlMateralsValidate struct {
	Id string `json:"id" alias:"ID" valid:"Required;MaxSize(36)"` // ID
	AddHdlMateralsValidate
}

type HdlMateralsPaginationValidate struct {
	CurrentPage int    `json:"current_page"  alias:"当前页" valid:"Required;Min(1)"`
	PerPage     int    `json:"per_page"  alias:"每页页数" valid:"Required;Max(10000)"`
	Name        string `json:"protocol_type,omitempty" alias:"物料名称" valid:"MaxSize(500)"`
	Resource    string `json:"device_type,omitempty" alias:"物料来源" valid:"MaxSize(50)"`
	Id          string `json:"id,omitempty" alias:"Id" valid:"MaxSize(36)"`
}
type RspHdlMateralsPaginationValidate struct {
	CurrentPage int                  `json:"current_page"  alias:"当前页" valid:"Required;Min(1)"`
	PerPage     int                  `json:"per_page"  alias:"每页页数" valid:"Required;Max(10000)"`
	Data        []models.HdlMaterals `json:"data" alias:"返回数据"`
	Total       int64                `json:"total" alias:"总数" valid:"Max(10000)"`
}

type HdlMateralsIdValidate struct {
	Id string `json:"id"  gorm:"primaryKey" valid:"Required;MaxSize(36)"`
}
