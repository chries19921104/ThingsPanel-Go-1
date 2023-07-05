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

type HdlRecipeService struct {
	//可搜索字段
	SearchField []string
	//可作为条件的字段
	WhereField []string
	//可做为时间范围查询的字段
	TimeField []string
}

func (*HdlRecipeService) GetHdlRecipeDetail(hdl_recipe_id string) []models.HdlRecipe {
	var hdl_recipe []models.HdlRecipe
	psql.Mydb.First(&hdl_recipe, "id = ?", hdl_recipe_id)
	return hdl_recipe
}

// 获取列表
func (*HdlRecipeService) GetHdlRecipeList(PaginationValidate valid.HdlRecipePaginationValidate) (bool, []models.HdlRecipe, int64) {
	var HdlRecipes []models.HdlRecipe
	offset := (PaginationValidate.CurrentPage - 1) * PaginationValidate.PerPage
	db := psql.Mydb.Model(&models.HdlRecipe{})
	if PaginationValidate.BottomPot != "" {
		db.Where("name like ?", "%"+PaginationValidate.BottomPot+"%")
	}
	if PaginationValidate.Id != "" {
		db.Where("id = ?", PaginationValidate.Id)
	}
	var count int64
	db.Count(&count)
	result := db.Limit(PaginationValidate.PerPage).Offset(offset).Order("bottom_pot asc").Find(&HdlRecipes)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return false, HdlRecipes, 0
	}
	return true, HdlRecipes, count
}

// 新增数据
func (*HdlRecipeService) AddHdlRecipe(hdl_recipe valid.AddHdlRecipeValidate) (models.HdlRecipe, error) {
	var HdlRecipe models.HdlRecipe = models.HdlRecipe{
		Id:               uuid.GetUuid(),
		BottomPotId:      hdl_recipe.BottomPotId,
		BottomPot:        hdl_recipe.BottomPot,
		BottomProperties: hdl_recipe.BottomProperties,
		HdlPotTypeId:     hdl_recipe.HdlPotTypeId,
		CreateAt:         time.Now().Unix(),
		UpdateAt:         time.Now().Unix(),
		Remark:           hdl_recipe.Remark,
	}
	result := psql.Mydb.Create(&HdlRecipe)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return HdlRecipe, result.Error
	}
	return HdlRecipe, nil
}

// 新增配方整体数据，包括物料和口味
func (*HdlRecipeService) AddHdlRecipeAndMateralsAndTaste(hdl_recipe valid.AddEntireHdlRecipeValidate) (models.HdlRecipe, error) {
	// 创建事务
	tx := psql.Mydb.Begin()
	// 创建配方
	var HdlRecipe models.HdlRecipe = models.HdlRecipe{
		Id:               uuid.GetUuid(),
		BottomPotId:      hdl_recipe.BottomPotId,
		BottomPot:        hdl_recipe.BottomPot,
		BottomProperties: hdl_recipe.BottomProperties,
		HdlPotTypeId:     hdl_recipe.HdlPotTypeId,
		CreateAt:         time.Now().Unix(),
		UpdateAt:         time.Now().Unix(),
		Remark:           hdl_recipe.Remark,
	}
	result := tx.Create(&HdlRecipe)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		tx.Rollback()
		return HdlRecipe, result.Error
	}
	// 创建物料
	for _, v := range hdl_recipe.MaterialsList {
		var hdlMaterials models.HdlMaterials = models.HdlMaterials{
			Name:      v.Name,
			Unit:      v.Unit,
			Dosage:    v.Dosage,
			WaterLine: v.WaterLine,
			Station:   v.Station,
			Resource:  v.Resource,
			Remark:    v.Remark,
		}
		// 判断物料id是否存在
		if hdlMaterials.Id != "" {
			// 存在则更新
			result := tx.Model(&models.HdlMaterials{}).Where("id = ?", hdlMaterials.Id).Updates(&hdlMaterials)
			if result.Error != nil {
				logs.Error(result.Error, gorm.ErrRecordNotFound)
				tx.Rollback()
				return HdlRecipe, result.Error
			}
		} else {
			hdlMaterials.Id = uuid.GetUuid()
			hdlMaterials.CreateAt = time.Now().Unix()
			// 不存在则新增
			result := tx.Create(&hdlMaterials)
			if result.Error != nil {
				logs.Error(result.Error, gorm.ErrRecordNotFound)
				tx.Rollback()
				return HdlRecipe, result.Error
			}
			// 创建配方物料关联
			var hdlRecipeMaterials models.HdlRRecipeMaterials = models.HdlRRecipeMaterials{
				HdlRecipeId:    HdlRecipe.Id,
				HdlMaterialsId: hdlMaterials.Id,
			}
			result = tx.Create(&hdlRecipeMaterials)
			if result.Error != nil {
				logs.Error(result.Error, gorm.ErrRecordNotFound)
				tx.Rollback()
				return HdlRecipe, result.Error
			}
		}

	}
	// 创建口味
	for _, v := range hdl_recipe.TasteList {
		var hdlTaste models.HdlTaste = models.HdlTaste{
			Name:       v.Name,
			TasteId:    v.TasteId,
			UpdateTime: time.Now().Unix(),
			Remark:     v.Remark,
		}
		if hdlTaste.Id != "" {
			// 存在则更新
			result := tx.Model(&models.HdlTaste{}).Where("id = ?", hdlTaste.Id).Updates(&hdlTaste)
			if result.Error != nil {
				logs.Error(result.Error, gorm.ErrRecordNotFound)
				tx.Rollback()
				return HdlRecipe, result.Error
			}
			// 创建口味物料
			for _, v := range v.MaterialsList {
				var hdlTasteMaterials models.HdlMaterials = models.HdlMaterials{
					Name:      v.Name,
					Unit:      v.Unit,
					Dosage:    v.Dosage,
					WaterLine: v.WaterLine,
					Station:   v.Station,
					Resource:  v.Resource,
					Remark:    v.Remark,
				}
				// 判断物料id是否存在
				if hdlTasteMaterials.Id != "" {
					// 存在则更新
					result := tx.Model(&models.HdlMaterials{}).Where("id = ?", hdlTasteMaterials.Id).Updates(&hdlTasteMaterials)
					if result.Error != nil {
						logs.Error(result.Error, gorm.ErrRecordNotFound)
						tx.Rollback()
						return HdlRecipe, result.Error
					}
				} else {
					hdlTasteMaterials.Id = uuid.GetUuid()
					hdlTasteMaterials.CreateAt = time.Now().Unix()
					// 不存在则新增
					result := tx.Create(&hdlTasteMaterials)
					if result.Error != nil {
						logs.Error(result.Error, gorm.ErrRecordNotFound)
						tx.Rollback()
						return HdlRecipe, result.Error
					}
					// 创建口味物料关联
					var hdlTasteMaterials models.HdlRTasteMaterials = models.HdlRTasteMaterials{
						HdlTasteId:     hdlTaste.Id,
						HdlMaterialsId: hdlTasteMaterials.Id,
					}
					result = tx.Create(&hdlTasteMaterials)
					if result.Error != nil {
						logs.Error(result.Error, gorm.ErrRecordNotFound)
						tx.Rollback()
						return HdlRecipe, result.Error
					}
				}
			}
		} else {
			// 不存在则新增
			hdlTaste.Id = uuid.GetUuid()
			hdlTaste.CreateAt = time.Now().Unix()
			result := tx.Create(&hdlTaste)
			if result.Error != nil {
				logs.Error(result.Error, gorm.ErrRecordNotFound)
				tx.Rollback()
				return HdlRecipe, result.Error
			}
			// 创建配方口味关联
			var hdlRecipeTaste models.HdlRRecipeTaste = models.HdlRRecipeTaste{
				HdlRecipeId: HdlRecipe.Id,
				HdlTasteId:  hdlTaste.Id,
			}
			// 创建口味物料
			result = tx.Create(&hdlRecipeTaste)
			if result.Error != nil {
				logs.Error(result.Error, gorm.ErrRecordNotFound)
				tx.Rollback()
				return HdlRecipe, result.Error
			}
			// 创建口味物料
			for _, v := range v.MaterialsList {
				var hdlTasteMaterials models.HdlMaterials = models.HdlMaterials{
					Name:      v.Name,
					Unit:      v.Unit,
					Dosage:    v.Dosage,
					WaterLine: v.WaterLine,
					Station:   v.Station,
					Resource:  v.Resource,
					Remark:    v.Remark,
				}
				// 判断物料id是否存在
				if hdlTasteMaterials.Id != "" {
					// 存在则更新
					result := tx.Model(&models.HdlMaterials{}).Where("id = ?", hdlTasteMaterials.Id).Updates(&hdlTasteMaterials)
					if result.Error != nil {
						logs.Error(result.Error, gorm.ErrRecordNotFound)
						tx.Rollback()
						return HdlRecipe, result.Error
					}
				} else {
					hdlTasteMaterials.Id = uuid.GetUuid()
					hdlTasteMaterials.CreateAt = time.Now().Unix()
					// 不存在则新增
					result := tx.Create(&hdlTasteMaterials)
					if result.Error != nil {
						logs.Error(result.Error, gorm.ErrRecordNotFound)
						tx.Rollback()
						return HdlRecipe, result.Error
					}
					// 创建口味物料关联
					var hdlTasteMaterials models.HdlRTasteMaterials = models.HdlRTasteMaterials{
						HdlTasteId:     hdlTaste.Id,
						HdlMaterialsId: hdlTasteMaterials.Id,
					}
					result = tx.Create(&hdlTasteMaterials)
					if result.Error != nil {
						logs.Error(result.Error, gorm.ErrRecordNotFound)
						tx.Rollback()
						return HdlRecipe, result.Error
					}
				}
			}
		}
	}
	// 提交事务
	tx.Commit()
	// 清理未被使用的口味和物料
	var HdlRecipeService HdlRecipeService
	err := HdlRecipeService.ClearHdlTasteAndMaterials()
	if err != nil {
		logs.Error(err.Error())
		return HdlRecipe, err
	}
	// 返回数据
	return HdlRecipe, nil
}

// 清理未被使用的口味和物料
func (*HdlRecipeService) ClearHdlTasteAndMaterials() error {
	// 创建事务
	tx := psql.Mydb.Begin()
	// 清理未被使用的口味，即清理配方口味关系表中不存在的口味
	var hdlTaste []models.HdlTaste
	result := tx.Model(&models.HdlTaste{}).Where("id not in (select hdl_taste_id from hdl_r_recipe_taste)").Find(&hdlTaste)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		tx.Rollback()
		return result.Error
	}
	for _, v := range hdlTaste {
		// 清理口味物料关系
		result := tx.Model(&models.HdlRTasteMaterials{}).Where("hdl_taste_id = ?", v.Id).Delete(&models.HdlRTasteMaterials{})
		if result.Error != nil {
			logs.Error(result.Error.Error())
			tx.Rollback()
			return result.Error
		}
		// 清理口味
		result = tx.Model(&models.HdlTaste{}).Where("id = ?", v.Id).Delete(&models.HdlTaste{})
		if result.Error != nil {
			logs.Error(result.Error.Error())
			tx.Rollback()
			return result.Error
		}
	}
	// 清理未被使用的物料，即清理配方物料关系表中不存在且物料的Resource值是material的物料
	var hdlMaterials []models.HdlMaterials
	result = tx.Model(&models.HdlMaterials{}).Where("id not in (select hdl_materials_id from hdl_r_recipe_materials) and resource = 'material'").Find(&hdlMaterials)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		tx.Rollback()
		return result.Error
	}
	for _, v := range hdlMaterials {
		// 清理物料
		result := tx.Model(&models.HdlMaterials{}).Where("id = ?", v.Id).Delete(&models.HdlMaterials{})
		if result.Error != nil {
			logs.Error(result.Error.Error())
			tx.Rollback()
			return result.Error
		}
	}
	// 清理未被使用的物料，即清理口味物料关系表中不存在且物料的Resource值是taste的物料
	result = tx.Model(&models.HdlMaterials{}).Where("id not in (select hdl_materials_id from hdl_r_taste_materials) and resource = 'taste'").Find(&hdlMaterials)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		tx.Rollback()
		return result.Error
	}
	for _, v := range hdlMaterials {
		// 清理物料
		result := tx.Model(&models.HdlMaterials{}).Where("id = ?", v.Id).Delete(&models.HdlMaterials{})
		if result.Error != nil {
			logs.Error(result.Error.Error())
			tx.Rollback()
			return result.Error
		}
	}
	// 提交事务
	tx.Commit()
	return nil
}

// 修改数据
func (*HdlRecipeService) EditHdlRecipe(HdlRecipe valid.EditHdlRecipeValidate) error {
	// 验证id是否存在
	var HdlRecipeModel models.HdlRecipe
	result := psql.Mydb.First(&HdlRecipeModel, "id = ?", HdlRecipe.Id)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		return result.Error
	}
	// 将需要修改的数据组合到结构体中
	HdlRecipeModel = models.HdlRecipe{
		BottomPotId:      HdlRecipe.BottomPotId,
		BottomPot:        HdlRecipe.BottomPot,
		BottomProperties: HdlRecipe.BottomProperties,
		HdlPotTypeId:     HdlRecipe.HdlPotTypeId,
		UpdateAt:         time.Now().Unix(),
		Remark:           HdlRecipe.Remark,
	}
	result = psql.Mydb.Model(&models.HdlRecipe{}).Where("id = ?", HdlRecipe.Id).Updates(&HdlRecipeModel)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		return result.Error
	}
	return nil
}

// 删除数据
func (*HdlRecipeService) DeleteHdlRecipe(HdlRecipe models.HdlRecipe) error {
	result := psql.Mydb.Delete(&HdlRecipe)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return result.Error
	}
	return nil
}
