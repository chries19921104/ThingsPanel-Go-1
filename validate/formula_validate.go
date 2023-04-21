package valid

import (
	"ThingsPanel-Go/models"
)

type AddRecipeValidator struct {
	Id          string `json:"id"`
	BottomPotId string `json:"BottomPotId" alias:"锅底ID" valid:"Required"`
	BottomPot   string `json:"BottomPot" alias:"锅底" valid:"Required"`
	PotTypeId   string `json:"PotTypeId" alias:"锅型ID" valid:"Required"`
	//PotTypeName      string     `json:"PotTypeName" alias:"锅型名称" valid:"Required"`
	MaterialsId      string     `json:"MaterialsId" alias:"物料ID"`
	TasteId          string     `json:"TasteId" alias:"口味ID"`
	Tastes           string     `json:"Taste" alias:"口味"`
	TastesArr        []Taste    `json:"TastesArr" alias:"口味" valid:"Required"`
	BottomProperties string     `json:"BottomProperties" alias:"锅底属性" valid:"Required"`
	SoupStandard     int64      `json:"SoupStandard" alias:"加汤水位标准" `
	MaterialsArr     []Material `json:"MaterialsArr" alias:"物料"`
	Materials        string     `json:"Materials" alias:"物料"`
	CurrentWaterLine int64      `json:"CurrentWaterLine" alias:"当前加汤水位线"`
}

type Taste struct {
	Name          string `json:"name"`
	TasteId       string `json:"taste_id"`
	MaterialsName string `json:"materials_name"`
	Dosage        int    `json:"dosage"`
	Unit          string `json:"unit"`
	WaterLine     int    `json:"water_line"`
	Station       string `json:"station"`
}

type Material struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Dosage    int    `json:"dosage"`
	Unit      string `json:"unit"`
	WaterLine int    `json:"water_line"`
	Station   string `json:"station"`
}

type RecipePaginationValidate struct {
	CurrentPage int    `json:"current_page"  alias:"当前页" valid:"Required;Min(1)"`
	PerPage     int    `json:"per_page"  alias:"每页页数" valid:"Required;Max(10000)"`
	Id          string `json:"id" alias:"配方ID"`
}

type RspRecipePaginationValidate struct {
	CurrentPage int                  `json:"current_page"  alias:"当前页" valid:"Required;Min(1)"`
	PerPage     int                  `json:"per_page"  alias:"每页页数" valid:"Required;Max(10000)"`
	Data        []models.RecipeValue `json:"data" alias:"返回数据"`
	Total       int64                `json:"total" alias:"总数" valid:"Max(10000)"`
}

type DelRecipeValidator struct {
	Id string `json:"id" valid:"Required"`
}


type SearchMaterialNameValidator struct {
	Keyword string `json:"keyword"`
}