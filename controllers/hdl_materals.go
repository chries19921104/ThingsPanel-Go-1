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

type HdlMaterialsController struct {
	beego.Controller
}

// 列表
func (c *HdlMaterialsController) List() {
	reqData := valid.HdlMaterialsPaginationValidate{}
	if err := valid.ParseAndValidate(&c.Ctx.Input.RequestBody, &reqData); err != nil {
		utils.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
		return
	}
	var HdlMaterialsService services.HdlMaterialsService
	isSuccess, d, t := HdlMaterialsService.GetHdlMaterialsList(reqData)
	if !isSuccess {
		utils.SuccessWithMessage(1000, "查询失败", (*context2.Context)(c.Ctx))
		return
	}
	dd := valid.RspHdlMaterialsPaginationValidate{
		CurrentPage: reqData.CurrentPage,
		Data:        d,
		Total:       t,
		PerPage:     reqData.PerPage,
	}
	utils.SuccessWithDetailed(200, "success", dd, map[string]string{}, (*context2.Context)(c.Ctx))
}

// 编辑
func (c *HdlMaterialsController) Edit() {
	reqData := valid.EditHdlMaterialsValidate{}
	if err := valid.ParseAndValidate(&c.Ctx.Input.RequestBody, &reqData); err != nil {
		utils.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
		return
	}
	var HdlMaterialsService services.HdlMaterialsService
	err := HdlMaterialsService.EditHdlMaterials(reqData)
	if err == nil {
		d := HdlMaterialsService.GetHdlMaterialsDetail(reqData.Id)
		utils.SuccessWithDetailed(200, "success", d, map[string]string{}, (*context2.Context)(c.Ctx))
	} else {
		utils.SuccessWithMessage(400, err.Error(), (*context2.Context)(c.Ctx))
	}
}

// 新增
func (c *HdlMaterialsController) Add() {
	reqData := valid.AddHdlMaterialsValidate{}
	if err := valid.ParseAndValidate(&c.Ctx.Input.RequestBody, &reqData); err != nil {
		utils.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
		return
	}
	var HdlMaterialsService services.HdlMaterialsService
	d, rsp_err := HdlMaterialsService.AddHdlMaterials(reqData)
	if rsp_err == nil {
		utils.SuccessWithDetailed(200, "success", d, map[string]string{}, (*context2.Context)(c.Ctx))
	} else {
		utils.SuccessWithMessage(400, rsp_err.Error(), (*context2.Context)(c.Ctx))
	}
}

// 删除
func (HdlMaterialsController *HdlMaterialsController) Delete() {
	HdlMaterialsIdValidate := valid.HdlMaterialsIdValidate{}
	err := json.Unmarshal(HdlMaterialsController.Ctx.Input.RequestBody, &HdlMaterialsIdValidate)
	if err != nil {
		fmt.Println("参数解析失败", err.Error())
	}
	v := validation.Validation{}
	status, _ := v.Valid(HdlMaterialsIdValidate)
	if !status {
		for _, err := range v.Errors {
			// 获取字段别称
			alias := gvalid.GetAlias(HdlMaterialsIdValidate, err.Field)
			message := strings.Replace(err.Message, err.Field, alias, 1)
			utils.SuccessWithMessage(1000, message, (*context2.Context)(HdlMaterialsController.Ctx))
			break
		}
		return
	}
	if HdlMaterialsIdValidate.Id == "" {
		utils.SuccessWithMessage(1000, "id不能为空", (*context2.Context)(HdlMaterialsController.Ctx))
	}
	var HdlMaterialsService services.HdlMaterialsService
	HdlMaterials := models.HdlMaterials{
		Id: HdlMaterialsIdValidate.Id,
	}
	req_err := HdlMaterialsService.DeleteHdlMaterials(HdlMaterials)
	if req_err == nil {
		utils.SuccessWithMessage(200, "success", (*context2.Context)(HdlMaterialsController.Ctx))
	} else {
		utils.SuccessWithMessage(400, "删除失败", (*context2.Context)(HdlMaterialsController.Ctx))
	}
}
