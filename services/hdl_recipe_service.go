package services

import (
	"ThingsPanel-Go/initialize/psql"
	"ThingsPanel-Go/models"
	sendmqtt "ThingsPanel-Go/modules/dataService/mqtt/sendMqtt"
	uuid "ThingsPanel-Go/utils"
	valid "ThingsPanel-Go/validate"
	"encoding/json"
	"math/rand"
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

// 获取列表
func (*HdlRecipeService) GetHdlRecipeEntireList(PaginationValidate valid.HdlRecipePaginationValidate, tenantId string) (bool, []map[string]interface{}, int64) {
	var hdlRecipesMap []map[string]interface{}
	offset := (PaginationValidate.CurrentPage - 1) * PaginationValidate.PerPage
	db := psql.Mydb.Model(&models.HdlRecipe{})
	db.Where("tenant_id = ?", tenantId)
	if PaginationValidate.BottomPot != "" {
		db.Where("bottom_pot like ?", "%"+PaginationValidate.BottomPot+"%")
	}
	if PaginationValidate.Id != "" {
		db.Where("id = ?", PaginationValidate.Id)
	}
	var count int64
	db.Count(&count)
	// 通过hdl_pot_type_id连接查询
	db = db.Joins("left join hdl_pot_type on hdl_pot_type.id = hdl_recipe.hdl_pot_type_id")
	// 设置查询字段
	db = db.Select("hdl_recipe.*,hdl_pot_type.name as pot_type_name,hdl_pot_type.pot_type_id as pot_type_id")
	result := db.Limit(PaginationValidate.PerPage).Offset(offset).Order("bottom_pot asc").Find(&hdlRecipesMap)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		return false, hdlRecipesMap, 0
	}
	// 查询配方的物料列表
	for k, v := range hdlRecipesMap {
		var hdlRecipeMaterials []models.HdlMaterials
		result := psql.Mydb.Model(&models.HdlMaterials{}).Joins("left join hdl_r_recipe_materials on hdl_r_recipe_materials.hdl_materials_id = hdl_materials.id").Where("hdl_r_recipe_materials.hdl_recipe_id = ?", v["id"]).Find(&hdlRecipeMaterials)
		if result.Error != nil {
			logs.Error(result.Error, gorm.ErrRecordNotFound)
			return false, hdlRecipesMap, 0
		}
		hdlRecipesMap[k]["materials_list"] = hdlRecipeMaterials
		// 查询配方的口味列表
		var hdlRecipeTasteMap []map[string]interface{}
		result = psql.Mydb.Model(&models.HdlTaste{}).Joins("left join hdl_r_recipe_taste on hdl_r_recipe_taste.hdl_taste_id = hdl_taste.id").Where("hdl_r_recipe_taste.hdl_recipe_id = ?", v["id"]).Find(&hdlRecipeTasteMap)
		if result.Error != nil {
			logs.Error(result.Error, gorm.ErrRecordNotFound)
			return false, hdlRecipesMap, 0
		}
		for k1, v1 := range hdlRecipeTasteMap {
			// 查询口味的物料列表
			var hdlTasteMaterials []models.HdlMaterials
			result := psql.Mydb.Model(&models.HdlMaterials{}).Joins("left join hdl_r_taste_materials on hdl_r_taste_materials.hdl_materials_id = hdl_materials.id").Where("hdl_r_taste_materials.hdl_taste_id = ?", v1["id"]).Find(&hdlTasteMaterials)
			if result.Error != nil {
				logs.Error(result.Error, gorm.ErrRecordNotFound)
				return false, hdlRecipesMap, 0
			}
			hdlRecipeTasteMap[k1]["materials_list"] = hdlTasteMaterials
		}
		hdlRecipesMap[k]["taste_list"] = hdlRecipeTasteMap
	}
	return true, hdlRecipesMap, count
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
func (*HdlRecipeService) AddHdlRecipeAndMateralsAndTaste(hdl_recipe valid.AddEntireHdlRecipeValidate, tenantId string) (models.HdlRecipe, error) {
	// 创建事务
	tx := psql.Mydb.Begin()
	// 创建配方
	var HdlRecipe models.HdlRecipe = models.HdlRecipe{
		Id:               hdl_recipe.Id,
		BottomPotId:      hdl_recipe.BottomPotId,
		BottomPot:        hdl_recipe.BottomPot,
		BottomProperties: hdl_recipe.BottomProperties,
		HdlPotTypeId:     hdl_recipe.HdlPotTypeId,
		UpdateAt:         time.Now().Unix(),
		Remark:           hdl_recipe.Remark,
		TenantId:         tenantId,
	}
	// 判断配方id是否存在
	if HdlRecipe.Id != "" {
		// 存在则更新
		result := tx.Model(&models.HdlRecipe{}).Where("id = ?", HdlRecipe.Id).Updates(&HdlRecipe)
		if result.Error != nil {
			logs.Error(result.Error.Error())
			tx.Rollback()
			return HdlRecipe, result.Error
		}
	} else {
		// 不存在则新增
		HdlRecipe.Id = uuid.GetUuid()
		HdlRecipe.CreateAt = time.Now().Unix()
		result := tx.Create(&HdlRecipe)
		if result.Error != nil {
			logs.Error(result.Error.Error())
			tx.Rollback()
			return HdlRecipe, result.Error
		}
	}
	// 首先删除配方物料关联
	result := tx.Model(&models.HdlRRecipeMaterials{}).Where("hdl_recipe_id = ?", HdlRecipe.Id).Delete(&models.HdlRRecipeMaterials{})
	if result.Error != nil {
		logs.Error(result.Error.Error())
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
				logs.Error(result.Error.Error())
				tx.Rollback()
				return HdlRecipe, result.Error
			}
		} else {
			hdlMaterials.Id = uuid.GetUuid()
			hdlMaterials.CreateAt = time.Now().Unix()
			// 不存在则新增
			result := tx.Create(&hdlMaterials)
			if result.Error != nil {
				logs.Error(result.Error.Error())
				tx.Rollback()
				return HdlRecipe, result.Error
			}
		}
		// 创建配方物料关联
		var hdlRecipeMaterials models.HdlRRecipeMaterials = models.HdlRRecipeMaterials{
			HdlRecipeId:    HdlRecipe.Id,
			HdlMaterialsId: hdlMaterials.Id,
		}
		result = tx.Create(&hdlRecipeMaterials)
		if result.Error != nil {
			logs.Error(result.Error.Error())
			tx.Rollback()
			return HdlRecipe, result.Error
		}

	}
	// 首先删除配方口味关联
	result = tx.Model(&models.HdlRRecipeTaste{}).Where("hdl_recipe_id = ?", HdlRecipe.Id).Delete(&models.HdlRRecipeTaste{})
	if result.Error != nil {
		logs.Error(result.Error.Error())
		tx.Rollback()
		return HdlRecipe, result.Error
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
				logs.Error(result.Error.Error())
				tx.Rollback()
				return HdlRecipe, result.Error
			}
			// 首先删除口味物料关联
			result = tx.Model(&models.HdlRTasteMaterials{}).Where("hdl_taste_id = ?", hdlTaste.Id).Delete(&models.HdlRTasteMaterials{})
			if result.Error != nil {
				logs.Error(result.Error.Error())
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
						logs.Error(result.Error.Error())
						tx.Rollback()
						return HdlRecipe, result.Error
					}
				} else {
					hdlTasteMaterials.Id = uuid.GetUuid()
					hdlTasteMaterials.CreateAt = time.Now().Unix()
					// 不存在则新增
					result := tx.Create(&hdlTasteMaterials)
					if result.Error != nil {
						logs.Error(result.Error.Error())
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
						logs.Error(result.Error.Error())
						tx.Rollback()
						return HdlRecipe, result.Error
					}
				}
				// 创建口味物料关联
				var hdlRTasteMaterials models.HdlRTasteMaterials = models.HdlRTasteMaterials{
					HdlTasteId:     hdlTaste.Id,
					HdlMaterialsId: hdlTasteMaterials.Id,
				}
				result = tx.Create(&hdlRTasteMaterials)
				if result.Error != nil {
					logs.Error(result.Error.Error())
					tx.Rollback()
					return HdlRecipe, result.Error
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
			// 首先删除口味物料关联
			result = tx.Model(&models.HdlRTasteMaterials{}).Where("hdl_taste_id = ?", hdlTaste.Id).Delete(&models.HdlRTasteMaterials{})
			if result.Error != nil {
				logs.Error(result.Error.Error())
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

				}
				// 创建口味物料关联
				var hdlRTasteMaterials models.HdlRTasteMaterials = models.HdlRTasteMaterials{
					HdlTasteId:     hdlTaste.Id,
					HdlMaterialsId: hdlTasteMaterials.Id,
				}
				result = tx.Create(&hdlRTasteMaterials)
				if result.Error != nil {
					logs.Error(result.Error, gorm.ErrRecordNotFound)
					tx.Rollback()
					return HdlRecipe, result.Error
				}
			}
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
func (*HdlRecipeService) DeleteHdlRecipe(HdlRecipe models.HdlRecipe, tenantId string) error {
	// 根据租户id和配方id删除
	result := psql.Mydb.Where("tenant_id = ? and id = ?", tenantId, HdlRecipe.Id).Delete(&models.HdlRecipe{})
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return result.Error
	}
	// 清理未被使用的口味和物料
	var HdlRecipeService HdlRecipeService
	err := HdlRecipeService.ClearHdlTasteAndMaterials()
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	return nil
}

type SendConfig struct {
	PotType   []*SendPotType
	Recipe    []*SendRecipe
	Materials []*SendMaterials
	Taste     []*SendTaste
}
type SendPotType struct {
	Name         string `json:"Name"`
	SoupStandard int    `json:"SoupStandard"`
	PotTypeId    string `json:"PotTypeId"`
}
type SendRecipe struct {
	BottomPotId      string   `json:"BottomPotId"`
	BottomPot        string   `json:"BottomPot"`
	PotTypeId        string   `json:"PotTypeId"`
	MaterialIdList   []string `json:"MaterialIdList"`
	BottomProperties string   `json:"BottomProperties"`
}
type SendMaterials struct {
	Id        string `json:"Id"`
	Name      string `json:"Name"`
	Dosage    int    `json:"Dosage"`
	Unit      string `json:"Unit"`
	WaterLine int    `json:"WaterLine"`
	Station   string `json:"Station"`
	Resource  string `json:"Resource"`
}
type SendTaste struct {
	Name           string   `json:"Name"`
	TasteId        string   `json:"TasteId"`
	PotTypeId      string   `json:"PotTypeId"`
	MaterialIdList []string `json:"MaterialIdList"`
	BottomPotId    string   `json:"BottomPotId"`
}

// 下发配方
func (*HdlRecipeService) SendHdlRecipe(SendHdlRecipeValidate valid.SendHdlRecipeValidate, tenantId string) error {
	// 定义下发配置
	// 定义下发配方
	var sendRecipeList []*SendRecipe
	// 定义下发口味
	var sendTasteList []*SendTaste
	// 根据租户id获取租户下所有配方
	var hdlRecipe []models.HdlRecipe
	result := psql.Mydb.Model(&models.HdlRecipe{}).Where("tenant_id = ?", tenantId).Find(&hdlRecipe)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		return result.Error
	}
	// 遍历hdlRecipe，通过配方id获取配方物料关系
	for i, v := range hdlRecipe {
		// sendRecipeList赋值
		var sendRecipe = &SendRecipe{
			BottomPotId:      v.BottomPotId,
			BottomPot:        v.BottomPot,
			PotTypeId:        v.HdlPotTypeId,
			BottomProperties: v.BottomProperties,
		}
		// 加到sendRecipeList中
		sendRecipeList = append(sendRecipeList, sendRecipe)
		var hdlRRecipeMaterials []models.HdlRRecipeMaterials
		result := psql.Mydb.Model(&models.HdlRRecipeMaterials{}).Where("hdl_recipe_id = ?", v.Id).Find(&hdlRRecipeMaterials)
		if result.Error != nil {
			logs.Error(result.Error.Error())
			return result.Error
		}
		var materialsIdList []string
		// 遍历hdlRRecipeMaterials,获取物料id
		for _, v := range hdlRRecipeMaterials {
			materialsIdList = append(materialsIdList, v.HdlMaterialsId)
		}
		sendRecipeList[i].MaterialIdList = materialsIdList
		// 获取配方关联的口味列表
		var hdlTaste []models.HdlTaste
		result = psql.Mydb.Model(&models.HdlTaste{}).Joins("left join hdl_r_recipe_taste on hdl_r_recipe_taste.hdl_taste_id = hdl_taste.id").Where("hdl_r_recipe_taste.hdl_recipe_id = ?", v.Id).Find(&hdlTaste)
		if result.Error != nil {
			logs.Error(result.Error.Error())
			return result.Error
		}
		// 遍历hdlTaste，给sendTasteList赋值
		for _, vv := range hdlTaste {
			var sendTaste = &SendTaste{
				Name:        vv.Name,
				TasteId:     vv.TasteId,
				PotTypeId:   v.HdlPotTypeId,
				BottomPotId: v.BottomPotId,
			}
			// 加到sendTasteList中
			sendTasteList = append(sendTasteList, sendTaste)
		}

	}

	// 定义锅型
	var SendPotTypeList []*SendPotType
	// 获取所有锅型
	var hdlPotType []models.HdlPotType
	result = psql.Mydb.Model(&models.HdlPotType{}).Find(&hdlPotType)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		return result.Error
	}
	// 遍历hdlPotType，给SendPotTypeList赋值
	for _, v := range hdlPotType {
		var SendPotType = &SendPotType{
			Name:         v.Name,
			SoupStandard: int(v.SoupStandard),
			PotTypeId:    v.Id,
		}
		// 加到SendPotTypeList中
		SendPotTypeList = append(SendPotTypeList, SendPotType)
	}
	// 定义物料
	var SendMaterialsList []*SendMaterials
	// 获取租户下所有配方关联的不重复的物料
	var hdlMaterials []models.HdlMaterials
	result = psql.Mydb.Model(&models.HdlMaterials{}).Joins("left join hdl_r_recipe_materials on hdl_r_recipe_materials.hdl_materials_id = hdl_materials.id").Where("hdl_r_recipe_materials.hdl_recipe_id in (select id from hdl_recipe where tenant_id = ?)", tenantId).Find(&hdlMaterials)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		return result.Error
	}
	// 遍历hdlMaterials，给SendMaterialsList赋值
	for _, v := range hdlMaterials {
		var sendMaterials = &SendMaterials{
			Id:        v.Id,
			Name:      v.Name,
			Dosage:    int(v.Dosage),
			Unit:      v.Unit,
			WaterLine: int(v.WaterLine),
			Station:   v.Station,
			Resource:  v.Resource,
		}
		// 加到SendMaterialsList中
		SendMaterialsList = append(SendMaterialsList, sendMaterials)
	}
	// 获取租户下所有配方关联的口味，这些口味关联下的不重复的物料
	result = psql.Mydb.Model(&models.HdlMaterials{}).Joins("left join hdl_r_taste_materials on hdl_r_taste_materials.hdl_materials_id = hdl_materials.id").Where("hdl_r_taste_materials.hdl_taste_id in (select id from hdl_taste where id in (select hdl_taste_id from hdl_r_recipe_taste where hdl_recipe_id in (select id from hdl_recipe where tenant_id = ?)))", tenantId).Find(&hdlMaterials)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		return result.Error
	}
	// 遍历hdlMaterials，给SendMaterialsList赋值
	for _, v := range hdlMaterials {
		var sendMaterials = &SendMaterials{
			Id:        v.Id,
			Name:      v.Name,
			Dosage:    int(v.Dosage),
			Unit:      v.Unit,
			WaterLine: int(v.WaterLine),
			Station:   v.Station,
			Resource:  v.Resource,
		}
		// 加到SendMaterialsList中
		SendMaterialsList = append(SendMaterialsList, sendMaterials)
	}
	// 定义下发配置
	var sendConfig = &SendConfig{
		PotType:   SendPotTypeList,
		Recipe:    sendRecipeList,
		Materials: SendMaterialsList,
		Taste:     sendTasteList,
	}
	// 根据deviceId获取token
	var hdlDevice models.Device
	result = psql.Mydb.Model(&models.Device{}).Where("id = ?", SendHdlRecipeValidate.DeviceId).Find(&hdlDevice)
	if result.Error != nil {
		logs.Error(result.Error.Error())
		return result.Error
	}
	// 拆分下发配置
	var HdlRecipeService HdlRecipeService
	err := HdlRecipeService.SplitSendMqtt(sendConfig, hdlDevice.Token, 5)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	return nil
}

type PotTypeConfig struct {
	SendId  string
	PotType []*SendPotType
}
type TasteConfig struct {
	SendId string
	Taste  []*SendTaste
}
type MaterialsConfig struct {
	SendId    string
	Materials []*SendMaterials
}
type RecipeConfig struct {
	SendId string
	Recipe []*SendRecipe
}

// 随机字符串
func GetRandomString(l int) string {
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())
	// 生成随机字符串
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, l)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	// 输出结果
	return string(result)
}

// 拆分下发配置
func (*HdlRecipeService) SplitSendMqtt(data *SendConfig, token string, intervalTime int) error {
	//随机6位字符串
	sendId := GetRandomString(6)
	// 锅型配置
	potTypeConfig := &PotTypeConfig{
		SendId:  sendId,
		PotType: data.PotType,
	}
	//发送给mqtt
	bytes, err := json.Marshal(potTypeConfig)
	if err != nil {
		return err
	}
	if err := sendmqtt.SendToHDL(bytes, token); err != nil {
		return err
	}
	//等待时间
	time.Sleep(time.Second * time.Duration(intervalTime))
	//口味配置,每10个口味下发一次
	for i := 0; i < len(data.Taste); i += 10 {
		end := i + 10
		if end > len(data.Taste) {
			end = len(data.Taste)
		}
		tasteConfig := &TasteConfig{
			SendId: sendId,
			Taste:  data.Taste[i:end],
		}
		bytes, err := json.Marshal(tasteConfig)
		if err != nil {
			return err
		}
		if err := sendmqtt.SendToHDL(bytes, token); err != nil {
			return err
		}
		//等待时间
		time.Sleep(time.Second * time.Duration(intervalTime))
	}
	//食材配置,每10个食材下发一次
	for i := 0; i < len(data.Materials); i += 10 {
		end := i + 10
		if end > len(data.Materials) {
			end = len(data.Materials)
		}
		materialsConfig := &MaterialsConfig{
			SendId:    sendId,
			Materials: data.Materials[i:end],
		}
		bytes, err := json.Marshal(materialsConfig)
		if err != nil {
			return err
		}
		if err := sendmqtt.SendToHDL(bytes, token); err != nil {
			return err
		}
		//等待时间
		time.Sleep(time.Second * time.Duration(intervalTime))
	}
	//食谱配置,每10个食谱下发一次
	for i := 0; i < len(data.Recipe); i += 10 {
		end := i + 10
		if end > len(data.Recipe) {
			end = len(data.Recipe)
		}
		recipeConfig := &RecipeConfig{
			SendId: sendId,
			Recipe: data.Recipe[i:end],
		}
		bytes, err := json.Marshal(recipeConfig)
		if err != nil {
			return err
		}
		if err := sendmqtt.SendToHDL(bytes, token); err != nil {
			return err
		}
		//等待时间
		time.Sleep(time.Second * time.Duration(intervalTime))
	}
	return nil

}
