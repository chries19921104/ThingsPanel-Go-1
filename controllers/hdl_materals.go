package controllers

import (
	gvalid "ThingsPanel-Go/initialize/validate"
	"ThingsPanel-Go/models"
	"ThingsPanel-Go/services"
	"ThingsPanel-Go/utils"
	valid "ThingsPanel-Go/validate"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	context2 "github.com/beego/beego/v2/server/web/context"
)

type HdlMateralsController struct {
	beego.Controller
}

// 列表
func (c *HdlMateralsController) List() {
	reqData := valid.HdlMateralsPaginationValidate{}
	if err := valid.ParseAndValidate(&c.Ctx.Input.RequestBody, &reqData); err != nil {
		utils.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
		return
	}
	var HdlMateralsService services.HdlMateralsService
	isSuccess, d, t := HdlMateralsService.GetHdlMateralsList(reqData)
	if !isSuccess {
		utils.SuccessWithMessage(1000, "查询失败", (*context2.Context)(c.Ctx))
		return
	}
	dd := valid.RspHdlMateralsPaginationValidate{
		CurrentPage: reqData.CurrentPage,
		Data:        d,
		Total:       t,
		PerPage:     reqData.PerPage,
	}
	utils.SuccessWithDetailed(200, "success", dd, map[string]string{}, (*context2.Context)(c.Ctx))
}

// 编辑
func (HdlMateralsController *HdlMateralsController) Edit() {
	HdlMateralsValidate := valid.EditHdlMateralsValidate{}
	err := json.Unmarshal(HdlMateralsController.Ctx.Input.RequestBody, &HdlMateralsValidate)
	if err != nil {
		fmt.Println("参数解析失败", err.Error())
	}
	v := validation.Validation{}
	status, _ := v.Valid(HdlMateralsValidate)
	if !status {
		for _, err := range v.Errors {
			// 获取字段别称
			alias := gvalid.GetAlias(HdlMateralsValidate, err.Field)
			message := strings.Replace(err.Message, err.Field, alias, 1)
			utils.SuccessWithMessage(1000, message, (*context2.Context)(HdlMateralsController.Ctx))
			break
		}
		return
	}
	if HdlMateralsValidate.Id == "" {
		utils.SuccessWithMessage(1000, "id不能为空", (*context2.Context)(HdlMateralsController.Ctx))
	}
	var HdlMateralsService services.HdlMateralsService
	isSucess := HdlMateralsService.EditHdlMaterals(HdlMateralsValidate)
	if isSucess {
		d := HdlMateralsService.GetHdlMateralsDetail(HdlMateralsValidate.Id)
		utils.SuccessWithDetailed(200, "success", d, map[string]string{}, (*context2.Context)(HdlMateralsController.Ctx))
	} else {
		utils.SuccessWithMessage(400, "编辑失败", (*context2.Context)(HdlMateralsController.Ctx))
	}
}

// 新增
func (c *HdlMateralsController) Add() {
	reqData := valid.AddHdlMateralsValidate{}
	if err := valid.ParseAndValidate(&c.Ctx.Input.RequestBody, &reqData); err != nil {
		utils.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
		return
	}
	var HdlMateralsService services.HdlMateralsService
	d, rsp_err := HdlMateralsService.AddHdlMaterals(reqData)
	if rsp_err == nil {
		utils.SuccessWithDetailed(200, "success", d, map[string]string{}, (*context2.Context)(c.Ctx))
	} else {
		utils.SuccessWithMessage(400, rsp_err.Error(), (*context2.Context)(c.Ctx))
	}
}

// 删除
func (HdlMateralsController *HdlMateralsController) Delete() {
	HdlMateralsIdValidate := valid.HdlMateralsIdValidate{}
	err := json.Unmarshal(HdlMateralsController.Ctx.Input.RequestBody, &HdlMateralsIdValidate)
	if err != nil {
		fmt.Println("参数解析失败", err.Error())
	}
	v := validation.Validation{}
	status, _ := v.Valid(HdlMateralsIdValidate)
	if !status {
		for _, err := range v.Errors {
			// 获取字段别称
			alias := gvalid.GetAlias(HdlMateralsIdValidate, err.Field)
			message := strings.Replace(err.Message, err.Field, alias, 1)
			utils.SuccessWithMessage(1000, message, (*context2.Context)(HdlMateralsController.Ctx))
			break
		}
		return
	}
	if HdlMateralsIdValidate.Id == "" {
		utils.SuccessWithMessage(1000, "id不能为空", (*context2.Context)(HdlMateralsController.Ctx))
	}
	var HdlMateralsService services.HdlMateralsService
	HdlMaterals := models.HdlMaterals{
		Id: HdlMateralsIdValidate.Id,
	}
	req_err := HdlMateralsService.DeleteHdlMaterals(HdlMaterals)
	if req_err == nil {
		utils.SuccessWithMessage(200, "success", (*context2.Context)(HdlMateralsController.Ctx))
	} else {
		utils.SuccessWithMessage(400, "删除失败", (*context2.Context)(HdlMateralsController.Ctx))
	}
}
