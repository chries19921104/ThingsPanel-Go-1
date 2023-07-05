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

type HdlPotTypeService struct {
	//可搜索字段
	SearchField []string
	//可作为条件的字段
	WhereField []string
	//可做为时间范围查询的字段
	TimeField []string
}

func (*HdlPotTypeService) GetHdlPotTypeDetail(pot_id string) []models.HdlPotType {
	var pot []models.HdlPotType
	psql.Mydb.First(&pot, "id = ?", pot_id)
	return pot
}

// 获取列表
func (*HdlPotTypeService) GetHdlPotTypeList(PaginationValidate valid.HdlPotTypePaginationValidate) (bool, []models.HdlPotType, int64) {
	var HdlPotTypes []models.HdlPotType
	offset := (PaginationValidate.CurrentPage - 1) * PaginationValidate.PerPage
	db := psql.Mydb.Model(&models.HdlPotType{})
	if PaginationValidate.Name != "" {
		db.Where("name like ?", "%"+PaginationValidate.Name+"%")
	}
	if PaginationValidate.Id != "" {
		db.Where("id = ?", PaginationValidate.Id)
	}
	var count int64
	db.Count(&count)
	result := db.Limit(PaginationValidate.PerPage).Offset(offset).Order("create_at desc").Find(&HdlPotTypes)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return false, HdlPotTypes, 0
	}
	return true, HdlPotTypes, count
}

// 新增数据
func (*HdlPotTypeService) AddHdlPotType(pot valid.AddHdlPotTypeValidate) (models.HdlPotType, error) {
	var HdlPotType models.HdlPotType = models.HdlPotType{
		Id:           uuid.GetUuid(),
		Name:         pot.Name,
		Image:        pot.Image,
		SoupStandard: pot.SoupStandard,
		CreateAt:     time.Now().Unix(),
		UpdateAt:     time.Now().Unix(),
		Remark:       pot.Remark,
		PotTypeId:    pot.PotTypeId,
	}
	result := psql.Mydb.Create(&HdlPotType)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return HdlPotType, result.Error
	}
	return HdlPotType, nil
}

// 修改数据
func (*HdlPotTypeService) EditHdlPotType(HdlPotType valid.EditHdlPotTypeValidate) error {
	// 验证id是否存在
	var HdlPotTypeModel models.HdlPotType
	result := psql.Mydb.First(&HdlPotTypeModel, "id = ?", HdlPotType.Id)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		return result.Error
	}
	// 将需要修改的数据组合到结构体中
	HdlPotTypeModel = models.HdlPotType{
		Name:         HdlPotType.Name,
		Image:        HdlPotType.Image,
		UpdateAt:     time.Now().Unix(),
		SoupStandard: HdlPotType.SoupStandard,
		Remark:       HdlPotType.Remark,
		PotTypeId:    HdlPotType.PotTypeId,
	}
	result = psql.Mydb.Model(&models.HdlPotType{}).Where("id = ?", HdlPotType.Id).Updates(&HdlPotTypeModel)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		return result.Error
	}
	return nil
}

// 删除数据
func (*HdlPotTypeService) DeleteHdlPotType(HdlPotType models.HdlPotType) error {
	result := psql.Mydb.Delete(&HdlPotType)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return result.Error
	}
	return nil
}
