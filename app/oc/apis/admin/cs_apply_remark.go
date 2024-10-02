package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/oc/models"
	service "go-admin/app/oc/service/admin"
	"go-admin/app/oc/service/admin/dto"
	"go-admin/common/actions"
)

type CsApplyRemark struct {
	api.Api
}

// GetPage 获取售后备注列表
// @Summary 获取售后备注列表
// @Description 获取售后备注列表
// @Tags 售后备注
// @Param type query string false "类型{sale_order:销售售后，purchase_order:采购售后}"
// @Param dataId query string true "订单号 必填"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.DataRemark}} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/admin/cs-apply/remark [get]
// @Security Bearer
func (e CsApplyRemark) GetPage(c *gin.Context) {
	req := dto.DataRemarkGetPageReq{}
	s := service.DataRemark{}
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
	list := make([]models.DataRemark, 0)
	var count int64
	req.Type = "sale_order"
	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取操作备注记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Insert 创建售后备注记录
// @Summary 创建售后备注记录
// @Description 创建售后备注记录
// @Tags 售后备注
// @Accept application/json
// @Product application/json
// @Param data body dto.DataRemarkInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/oc/admin/cs-apply/remark [post]
// @Security Bearer
func (e CsApplyRemark) Insert(c *gin.Context) {
	req := dto.DataRemarkInsertReq{}
	s := service.DataRemark{}
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
	req.Type = "sale_order"
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建操作备注记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改操作备注记录
// @Summary 修改操作备注记录
// @Description 修改操作备注记录
// @Tags 操作备注记录
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.DataRemarkUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/oc/data-remark/{id} [put]
// @Security Bearer
func (e CsApplyRemark) Update(c *gin.Context) {
	req := dto.DataRemarkUpdateReq{}
	s := service.DataRemark{}
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
		e.Error(500, err, fmt.Sprintf("修改操作备注记录失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除操作备注记录
// @Summary 删除操作备注记录
// @Description 删除操作备注记录
// @Tags 操作备注记录
// @Param data body dto.DataRemarkDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/oc/data-remark [delete]
// @Security Bearer
func (e CsApplyRemark) Delete(c *gin.Context) {
	s := service.DataRemark{}
	req := dto.DataRemarkDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除操作备注记录失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
