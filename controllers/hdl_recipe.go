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

type HdlRecipeController struct {
	beego.Controller
}

// 列表
func (c *HdlRecipeController) List() {
	reqData := valid.HdlRecipePaginationValidate{}
	if err := valid.ParseAndValidate(&c.Ctx.Input.RequestBody, &reqData); err != nil {
		utils.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
		return
	}
	var HdlRecipeService services.HdlRecipeService
	isSuccess, d, t := HdlRecipeService.GetHdlRecipeList(reqData)
	if !isSuccess {
		utils.SuccessWithMessage(1000, "查询失败", (*context2.Context)(c.Ctx))
		return
	}
	dd := valid.RspHdlRecipePaginationValidate{
		CurrentPage: reqData.CurrentPage,
		Data:        d,
		Total:       t,
		PerPage:     reqData.PerPage,
	}
	utils.SuccessWithDetailed(200, "success", dd, map[string]string{}, (*context2.Context)(c.Ctx))
}

// 列表查询包括物料、口味列表
func (c *HdlRecipeController) EntireList() {
	reqData := valid.HdlRecipePaginationValidate{}
	if err := valid.ParseAndValidate(&c.Ctx.Input.RequestBody, &reqData); err != nil {
		utils.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
		return
	}
	// 获取用户租户id
	tenantId, ok := c.Ctx.Input.GetData("tenant_id").(string)
	if !ok {
		utils.SuccessWithMessage(400, "代码逻辑错误,未获取到租户id", (*context2.Context)(c.Ctx))
		return
	}
	var HdlRecipeService services.HdlRecipeService
	isSuccess, d, t := HdlRecipeService.GetHdlRecipeEntireList(reqData, tenantId)
	if !isSuccess {
		utils.SuccessWithMessage(1000, "查询失败", (*context2.Context)(c.Ctx))
		return
	}
	dd := valid.RspHdlEntireRecipePaginationValidate{
		CurrentPage: reqData.CurrentPage,
		Data:        d,
		Total:       t,
		PerPage:     reqData.PerPage,
	}
	utils.SuccessWithDetailed(200, "success", dd, map[string]string{}, (*context2.Context)(c.Ctx))
}

// 编辑
func (c *HdlRecipeController) Edit() {
	reqData := valid.EditHdlRecipeValidate{}
	if err := valid.ParseAndValidate(&c.Ctx.Input.RequestBody, &reqData); err != nil {
		utils.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
		return
	}
	var HdlRecipeService services.HdlRecipeService
	err := HdlRecipeService.EditHdlRecipe(reqData)
	if err == nil {
		d := HdlRecipeService.GetHdlRecipeDetail(reqData.Id)
		utils.SuccessWithDetailed(200, "success", d, map[string]string{}, (*context2.Context)(c.Ctx))
	} else {
		utils.SuccessWithMessage(400, err.Error(), (*context2.Context)(c.Ctx))
	}
}

// 新增整个配方
func (c *HdlRecipeController) EntireAdd() {
	reqData := valid.AddEntireHdlRecipeValidate{}
	if err := valid.ParseAndValidate(&c.Ctx.Input.RequestBody, &reqData); err != nil {
		utils.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
		return
	}
	// 获取用户租户id
	tenantId, ok := c.Ctx.Input.GetData("tenant_id").(string)
	if !ok {
		utils.SuccessWithMessage(400, "代码逻辑错误,未获取到租户id", (*context2.Context)(c.Ctx))
		return
	}
	var HdlRecipeService services.HdlRecipeService
	d, rsp_err := HdlRecipeService.AddHdlRecipeAndMateralsAndTaste(reqData, tenantId)
	if rsp_err == nil {
		utils.SuccessWithDetailed(200, "success", d, map[string]string{}, (*context2.Context)(c.Ctx))
	} else {
		utils.SuccessWithMessage(400, rsp_err.Error(), (*context2.Context)(c.Ctx))
	}
}

// 新增
func (c *HdlRecipeController) Add() {
	reqData := valid.AddHdlRecipeValidate{}
	if err := valid.ParseAndValidate(&c.Ctx.Input.RequestBody, &reqData); err != nil {
		utils.SuccessWithMessage(1000, err.Error(), (*context2.Context)(c.Ctx))
		return
	}
	var HdlRecipeService services.HdlRecipeService
	d, rsp_err := HdlRecipeService.AddHdlRecipe(reqData)
	if rsp_err == nil {
		utils.SuccessWithDetailed(200, "success", d, map[string]string{}, (*context2.Context)(c.Ctx))
	} else {
		utils.SuccessWithMessage(400, rsp_err.Error(), (*context2.Context)(c.Ctx))
	}
}

// 删除
func (HdlRecipeController *HdlRecipeController) Delete() {
	HdlRecipeIdValidate := valid.HdlRecipeIdValidate{}
	err := json.Unmarshal(HdlRecipeController.Ctx.Input.RequestBody, &HdlRecipeIdValidate)
	if err != nil {
		fmt.Println("参数解析失败", err.Error())
	}
	v := validation.Validation{}
	status, _ := v.Valid(HdlRecipeIdValidate)
	if !status {
		for _, err := range v.Errors {
			// 获取字段别称
			alias := gvalid.GetAlias(HdlRecipeIdValidate, err.Field)
			message := strings.Replace(err.Message, err.Field, alias, 1)
			utils.SuccessWithMessage(1000, message, (*context2.Context)(HdlRecipeController.Ctx))
			break
		}
		return
	}
	if HdlRecipeIdValidate.Id == "" {
		utils.SuccessWithMessage(1000, "id不能为空", (*context2.Context)(HdlRecipeController.Ctx))
	}
	var HdlRecipeService services.HdlRecipeService
	HdlRecipe := models.HdlRecipe{
		Id: HdlRecipeIdValidate.Id,
	}
	req_err := HdlRecipeService.DeleteHdlRecipe(HdlRecipe)
	if req_err == nil {
		utils.SuccessWithMessage(200, "success", (*context2.Context)(HdlRecipeController.Ctx))
	} else {
		utils.SuccessWithMessage(400, "删除失败", (*context2.Context)(HdlRecipeController.Ctx))
	}
}
