package services

import (
	"ThingsPanel-Go/initialize/psql"
	"ThingsPanel-Go/models"
	"ThingsPanel-Go/utils"
	valid "ThingsPanel-Go/validate"
	"errors"
	"fmt"

	"github.com/beego/beego/v2/core/logs"
)

type SoupDataService struct {
	//可搜索字段
	SearchField []string
	//可作为条件的字段
	WhereField []string
	//可做为时间范围查询的字段
	TimeField []string
}

// 获取列表
func (*SoupDataService) GetList(PaginationValidate valid.SoupDataPaginationValidate, tenantId string) ([]map[string]interface{}, int64, error) {
	// 根据租户id获取租户管理员的备注，租户管理员在user表存贮
	var tenantAdmin models.Users
	if tenantId != "SYS_ADMIN" {
		if err := psql.Mydb.Model(&models.Users{}).Where("tenant_id = ? and authority = 'TENANT_ADMIN'", tenantId).First(&tenantAdmin).Error; err != nil {
			return nil, 0, err
		}
	}
	// 店铺名称是租户管理员的名称，店铺id是租户管理员的备注
	var SoupData []models.HdlAddSoupData
	offset := (PaginationValidate.CurrentPage - 1) * PaginationValidate.PerPage
	db := psql.Mydb.Model(&models.HdlAddSoupData{})
	// 如果店铺名称不为空，则根据店铺名称查询店铺id
	if PaginationValidate.ShopName != "" {
		db.Where("shop_name like ?", "%"+PaginationValidate.ShopName+"%")
	}
	// 根据店铺id查询店铺订单数量
	if tenantId != "SYS_ADMIN" {
		if tenantAdmin.Remark != "" {
			db = db.Where("shop_id = ?", tenantAdmin.Remark)
		} else {
			return nil, 0, errors.New("店铺id为空，请先设置店铺id")
		}
	}
	// 如果桌号不为空，则根据桌号查询
	if PaginationValidate.TableNumber != "" {
		db = db.Where("table_number = ?", PaginationValidate.TableNumber)
	}
	// 如果订单检索开始时间不为空，则根据订单检索开始时间查询
	if PaginationValidate.StartTime != "" {
		db = db.Where("creation_time >= ?", PaginationValidate.StartTime)
	}
	// 如果订单检索结束时间不为空，则根据订单检索结束时间查询
	if PaginationValidate.EndTime != "" {
		db = db.Where("creation_time <= ?", PaginationValidate.EndTime)
	}
	// 定义店铺数量
	var count int64
	// 根据店铺id查询店铺订单数量
	db.Count(&count)
	// 根据店铺id查询店铺订单
	result := db.Limit(PaginationValidate.PerPage).Offset(offset).Order("creation_time desc").Find(&SoupData)
	if result.Error != nil {
		logs.Error(result.Error)
		return nil, 0, result.Error
	}
	//定义一个map列表
	var SoupDataMapList []map[string]interface{}
	for _, v := range SoupData {
		// 定义一个map
		SoupDataMap := make(map[string]interface{})

		SoupDataMap["order_sn"] = v.OrderSn         //订单号
		SoupDataMap["table_number"] = v.TableNumber //桌号
		SoupDataMap["shop_id"] = v.ShopId           //店铺id
		SoupDataMap["bottom_id"] = v.BottomId       //锅底id
		SoupDataMap["bottom_pot"] = v.BottomPot     //锅底名称
		SoupDataMap["shop_name"] = v.ShopName       //店铺名称
		SoupDataMap["order_time"] = v.OrderTime     //下单时间
		//SoupDataMap["feeding_end_time"] = v.FeedingEndTime        //投料结束时间
		SoupDataMap["soup_start_time"] = v.SoupStartTime          //加汤开始时间
		SoupDataMap["turning_pot_end_time"] = v.TurningPotEndTime //转锅结束时间
		SoupDataMap["soup_end_time"] = v.SoupEndTime              //加汤结束时间
		//SoupDataMap["feeding_start_time"] = v.FeedingStartTime//投料开始时间
		SoupDataMap["creation_time"] = v.CreationTime //订单创建时间
		SoupDataMapList = append(SoupDataMapList, SoupDataMap)

	}
	return SoupDataMapList, count, nil
	// if PaginationValidate.ShopName != "" {
	// 	asset := &models.Asset{}
	// 	if err := psql.Mydb.Model(&models.Asset{}).Where("name like ?", "%"+PaginationValidate.ShopName+"%").First(asset).Error; err != nil {
	// 		return false, nil, 0
	// 	}
	// 	db = db.Where("add_soup_data.shop_id = ?", asset.ID)
	// }

	// var count int64
	// db.Count(&count)
	// result := db.Model(new(models.AddSoupData)).Select("add_soup_data.bottom_pot,add_soup_data.order_sn,add_soup_data.table_number,add_soup_data.order_time,add_soup_data.soup_start_time,add_soup_data.soup_end_time,add_soup_data.feeding_start_time,add_soup_data.feeding_end_time,add_soup_data.turning_pot_end_time,add_soup_data.turning_pot_end_time,add_soup_data.name order by create_at desc").Limit(PaginationValidate.PerPage).Offset(offset).Find(&SoupData)
	//result := db.Model(new(models.AddSoupData)).Select("add_soup_data.bottom_pot,add_soup_data.order_sn,add_soup_data.table_number,add_soup_data.order_time,add_soup_data.soup_start_time,add_soup_data.soup_end_time,add_soup_data.feeding_start_time,add_soup_data.feeding_end_time,add_soup_data.turning_pot_end_time,add_soup_data.turning_pot_end_time,asset.name").Joins("left join recipe on add_soup_data.bottom_id = recipe.bottom_pot_id").Joins("left join asset on add_soup_data.shop_id = asset.id").Limit(PaginationValidate.PerPage).Offset(offset).Find(&SoupData)
	// if result.Error != nil {
	// 	logs.Error(result.Error, gorm.ErrRecordNotFound)
	// 	return false, SoupData, 0
	// }
	//return true, SoupData, count
}

// 分页查询数据
func (*SoupDataService) Paginate(shopName string, limit int, offset int) ([]models.AddSoupDataValue, int64) {
	tSKVs := []models.SoupDataKVResult{}
	tsk := []models.AddSoupDataValue{}
	var count int64
	result := psql.Mydb
	result2 := psql.Mydb
	if limit <= 0 {
		limit = 1000000
	}
	if offset <= 0 {
		offset = 0
	}
	filters := map[string]interface{}{}

	// if shopName != "" { //店铺id
	// 	Asset := models.Asset{}
	// 	if err := psql.Mydb.Where("name like ?", "%"+shopName+"%").First(&Asset).Error; err != nil {
	// 		return nil, 0
	// 	}
	// 	filters["asd.shop_id"] = Asset.ID
	// }

	SQLWhere, params := utils.TsKvFilterToSql(filters)

	//countsql := "SELECT Count(*) AS count FROM add_soup_data as asd LEFT JOIN asset as a ON asd.shop_id=a.id   " + SQLWhere
	countsql := "SELECT Count(*) AS count FROM add_soup_data" + SQLWhere

	if err := result2.Raw(countsql, params...).Count(&count).Error; err != nil {
		logs.Info(err.Error())
		return tsk, 0
	}
	fmt.Println(countsql)
	//select business.name bname,ts_kv.*,concat_ws('-',asset.name,device.name) AS name,device.token
	//FROM ts_kv LEFT join device on device.id=ts_kv.entity_id
	//LEFT JOIN asset  ON asset.id=device.asset_id
	//LEFT JOIN business ON business.id=asset.business_id
	//WHERE 1=1  and ts_kv.ts >= 1654790400000000 and ts_kv.ts < 1655481599000000 ORDER BY ts_kv.ts DESC limit 10 offset 0
	//SQL := `select add_soup_data.order_sn,asset.name,add_soup_data.table_number,add_soup_data.order_time,add_soup_data.bottom_pot,
	//add_soup_data.soup_start_time,add_soup_data.soup_end_time,add_soup_data.feeding_start_time,add_soup_data.feeding_end_time,add_soup_data.turning_pot_end_time ,asset.name FROM add_soup_data  LEFT JOIN asset  ON add_soup_data.shop_id=asset.id LEFT JOIN recipe on add_soup_data.bottom_id = recipe.bottom_pot_id` + SQLWhere
	SQL := `select add_soup_data.order_sn,add_soup_data.table_number,add_soup_data.order_time,add_soup_data.bottom_pot,
add_soup_data.soup_start_time,add_soup_data.soup_end_time,add_soup_data.feeding_start_time,add_soup_data.feeding_end_time,add_soup_data.turning_pot_end_time ,add_soup_data.name FROM add_soup_data ` + SQLWhere

	if limit > 0 && offset >= 0 {
		SQL = fmt.Sprintf("%s limit ? offset ? ", SQL)
		params = append(params, limit, offset)
	}
	if err := result.Raw(SQL, params...).Scan(&tSKVs).Error; err != nil {
		logs.Error(err.Error())
		return tsk, 0
	}
	for _, v := range tSKVs {
		ts := models.AddSoupDataValue{
			ShopName:         v.ShopName,
			OrderSn:          v.OrderSn,
			BottomPot:        v.BottomPot,
			TableNumber:      v.TableNumber,
			OrderTime:        v.OrderTime,
			SoupStartTime:    v.SoupStartTime,
			SoupEndTime:      v.SoupEndTime,
			FeedingStartTime: v.FeedingStartTime,
			FeedingEndTime:   v.FeedingEndTime,
			TurningPotEnd:    v.TurningPotEnd,
		}
		tsk = append(tsk, ts)
	}
	return tsk, count
}
