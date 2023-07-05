package services

import (
	"ThingsPanel-Go/initialize/psql"
	"ThingsPanel-Go/models"
	uuid "ThingsPanel-Go/utils"
	valid "ThingsPanel-Go/validate"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"gorm.io/gorm"
)

type HdlTasteService struct {
	//可搜索字段
	SearchField []string
	//可作为条件的字段
	WhereField []string
	//可做为时间范围查询的字段
	TimeField []string
}

func (*HdlTasteService) GetHdlTasteDetail(hdl_taste_id string) []models.HdlTaste {
	var hdl_taste []models.HdlTaste
	psql.Mydb.First(&hdl_taste, "id = ?", hdl_taste_id)
	return hdl_taste
}

// 获取列表
func (*HdlTasteService) GetHdlTasteList(PaginationValidate valid.HdlTastePaginationValidate) (bool, []models.HdlTaste, int64) {
	var HdlTastes []models.HdlTaste
	offset := (PaginationValidate.CurrentPage - 1) * PaginationValidate.PerPage
	db := psql.Mydb.Model(&models.HdlTaste{})
	if PaginationValidate.Name != "" {
		db.Where("name like ?", "%"+PaginationValidate.Name+"%")
	}
	if PaginationValidate.Id != "" {
		db.Where("id = ?", PaginationValidate.Id)
	}
	var count int64
	db.Count(&count)
	result := db.Limit(PaginationValidate.PerPage).Offset(offset).Order("name asc").Find(&HdlTastes)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return false, HdlTastes, 0
	}
	return true, HdlTastes, count
}

// 新增数据
func (*HdlTasteService) AddHdlTaste(hdl_taste valid.AddHdlTasteValidate) (models.HdlTaste, error) {
	var HdlTaste models.HdlTaste = models.HdlTaste{
		Id:         uuid.GetUuid(),
		Name:       hdl_taste.Name,
		TasteId:    hdl_taste.TasteId,
		CreateAt:   time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
		Remark:     hdl_taste.Remark,
	}
	result := psql.Mydb.Create(&HdlTaste)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return HdlTaste, result.Error
	}
	return HdlTaste, nil
}

// 修改数据
func (*HdlTasteService) EditHdlTaste(HdlTaste valid.EditHdlTasteValidate) error {
	// 验证id是否存在
	var HdlTasteModel models.HdlTaste
	result := psql.Mydb.First(&HdlTasteModel, "id = ?", HdlTaste.Id)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		return result.Error
	}
	// 将需要修改的数据组合到结构体中
	HdlTasteModel = models.HdlTaste{
		Name:       HdlTaste.Name,
		TasteId:    HdlTaste.TasteId,
		UpdateTime: time.Now().Unix(),
		Remark:     HdlTaste.Remark,
	}
	result = psql.Mydb.Model(&models.HdlTaste{}).Where("id = ?", HdlTaste.Id).Updates(&HdlTasteModel)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		return result.Error
	}
	return nil
}

// 删除数据
func (*HdlTasteService) DeleteHdlTaste(HdlTaste models.HdlTaste) error {
	result := psql.Mydb.Delete(&HdlTaste)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return result.Error
	}
	return nil
}
