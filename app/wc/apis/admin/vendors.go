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

type Vendors struct {
	api.Api
}

// GetPage 获取货主列表
// @Summary 获取货主列表
// @Description 获取货主列表
// @Tags 货主
// @Param code query string false "货主编码"
// @Param nameZh query string false "货主中文名"
// @Param nameEn query string false "供应商英文名"
// @Param shortName query string false "供应商简称"
// @Param city query string false "市id"
// @Param province query string false "省id"
// @Param country query string false "国家id"
// @Param companyId query string false "关联公司"
// @Param status query string false "0无效 1有效"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Vendors}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/vendors [get]
// @Security Bearer
func (e Vendors) GetPage(c *gin.Context) {
	req := dto.VendorsGetPageReq{}
	s := service.Vendors{}
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
	list := make([]models.Vendors, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取货主失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取货主
// @Summary 获取货主
// @Description 获取货主
// @Tags 货主
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Vendors} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/vendors/{id} [get]
// @Security Bearer
func (e Vendors) Get(c *gin.Context) {
	req := dto.VendorsGetReq{}
	s := service.Vendors{}
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
	var object models.Vendors

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取货主失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建货主
// @Summary 创建货主
// @Description 创建货主
// @Tags 货主
// @Accept application/json
// @Product application/json
// @Param data body dto.VendorsInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/wc/admin/vendors [post]
// @Security Bearer
func (e Vendors) Insert(c *gin.Context) {
	req := dto.VendorsInsertReq{}
	s := service.Vendors{}
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

// Update 修改货主
// @Summary 修改货主
// @Description 修改货主
// @Tags 货主
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.VendorsUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/wc/admin/vendors/{id} [put]
// @Security Bearer
func (e Vendors) Update(c *gin.Context) {
	req := dto.VendorsUpdateReq{}
	s := service.Vendors{}
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
		e.Error(500, err, fmt.Sprintf("修改货主失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// GetPage 获取货主下拉列表
// @Summary 获取货主下拉列表
// @Description 获取货主下拉列表
// @Tags 货主
// @Param id query int false "货主id"
// @Param code query string false "货主编码"
// @Param nameZh query string false "货主中文名"
// @Param shortName query string false "货主简称"
// @Param status query string false "是否有效 0-否，1-是"
// @Param noCheckPermission query int false "是否不校验货主权限"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.VendorsSelectResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/vendors/select [get]
// @Security Bearer
func (e Vendors) Select(c *gin.Context) {
	req := dto.VendorsSelectReq{}
	s := service.Vendors{}
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
	list := make([]dto.VendorsSelectResp, 0)
	var count int64

	err = s.Select(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取货主失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetPage Inner批量获取货主信息
// @Summary Inner批量获取货主信息
// @Description Inner批量获取货主信息
// @Tags Inner货主
// @Param ids query string false "货主id 多个用','分割"
// @Param code query string false "货主编码 多个用','分割"
// @Param nameZh query string false "货主中文名  多个用','分割"
// @Param status query string false "0 无效 1 有效"
// @Success 200 {object} response.Response{data=[]models.Vendors} "{"code": 200, "data": [...]}"
// @Router /inner/wc/admin/vendors/list [get]
func (e Vendors) InnerGetList(c *gin.Context) {
	req := dto.InnerVendorsGetListReq{}
	s := service.Vendors{}
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

	list := make([]models.Vendors, 0)

	err = s.InnerGetList(&req, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("Inner批量获取货主信息，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(list, "成功")
}
