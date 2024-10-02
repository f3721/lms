package admin

import (
	"fmt"
	"go-admin/app/uc/service/admin"
	"go-admin/app/uc/service/admin/dto"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/uc/models"
	"go-admin/common/actions"
)

type CompanyIndividualitySwitch struct {
	api.Api
}

// GetPage 获取公司开关表列表
// @Summary 获取公司开关表列表
// @Description 获取公司开关表列表
// @Tags 公司开关表
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CompanyIndividualitySwitch}} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/company-individuality-switch [get]
// @Security Bearer
func (e CompanyIndividualitySwitch) GetPage(c *gin.Context) {
	req := dto.CompanyIndividualitySwitchGetPageReq{}
	s := admin.CompanyIndividualitySwitch{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.CompanyIndividualitySwitch, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取公司开关表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取公司开关表
// @Summary 获取公司开关表
// @Description 获取公司开关表
// @Tags 公司开关表
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.CompanyIndividualitySwitch} "{"code": 200, "data": [...]}"
// @Router /api/v1/uc/admin/company-individuality-switch/{id} [get]
// @Security Bearer
func (e CompanyIndividualitySwitch) Get(c *gin.Context) {
	req := dto.CompanyIndividualitySwitchGetReq{}
	s := admin.CompanyIndividualitySwitch{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object models.CompanyIndividualitySwitch

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取公司开关表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建公司开关表
// @Summary 创建公司开关表
// @Description 创建公司开关表
// @Tags 公司开关表
// @Accept application/json
// @Product application/json
// @Param data body dto.CompanyIndividualitySwitchInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/uc/admin/company-individuality-switch [post]
// @Security Bearer
func (e CompanyIndividualitySwitch) Insert(c *gin.Context) {
	req := dto.CompanyIndividualitySwitchInsertReq{}
	s := admin.CompanyIndividualitySwitch{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建公司开关表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改公司开关表
// @Summary 修改公司开关表
// @Description 修改公司开关表
// @Tags 公司开关表
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CompanyIndividualitySwitchUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/uc/admin/company-individuality-switch/{id} [put]
// @Security Bearer
func (e CompanyIndividualitySwitch) Update(c *gin.Context) {
	req := dto.CompanyIndividualitySwitchUpdateReq{}
	s := admin.CompanyIndividualitySwitch{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改公司开关表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除公司开关表
// @Summary 删除公司开关表
// @Description 删除公司开关表
// @Tags 公司开关表
// @Param data body dto.CompanyIndividualitySwitchDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/uc/admin/company-individuality-switch [delete]
// @Security Bearer
func (e CompanyIndividualitySwitch) Delete(c *gin.Context) {
	s := admin.CompanyIndividualitySwitch{}
	req := dto.CompanyIndividualitySwitchDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除公司开关表失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
