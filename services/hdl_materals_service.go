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

type HdlMateralsService struct {
	//可搜索字段
	SearchField []string
	//可作为条件的字段
	WhereField []string
	//可做为时间范围查询的字段
	TimeField []string
}

func (*HdlMateralsService) GetHdlMateralsDetail(hdl_materals_id string) []models.HdlMaterals {
	var hdl_materals []models.HdlMaterals
	psql.Mydb.First(&hdl_materals, "id = ?", hdl_materals_id)
	return hdl_materals
}

// 获取列表
func (*HdlMateralsService) GetHdlMateralsList(PaginationValidate valid.HdlMateralsPaginationValidate) (bool, []models.HdlMaterals, int64) {
	var HdlMateralss []models.HdlMaterals
	offset := (PaginationValidate.CurrentPage - 1) * PaginationValidate.PerPage
	db := psql.Mydb.Model(&models.HdlMaterals{})
	if PaginationValidate.Name != "" {
		db.Where("name = ?", PaginationValidate.Name)
	}
	if PaginationValidate.Id != "" {
		db.Where("id = ?", PaginationValidate.Id)
	}
	if PaginationValidate.Resource != "" {
		db.Where("resource = ?", PaginationValidate.Resource)
	}
	var count int64
	db.Count(&count)
	result := db.Limit(PaginationValidate.PerPage).Offset(offset).Order("created_at desc").Find(&HdlMateralss)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return false, HdlMateralss, 0
	}
	return true, HdlMateralss, count
}

// 新增数据
func (*HdlMateralsService) AddHdlMaterals(hdl_materals valid.AddHdlMateralsValidate) (models.HdlMaterals, error) {
	var hdlMaterals models.HdlMaterals = models.HdlMaterals{
		Id:        uuid.GetUuid(),
		Name:      hdl_materals.Name,
		Dosage:    hdl_materals.Dosage,
		Unit:      hdl_materals.Unit,
		WaterLine: hdl_materals.WaterLine,
		Station:   hdl_materals.Station,
		Resource:  hdl_materals.Resource,
		CreateAt:  time.Now().Unix(),
		Remark:    hdl_materals.Remark,
	}
	result := psql.Mydb.Create(&hdlMaterals)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return hdlMaterals, result.Error
	}
	return hdlMaterals, nil
}

// 修改数据
func (*HdlMateralsService) EditHdlMaterals(hdlMaterals valid.EditHdlMateralsValidate) bool {
	result := psql.Mydb.Model(&models.HdlMaterals{}).Where("id = ?", hdlMaterals.Id).Updates(&hdlMaterals)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return false
	}
	return true
}

// 删除数据
func (*HdlMateralsService) DeleteHdlMaterals(hdlMaterals models.HdlMaterals) error {
	result := psql.Mydb.Delete(&hdlMaterals)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return result.Error
	}
	return nil
}
