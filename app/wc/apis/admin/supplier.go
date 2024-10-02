package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/wc/models"
	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
)

type Supplier struct {
	api.Api
}

// GetPage 获取供应商列表
// @Summary 获取供应商列表
// @Description 获取供应商列表
// @Tags 供应商
// @Param code query string false "供应商编码"
// @Param nameZh query string false "供应商中文名"
// @Param nameEn query string false "供应商英文名"
// @Param shortName query string false "供应商简称"
// @Param cityId query string false "市id"
// @Param provinceId query string false "省id"
// @Param countryId query string false "国家id"
// @Param companyId query string false "关联公司"
// @Param status query string false "0无效 1有效"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Supplier}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/supplier [get]
// @Security Bearer
func (e Supplier) GetPage(c *gin.Context) {
	req := dto.SupplierGetPageReq{}
	s := service.Supplier{}
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
	list := make([]models.Supplier, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取货主失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取供应商
// @Summary 获取供应商
// @Description 获取供应商
// @Tags 供应商
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Supplier} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/supplier/{id} [get]
// @Security Bearer
func (e Supplier) Get(c *gin.Context) {
	req := dto.SupplierGetReq{}
	s := service.Supplier{}
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
	var object models.Supplier

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取货主失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建供应商
// @Summary 创建供应商
// @Description 创建供应商
// @Tags 供应商
// @Accept application/json
// @Product application/json
// @Param data body dto.SupplierInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/wc/admin/supplier [post]
// @Security Bearer
func (e Supplier) Insert(c *gin.Context) {
	req := dto.SupplierInsertReq{}
	s := service.Supplier{}
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
	req.SetCreateByName(user.GetUserName(c))

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建货主失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改供应商
// @Summary 修改供应商
// @Description 修改供应商
// @Tags 供应商
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.SupplierUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/wc/admin/supplier/{id} [put]
// @Security Bearer
func (e Supplier) Update(c *gin.Context) {
	req := dto.SupplierUpdateReq{}
	s := service.Supplier{}
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
	req.SetUpdateByName(user.GetUserName(c))

	p := actions.GetPermissionFromContext(c)

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改供应商失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Select 获取供应商下拉列表
// @Summary 获取供应商下拉列表
// @Description 获取供应商下拉列表
// @Tags 供应商
// @Param id query int false "供应商id"
// @Param code query string false "供应商编码"
// @Param nameZh query string false "供应商中文名"
// @Param shortName query string false "供应商简称"
// @Param status query string false "是否有效 0-否，1-是"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.SupplierSelectResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/supplier/select [get]
// @Security Bearer
func (e Supplier) Select(c *gin.Context) {
	req := dto.SupplierSelectReq{}
	s := service.Supplier{}
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
	list := make([]dto.SupplierSelectResp, 0)
	var count int64

	err = s.Select(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取供应商失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// InnerGetList Inner批量获取供应商信息
// @Summary Inner批量获取供应商信息
// @Description Inner批量获取供应商信息
// @Tags Inner供应商
// @Param ids query string false "供应商id 多个用','分割"
// @Param code query string false "供应商编码 多个用','分割"
// @Param nameZh query string false "供应商中文名  多个用','分割"
// @Param status query string false "0 无效 1 有效"
// @Success 200 {object} response.Response{data=[]models.Supplier} "{"code": 200, "data": [...]}"
// @Router /inner/wc/admin/supplier/list [get]
func (e Supplier) InnerGetList(c *gin.Context) {
	req := dto.InnerSupplierGetListReq{}
	s := service.Supplier{}
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

	list := make([]models.Supplier, 0)

	err = s.InnerGetList(&req, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("Inner批量获取供应商信息，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(list, "成功")
}
