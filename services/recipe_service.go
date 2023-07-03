package services

import (
	"ThingsPanel-Go/initialize/psql"
	"ThingsPanel-Go/models"
	"ThingsPanel-Go/modules/dataService/mqtt"
	"ThingsPanel-Go/utils"
	valid "ThingsPanel-Go/validate"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"gorm.io/gorm"
)

type RecipeService struct {
	//可搜索字段
	SearchField []string
	//可作为条件的字段
	WhereField []string
	//可做为时间范围查询的字段
	TimeField []string
}

func (*RecipeService) GetRecipeDetail(recipeId string) []models.Recipe {
	var recipe []models.Recipe
	psql.Mydb.First(&recipe, "recipe.id = ?", recipeId)
	return recipe
}

// 获取列表
func (*RecipeService) GetRecipeList(PaginationValidate valid.RecipePaginationValidate) (bool, []models.RecipeValue, int64) {
	var Recipe []models.RecipeValue
	offset := (PaginationValidate.CurrentPage - 1) * PaginationValidate.PerPage
	db := psql.Mydb.Model(&models.Recipe{})
	if PaginationValidate.Id != "" {
		db = db.Where("recipe.id = ?", PaginationValidate.Id)
	}
	db = db.Select("recipe.id,recipe.bottom_pot_id,recipe.bottom_pot,recipe.pot_type_id,recipe.taste_materials,recipe.bottom_properties,recipe.soup_standard,pot_type.name").Joins("left join pot_type on recipe.pot_type_id = pot_type.pot_type_id").Where("recipe.is_del", false)

	var count int64
	db.Count(&count)
	result := db.Limit(PaginationValidate.PerPage).Offset(offset).Order("recipe.create_at desc").Find(&Recipe)
	if result.Error != nil {
		logs.Error(result.Error, gorm.ErrRecordNotFound)
		return false, Recipe, 0
	}
	return true, Recipe, count
}

// 新增数据
func (*RecipeService) AddRecipe(pot models.Recipe, list1 []models.Materials, list2 []*models.Taste) (error, models.Recipe) {

	err := psql.Mydb.Transaction(func(tx *gorm.DB) error {

		result := tx.Create(&pot)
		if result.Error != nil {
			logs.Error(result.Error, gorm.ErrRecordNotFound)
			return result.Error
		}
		if err := tx.Create(list1).Error; err != nil {
			logs.Error(err)
			return err
		}
		if len(list2) > 0 {
			if err := tx.Create(list2).Error; err != nil {
				logs.Error(err)
				return err
			}
		}
		return nil

	})
	if err != nil {
		return err, pot
	}

	return nil, pot
}

// 修改数据
func (*RecipeService) EditRecipe(pot valid.EditRecipeValidator, list1 []models.Materials, list2 []*models.Taste, list3 []string, list4 []string) error {
	taste := ""
	//if len(pot.Tastes) == 0 {
	//	taste = ""
	//} else {
	//	taste = strings.Join(pot.Tastes, ",")
	//}

	err := psql.Mydb.Transaction(func(tx *gorm.DB) error {

		err := tx.Exec("UPDATE recipe SET bottom_pot_id = $1,"+
			"bottom_pot= $2,pot_type_id= $3,bottom_properties = $4,"+
			"soup_standard = $5 ,update_at = $6 ,taste_materials = $7 WHERE id = $8",
			pot.BottomPotId, pot.BottomPot, pot.PotTypeId, pot.BottomProperties, pot.SoupStandard, taste, time.Now().Format("2006-01-02 15:04:05"), pot.Id).Error
		if err != nil {
			return err
		}
		if len(list1) > 0 {
			if err = tx.Create(&list1).Error; err != nil {
				return err
			}
		}
		if len(list2) > 0 {
			if err = tx.Create(&list2).Error; err != nil {
				fmt.Println(err)
				return err
			}
		}

		if len(list3) > 0 {

			var taste models.Taste
			if err = tx.Where("id in (?)", list3).Delete(&taste).Error; err != nil {
				return err
			}

		}

		if len(list4) > 0 {
			var material models.Materials
			if err = tx.Where("id in (?)", list4).Delete(&material).Error; err != nil {
				fmt.Println(err)
				return err
			}

		}

		return nil
	})
	fmt.Println(err)
	return err

}

// 删除数据
func (*RecipeService) DeleteRecipe(pot models.Recipe) error {
	return psql.Mydb.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&models.Recipe{}).Where("id = ?", pot.Id).UpdateColumns(map[string]interface{}{"is_del": true, "delete_at": time.Now()}).Error
		if err != nil {
			return err
		}

		var material models.Materials
		err = tx.Where("recipe_id = ?", pot.Id).Where("pot_type_id = ?", pot.PotTypeId).Delete(&material).Error
		if err != nil {
			return err
		}

		var taste models.Taste
		err = tx.Where("recipe_id = ?", pot.Id).Where("pot_type_id = ?", pot.PotTypeId).Delete(&taste).Error
		if err != nil {
			return err
		}
		return nil
	})

}

// 拆分下发配置
func (*RecipeService) SplitSendMqtt(data *mqtt.SendConfig, token string, intervalTime int) error {
	//随机6位字符串
	sendId := utils.GetRandomString(6)
	// 锅型配置
	potTypeConfig := &mqtt.PotTypeConfig{
		SendId:  sendId,
		PotType: data.PotType,
	}
	//发送给mqtt
	bytes, err := json.Marshal(potTypeConfig)
	if err != nil {
		return err
	}
	if err := mqtt.SendToHDL(bytes, token); err != nil {
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
		tasteConfig := &mqtt.TasteConfig{
			SendId: sendId,
			Taste:  data.Taste[i:end],
		}
		bytes, err := json.Marshal(tasteConfig)
		if err != nil {
			return err
		}
		if err := mqtt.SendToHDL(bytes, token); err != nil {
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
		materialsConfig := &mqtt.MaterialsConfig{
			SendId:    sendId,
			Materials: data.Materials[i:end],
		}
		bytes, err := json.Marshal(materialsConfig)
		if err != nil {
			return err
		}
		if err := mqtt.SendToHDL(bytes, token); err != nil {
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
		recipeConfig := &mqtt.RecipeConfig{
			SendId: sendId,
			Recipe: data.Recipe[i:end],
		}
		bytes, err := json.Marshal(recipeConfig)
		if err != nil {
			return err
		}
		if err := mqtt.SendToHDL(bytes, token); err != nil {
			return err
		}
		//等待时间
		time.Sleep(time.Second * time.Duration(intervalTime))
	}
	return nil

}
func (*RecipeService) GetSendToMQTTData(assetId string) (*mqtt.SendConfig, error) {
	var Asset models.Asset
	err2 := psql.Mydb.Where("id = ?", assetId).First(&Asset).Error
	if err2 != nil {
		return nil, err2
	}

	tmpSendConfig := &mqtt.SendConfig{
		Shop: mqtt.ShopContent{
			Name:   Asset.Name,
			Number: Asset.ID,
		},
		PotType:   make([]*mqtt.PotType, 0),
		Taste:     make([]*mqtt.Taste, 0),
		Materials: make([]*mqtt.Materials, 0),
		Recipe:    make([]*mqtt.Recipe, 0),
	}
	var Recipe []*models.Recipe
	//	err := psql.Mydb.Where("is_del = ?", false).Where("asset_id = ?", Asset.ID).Find(&Recipe).Error
	err := psql.Mydb.Where("is_del = ?", false).Find(&Recipe).Error
	if err != nil {
		return nil, err
	}

	if len(Recipe) == 0 {
		return nil, errors.New("该店铺下不存在配方")
	}
	for _, v := range Recipe {
		tmpSendConfig.Recipe = append(tmpSendConfig.Recipe, &mqtt.Recipe{
			BottomPotId:      v.BottomPotId,
			BottomPot:        v.BottomPot,
			PotTypeId:        v.PotTypeId,
			BottomProperties: v.BottomProperties,
			//SoupStandard:     v.SoupStandard,
		})
	}

	recipeIdArr := make([]string, 0)
	potTypeArr := make([]string, 0)
	for _, v := range Recipe {
		recipeIdArr = append(recipeIdArr, v.Id)
		potTypeArr = append(potTypeArr, v.PotTypeId)
	}
	potTypeList := make([]*models.PotType, 0)
	err = psql.Mydb.Where("pot_type_id in (?)", potTypeArr).Find(&potTypeList).Error
	if err != nil {
		return nil, err
	}
	for _, v := range potTypeList {
		tmpSendConfig.PotType = append(tmpSendConfig.PotType, &mqtt.PotType{
			Name:         v.Name,
			SoupStandard: v.SoupStandard,
			PotTypeId:    v.PotTypeId,
		})
	}
	materialList := make([]*models.Materials, 0)
	err = psql.Mydb.Where("recipe_id in (?)", recipeIdArr).Where("resource = ?", "material").Find(&materialList).Error
	if err != nil {
		return nil, err
	}
	tasteMaterialList := make([]*models.Materials, 0)
	err = psql.Mydb.Where("recipe_id in (?)", recipeIdArr).Where("resource = ?", "taste").Find(&tasteMaterialList).Error
	if err != nil {
		return nil, err
	}

	tmpMaterialMap := make(map[string]*mqtt.Materials)
	materialIdList := make(map[string][]string, 0)
	for _, v := range materialList {
		if value, ok := tmpMaterialMap[fmt.Sprintf("%s%d%s%d", v.Name, v.Dosage, v.Unit, v.WaterLine)]; ok {
			materialIdList[v.RecipeID] = append(materialIdList[v.RecipeID], value.Id)
		} else {
			materialIdList[v.RecipeID] = append(materialIdList[v.RecipeID], v.Id)
			tmpMaterialMap[fmt.Sprintf("%s%d%s%d", v.Name, v.Dosage, v.Unit, v.WaterLine)] = &mqtt.Materials{
				Id:        v.Id,
				Name:      v.Name,
				Dosage:    v.Dosage,
				Unit:      v.Unit,
				WaterLine: v.WaterLine,
				Station:   v.Station,
				Resource:  v.Resource,
			}
		}

	}

	for _, v := range tasteMaterialList {
		if _, ok := tmpMaterialMap[fmt.Sprintf("%s%d%s%d", v.Name, v.Dosage, v.Unit, v.WaterLine)]; ok {
			//materialIdList[v.RecipeID] = append(materialIdList[v.RecipeID], value.Id)
		} else {
			//materialIdList[v.RecipeID] = append(materialIdList[v.RecipeID], v.Id)
			tmpMaterialMap[fmt.Sprintf("%s%d%s%d", v.Name, v.Dosage, v.Unit, v.WaterLine)] = &mqtt.Materials{
				Id:        v.Id,
				Name:      v.Name,
				Dosage:    v.Dosage,
				Unit:      v.Unit,
				WaterLine: v.WaterLine,
				Station:   v.Station,
				Resource:  v.Resource,
			}
		}
	}

	for _, v := range tmpMaterialMap {
		tmpSendConfig.Materials = append(tmpSendConfig.Materials, v)
	}

	tasteList := make([]*models.Taste, 0)

	err = psql.Mydb.Where("recipe_id in (?)", recipeIdArr).Where("is_del", false).Find(&tasteList).Error
	if err != nil {
		return nil, err
	}
	tmpTasteMap := make(map[string]*mqtt.Taste)
	tasteMaterialIdArr := make(map[string][]string, 0)
	for _, v := range tasteList {
		tmpTasteMap[v.TasteId] = &mqtt.Taste{
			Name:      v.Name,
			TasteId:   v.TasteId,
			PotTypeId: v.PotTypeId,
			//RecipeId:       v.RecipeID,
			MaterialIdList: strings.Split(v.MaterialIdList, ","),
			BottomPotId:    v.BottomPotId,
		}
		tasteMaterialIdArr[v.TasteId] = append(tasteMaterialIdArr[v.TasteId], v.Id)
	}

	for _, v := range tmpTasteMap {
		tmpSendConfig.Taste = append(tmpSendConfig.Taste, v)
	}

	for key, value := range Recipe {
		tmpSendConfig.Recipe[key].MaterialIdList = materialIdList[value.Id]
	}

	return tmpSendConfig, nil
}

func (*RecipeService) FindMaterialByName(potTypeId, resource string) ([]*models.Materials, error) {
	list := make([]*models.Materials, 0)
	db := psql.Mydb
	if potTypeId != "" {
		db = db.Where("pot_type_id = ?", potTypeId)
	}
	if resource != "" {
		db = db.Where("resource = ?", resource)
	}
	err := db.Find(&list).Error
	if err != nil {
		return nil, err
	}
	list1 := make(map[string]*models.Materials, 0)
	for _, v := range list {
		list1[MD5(fmt.Sprintf("%s%d%s%d", v.Name, v.Dosage, v.Unit, v.WaterLine))] = v
	}
	list2 := make([]*models.Materials, 0)
	for _, v := range list1 {
		list2 = append(list2, v)
	}

	return list2, nil
}

//func (*RecipeService) FindTasteMaterialList(potTypeId string) ([]*models.Taste, error) {
//	list := make([]*models.Taste, 0)
//	db := psql.Mydb
//	if potTypeId != "" {
//		db = db.Where("pot_type_id = ?", potTypeId)
//	}
//	err := db.Find(&list).Error
//	if err != nil {
//		return nil, err
//	}
//	list1 := make(map[string]*models.Taste, 0)
//	for _, v := range list {
//		if v.Material == "" {
//			continue
//		}
//		list1[MD5(fmt.Sprintf("%s%d%s%d", v.Material, v.Dosage, v.Unit, v.WaterLine))] = v
//	}
//	list2 := make([]*models.Taste, 0)
//	for _, v := range list1 {
//		list2 = append(list2, v)
//	}
//
//	return list2, nil
//}

func MD5(data string) string {
	h := md5.New()
	h.Write([]byte(data)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr) // 输出加密结果
}

func (*RecipeService) CreateMaterial(material *models.OriginalMaterials) (bool, error) {
	var createModel models.OriginalMaterials
	err := psql.Mydb.Where("name = ?", material.Name).First(&createModel).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			if psql.Mydb.Create(&material).Error != nil {
				return false, err
			}
			return false, nil
		}

	}
	return true, nil
}

func (*RecipeService) CreateTaste(taste *models.OriginalTaste, action string) error {
	var createModel models.OriginalTaste
	if action == "CHECK" {
		err := psql.Mydb.Where("name = ?", taste.Name).First(&createModel).Error
		if err != nil {
			fmt.Println(strings.Contains(err.Error(), "record not found"))
			if strings.Contains(err.Error(), "record not found") {
				return nil
			}
			return err
		}
	} else {
		if err := psql.Mydb.Create(&taste).Error; err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (*RecipeService) CheckBottomIdIsRepeat(bottomId, recipeId string, action string, potTypeId string) (bool, error) {
	var model models.Recipe
	var err error
	if action == "ADD" {
		err = psql.Mydb.Where("bottom_pot_id = ?", bottomId).Where("pot_type_id = ?", potTypeId).Where("is_del", false).First(&model).Error
	} else {
		err = psql.Mydb.Where("bottom_pot_id = ?", bottomId).Where("id <> ?", recipeId).Where("pot_type_id = ?", potTypeId).Where("is_del", false).First(&model).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil

}

func (*RecipeService) CheckPosTasteIdIsRepeat(list5 string, action string, potTypeIs string, BottomPotId string) (bool, error) {
	if len(list5) > 0 {
		list := make([]*models.Taste, 0)
		if action == "GET" {
			return false, nil
		}
		err := psql.Mydb.Where("taste_id = ?", list5).Where("pot_type_id = ? and bottom_pot_id = ?", potTypeIs, BottomPotId).First(&list).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return false, err
			}
		}
		if len(list) > 0 {
			return true, nil
		}
	}
	return false, nil
}

func FindDiff(arr1, arr2 []string) []string {

	diff := make([]string, 0)

	for _, str1 := range arr1 {
		found := false
		for _, str2 := range arr2 {
			if str1 == str2 {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, str1)
		}
	}

	for _, str2 := range arr2 {
		found := false
		for _, str1 := range arr1 {
			if str2 == str1 {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, str2)
		}
	}

	return diff
}
