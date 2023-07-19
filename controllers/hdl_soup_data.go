package controllers

import (
	gvalid "ThingsPanel-Go/initialize/validate"
	"ThingsPanel-Go/services"
	response "ThingsPanel-Go/utils"
	valid "ThingsPanel-Go/validate"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	context2 "github.com/beego/beego/v2/server/web/context"
	"github.com/mintance/go-uniqid"
	"github.com/xuri/excelize/v2"
)

type SoupDataController struct {
	beego.Controller
}

func (soup *SoupDataController) Index() {
	PaginationValidate := valid.SoupDataPaginationValidate{}
	err := json.Unmarshal(soup.Ctx.Input.RequestBody, &PaginationValidate)
	if err != nil {
		fmt.Println("参数解析失败", err.Error())
	}
	v := validation.Validation{}
	status, _ := v.Valid(PaginationValidate)
	if !status {
		for _, err := range v.Errors {
			// 获取字段别称
			alias := gvalid.GetAlias(PaginationValidate, err.Field)
			message := strings.Replace(err.Message, err.Field, alias, 1)
			response.SuccessWithMessage(1000, message, (*context2.Context)(soup.Ctx))
			break
		}
		return
	}
	var tenantId string
	// 判断是否是系统管理员
	authority, ok := soup.Ctx.Input.GetData("authority").(string)
	if !ok {
		response.SuccessWithMessage(400, "代码逻辑错误,未获取到用户权限", (*context2.Context)(soup.Ctx))
		return
	}
	if authority == "SYS_ADMIN" {
		tenantId = "SYS_ADMIN"
	} else {
		// 获取用户租户id
		tenantId, ok = soup.Ctx.Input.GetData("tenant_id").(string)
		if !ok {
			response.SuccessWithMessage(400, "代码逻辑错误,未获取到租户id", (*context2.Context)(soup.Ctx))
			return
		}
	}

	var SoupDataService services.SoupDataService
	d, t, err := SoupDataService.GetList(PaginationValidate, tenantId)
	if err != nil {
		response.SuccessWithMessage(1000, "查询失败", (*context2.Context)(soup.Ctx))
		return
	}
	dd := valid.RspSoupDataPaginationValidate{
		CurrentPage: PaginationValidate.CurrentPage,
		Data:        d,
		Total:       t,
		PerPage:     PaginationValidate.PerPage,
	}
	response.SuccessWithDetailed(200, "success", dd, map[string]string{}, (*context2.Context)(soup.Ctx))

}

func (soup *SoupDataController) Export() {
	SoupDataExcelValidate := valid.SoupDataPaginationValidate{}
	err := json.Unmarshal(soup.Ctx.Input.RequestBody, &SoupDataExcelValidate)
	if err != nil {
		fmt.Println("参数解析失败", err.Error())
	}
	v := validation.Validation{}
	status, _ := v.Valid(SoupDataExcelValidate)
	if !status {
		for _, err := range v.Errors {
			// 获取字段别称
			alias := gvalid.GetAlias(SoupDataExcelValidate, err.Field)
			message := strings.Replace(err.Message, err.Field, alias, 1)
			response.SuccessWithMessage(1000, message, (*context2.Context)(soup.Ctx))
			break
		}
		return
	}
	var tenantId string
	// 判断是否是系统管理员
	authority, ok := soup.Ctx.Input.GetData("authority").(string)
	if !ok {
		response.SuccessWithMessage(400, "代码逻辑错误,未获取到用户权限", (*context2.Context)(soup.Ctx))
		return
	}
	if authority == "SYS_ADMIN" {
		tenantId = "SYS_ADMIN"
	} else {
		// 获取用户租户id
		tenantId, ok = soup.Ctx.Input.GetData("tenant_id").(string)
		if !ok {
			response.SuccessWithMessage(400, "代码逻辑错误,未获取到租户id", (*context2.Context)(soup.Ctx))
			return
		}
	}
	var TSKVService services.SoupDataService
	//每次查10000条
	num := SoupDataExcelValidate.Limit / 10000
	excel_file := excelize.NewFile()
	index := excel_file.NewSheet("Sheet1")
	excel_file.SetActiveSheet(index)
	excel_file.SetCellValue("Sheet1", "A1", "店名")
	excel_file.SetCellValue("Sheet1", "B1", "订单号")
	excel_file.SetCellValue("Sheet1", "C1", "锅底名称")
	excel_file.SetCellValue("Sheet1", "D1", "桌号")
	excel_file.SetCellValue("Sheet1", "E1", "订单时间")
	excel_file.SetCellValue("Sheet1", "F1", "开始加汤时间")
	excel_file.SetCellValue("Sheet1", "G1", "加汤完毕时间")
	//excel_file.SetCellValue("Sheet1", "H1", "加料完成时间")
	excel_file.SetCellValue("Sheet1", "H1", "传锅完成时间")
	for i := 0; i <= num; i++ {
		var t []map[string]interface{}
		var c int64
		if (i+1)*10000 <= SoupDataExcelValidate.Limit {
			SoupDataExcelValidate.CurrentPage = i + 1
			SoupDataExcelValidate.PerPage = 10000
			t, c, err = TSKVService.GetList(SoupDataExcelValidate, tenantId)
		} else {
			SoupDataExcelValidate.CurrentPage = i + 1
			SoupDataExcelValidate.PerPage = SoupDataExcelValidate.Limit % 10000
			t, c, err = TSKVService.GetList(SoupDataExcelValidate, tenantId)
		}
		var i int
		if c > 0 {
			i = 1
			for _, tv := range t {
				i++
				is := strconv.Itoa(i)
				excel_file.SetCellValue("Sheet1", "A"+is, tv["shop_name"])
				excel_file.SetCellValue("Sheet1", "B"+is, tv["order_sn"])
				excel_file.SetCellValue("Sheet1", "C"+is, tv["bottom_pot"])
				excel_file.SetCellValue("Sheet1", "D"+is, tv["table_number"])
				excel_file.SetCellValue("Sheet1", "E"+is, tv["creation_time"])
				excel_file.SetCellValue("Sheet1", "F"+is, tv["soup_start_time"])
				excel_file.SetCellValue("Sheet1", "G"+is, tv["soup_end_time"])
				//excel_file.SetCellValue("Sheet1", "H"+is, tv["feeding_end_time"])
				excel_file.SetCellValue("Sheet1", "H"+is, tv["turning_pot_end_time"])
			}
		}
	}

	uploadDir := "./files/excel/"
	errs := os.MkdirAll(uploadDir, os.ModePerm)
	if errs != nil {
		response.SuccessWithMessage(1000, err.Error(), (*context2.Context)(soup.Ctx))
	}
	// 根据指定路径保存文件
	uniqid_str := uniqid.New(uniqid.Params{Prefix: "excel", MoreEntropy: true})
	excelName := "files/excel/数据列表" + uniqid_str + ".xlsx"
	if err := excel_file.SaveAs(excelName); err != nil {
		fmt.Println(err)
	}
	response.SuccessWithDetailed(200, "获取成功", excelName, map[string]string{}, (*context2.Context)(soup.Ctx))
}
