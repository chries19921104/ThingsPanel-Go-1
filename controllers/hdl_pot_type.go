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

type HdlPotTypeController struct {
	beego.Controller
}

// 列表
func (c *HdlPotTypeController) List() {
	reqData := valid.HdlPotTypePaginationValidate{}
	if err := valid.ParseAndValidate(&c.Ctx.Input.RequestBody, &reqData); err != nil {
		utils.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
		return
	}
	var HdlPotTypeService services.HdlPotTypeService
	isSuccess, d, t := HdlPotTypeService.GetHdlPotTypeList(reqData)
	if !isSuccess {
		utils.SuccessWithMessage(1000, "查询失败", (*context2.Context)(c.Ctx))
		return
	}
	dd := valid.RspHdlPotTypePaginationValidate{
		CurrentPage: reqData.CurrentPage,
		Data:        d,
		Total:       t,
		PerPage:     reqData.PerPage,
	}
	utils.SuccessWithDetailed(200, "success", dd, map[string]string{}, (*context2.Context)(c.Ctx))
}

// 编辑
func (c *HdlPotTypeController) Edit() {
	reqData := valid.EditHdlPotTypeValidate{}
	if err := valid.ParseAndValidate(&c.Ctx.Input.RequestBody, &reqData); err != nil {
		utils.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
		return
	}
	var HdlPotTypeService services.HdlPotTypeService
	err := HdlPotTypeService.EditHdlPotType(reqData)
	if err == nil {
		d := HdlPotTypeService.GetHdlPotTypeDetail(reqData.Id)
		utils.SuccessWithDetailed(200, "success", d, map[string]string{}, (*context2.Context)(c.Ctx))
	} else {
		utils.SuccessWithMessage(400, err.Error(), (*context2.Context)(c.Ctx))
	}
}

// 新增
func (c *HdlPotTypeController) Add() {
	reqData := valid.AddHdlPotTypeValidate{}
	if err := valid.ParseAndValidate(&c.Ctx.Input.RequestBody, &reqData); err != nil {
		utils.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
		return
	}
	var HdlPotTypeService services.HdlPotTypeService
	d, rsp_err := HdlPotTypeService.AddHdlPotType(reqData)
	if rsp_err == nil {
		utils.SuccessWithDetailed(200, "success", d, map[string]string{}, (*context2.Context)(c.Ctx))
	} else {
		utils.SuccessWithMessage(400, rsp_err.Error(), (*context2.Context)(c.Ctx))
	}
}

// 删除
func (HdlPotTypeController *HdlPotTypeController) Delete() {
	HdlPotTypeIdValidate := valid.HdlPotTypeIdValidate{}
	err := json.Unmarshal(HdlPotTypeController.Ctx.Input.RequestBody, &HdlPotTypeIdValidate)
	if err != nil {
		fmt.Println("参数解析失败", err.Error())
	}
	v := validation.Validation{}
	status, _ := v.Valid(HdlPotTypeIdValidate)
	if !status {
		for _, err := range v.Errors {
			// 获取字段别称
			alias := gvalid.GetAlias(HdlPotTypeIdValidate, err.Field)
			message := strings.Replace(err.Message, err.Field, alias, 1)
			utils.SuccessWithMessage(1000, message, (*context2.Context)(HdlPotTypeController.Ctx))
			break
		}
		return
	}
	if HdlPotTypeIdValidate.Id == "" {
		utils.SuccessWithMessage(1000, "id不能为空", (*context2.Context)(HdlPotTypeController.Ctx))
	}
	var HdlPotTypeService services.HdlPotTypeService
	HdlPotType := models.HdlPotType{
		Id: HdlPotTypeIdValidate.Id,
	}
	req_err := HdlPotTypeService.DeleteHdlPotType(HdlPotType)
	if req_err == nil {
		utils.SuccessWithMessage(200, "success", (*context2.Context)(HdlPotTypeController.Ctx))
	} else {
		utils.SuccessWithMessage(400, "删除失败", (*context2.Context)(HdlPotTypeController.Ctx))
	}
}
