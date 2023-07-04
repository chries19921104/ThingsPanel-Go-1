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

type HdlMaterialsService struct {
	//可搜索字段
	SearchField []string
	//可作为条件的字段
	WhereField []string
	//可做为时间范围查询的字段
	TimeField []string
}

func (*HdlMaterialsService) GetHdlMaterialsDetail(hdl_materials_id string) []models.HdlMaterials {
	var hdl_materials []models.HdlMaterials
	psql.Mydb.First(&hdl_materials, "id = ?", hdl_materials_id)
	return hdl_materials
}

// 获取列表
func (*HdlMaterialsService) GetHdlMaterialsList(PaginationValidate valid.HdlMaterialsPaginationValidate) (bool, []models.HdlMaterials, int64) {
	var HdlMaterialss []models.HdlMaterials
	offset := (PaginationValidate.CurrentPage - 1) * PaginationValidate.PerPage
	db := psql.Mydb.Model(&models.HdlMaterials{})
	if PaginationValidate.Name != "" {
		db.Where("name like ?", "%"+PaginationValidate.Name+"%")
	}
	if PaginationValidate.Id != "" {
		db.Where("id = ?", PaginationValidate.Id)
	}
	if PaginationValidate.Resource != "" {
		db.Where("resource = ?", PaginationValidate.Resource)
	}
	var count int64
	db.Count(&count)
	result := db.Limit(PaginationValidate.PerPage).Offset(offset).Order("name asc").Find(&HdlMaterialss)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return false, HdlMaterialss, 0
	}
	return true, HdlMaterialss, count
}

// 新增数据
func (*HdlMaterialsService) AddHdlMaterials(hdl_materials valid.AddHdlMaterialsValidate) (models.HdlMaterials, error) {
	var hdlMaterials models.HdlMaterials = models.HdlMaterials{
		Id:        uuid.GetUuid(),
		Name:      hdl_materials.Name,
		Dosage:    hdl_materials.Dosage,
		Unit:      hdl_materials.Unit,
		WaterLine: hdl_materials.WaterLine,
		Station:   hdl_materials.Station,
		Resource:  hdl_materials.Resource,
		CreateAt:  time.Now().Unix(),
		Remark:    hdl_materials.Remark,
	}
	result := psql.Mydb.Create(&hdlMaterials)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return hdlMaterials, result.Error
	}
	return hdlMaterials, nil
}

// 修改数据
func (*HdlMaterialsService) EditHdlMaterials(hdlMaterials valid.EditHdlMaterialsValidate) error {
	// 验证id是否存在
	var hdl_materials models.HdlMaterials
	result := psql.Mydb.First(&hdl_materials, "id = ?", hdlMaterials.Id)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		return result.Error
	}
	// 将需要修改的数据组合到结构体中
	hdl_materials = models.HdlMaterials{
		Name:      hdlMaterials.Name,
		Dosage:    hdlMaterials.Dosage,
		Unit:      hdlMaterials.Unit,
		WaterLine: hdlMaterials.WaterLine,
		Station:   hdlMaterials.Station,
		Resource:  hdlMaterials.Resource,
		Remark:    hdlMaterials.Remark,
	}
	result = psql.Mydb.Model(&models.HdlMaterials{}).Where("id = ?", hdlMaterials.Id).Updates(&hdlMaterials)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		return result.Error
	}
	return nil
}

// 删除数据
func (*HdlMaterialsService) DeleteHdlMaterials(hdlMaterials models.HdlMaterials) error {
	result := psql.Mydb.Delete(&hdlMaterials)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return result.Error
	}
	return nil
}
