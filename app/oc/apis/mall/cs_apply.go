package mall

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/oc/models"
	service "go-admin/app/oc/service/mall"
	"go-admin/app/oc/service/mall/dto"
	"go-admin/common/actions"
)

type CsApply struct {
	api.Api
}

// GetPage 获取商城-售后中心列表
// @Summary 获取商城-售后中心列表
// @Description 获取商城-售后中心列表
// @Tags 商城-售后中心
// @Param csNo query string false "售后申请编号"
// @Param csType query int false "售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限erp迁移老数据使用）"
// @Param orderId query string false "销售订单号"
// @Param csStatus query int false "售后状态：0-待处理、1-已确认、2-已驳回、99-完结"
// @Param userId query int false "提交人id"
// @Param warehouseCode query string false "退货实体仓库code"
// @Param logicWarehouseCode query string false "退货逻辑仓code"
// @Param csSource query string false "mall ,sxyz"
// @Param vendorId query int false "售后单所属货主id"
// @Param vendorSkuCode query string false "售后单所属货主sku"
// @Param isStatements query int false "订单是否存在对账单 0否 1是"
// @Param filterKeyword query int false "关键字搜索"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CsApply}} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/mall/cs-apply [get]
// @Security Bearer
func (e CsApply) GetPage(c *gin.Context) {
	req := dto.CsApplyGetPageReq{}
	s := service.CsApply{}
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
	list := make([]models.CsApply, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商城-售后中心失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetPageDetails 获取商城-售后中心列表(包含售后商品)
// @Summary 获取商城-售后中心列表(包含售后商品)
// @Description 获取商城-售后中心列表(包含售后商品)
// @Tags 商城-售后中心
// @Param csNo query string false "售后申请编号"
// @Param csType query int false "售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限erp迁移老数据使用）"
// @Param orderId query string false "销售订单号"
// @Param csStatus query int false "售后状态：0-待处理、1-已确认、2-已驳回、99-完结"
// @Param userId query int false "提交人id"
// @Param warehouseCode query string false "退货实体仓库code"
// @Param logicWarehouseCode query string false "退货逻辑仓code"
// @Param csSource query string false "mall ,sxyz"
// @Param vendorId query int false "售后单所属货主id"
// @Param vendorSkuCode query string false "售后单所属货主sku"
// @Param isStatements query int false "订单是否存在对账单 0否 1是"
// @Param filterKeyword query int false "关键字搜索"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CsApply}} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/mall/cs-apply/get-page-details [get]
// @Security Bearer
func (e CsApply) GetPageDetails(c *gin.Context) {
	req := dto.CsApplyGetPageReq{}
	s := service.CsApply{}
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
	list := make([]*dto.CsApplyGetPageGroupByOrderCsApply, 0)
	var count int64

	req.UserId = user.GetUserId(c)

	err = s.GetPageDetails(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商城-售后中心失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetPageGroupByOrder 获取商城-售后中心列表-根据订单分组
// @Summary 获取商城-售后中心列表-根据订单分组
// @Description 获取商城-售后中心列表-根据订单分组
// @Tags 商城-售后中心
// @Param csNo query string false "售后申请编号"
// @Param csType query int false "售后类型：0-退货、1-换货、2-退款、3-发票问题、4-技术及资料支持、5-技术及资料支持、6-缺货少配件、7-售后维修、8-其他、9-订单取消、10-补发货（仅限erp迁移老数据使用）"
// @Param orderId query string false "销售订单号"
// @Param csStatus query int false "售后状态：0-待处理、1-已确认、2-已驳回、99-完结"
// @Param userId query int false "提交人id"
// @Param warehouseCode query string false "退货实体仓库code"
// @Param logicWarehouseCode query string false "退货逻辑仓code"
// @Param csSource query string false "mall ,sxyz"
// @Param vendorId query int false "售后单所属货主id"
// @Param vendorSkuCode query string false "售后单所属货主sku"
// @Param isStatements query int false "订单是否存在对账单 0否 1是"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.CsApply}} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/mall/cs-apply/group-by-order [get]
// @Security Bearer
func (e CsApply) GetPageGroupByOrder(c *gin.Context) {
	req := dto.CsApplyGetPageReq{}
	s := service.CsApply{}
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
	list := make([]*dto.CsApplyGetPageGroupByOrder, 0)
	var count int64

	req.UserId = user.GetUserId(c)

	err = s.GetPageGroupByOrder(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取商城-售后中心失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Cancel 拒绝售后单
// @Summary 拒绝售后单
// @Description 拒绝售后单
// @Tags 商城-售后中心
// @Accept application/json
// @Product application/json
// @Param data body dto.CsApplyCancelReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/oc/mall/cs-apply/cancel [put]
// @Security Bearer
func (e CsApply) Cancel(c *gin.Context) {
	req := dto.CsApplyCancelReq{}
	s := service.CsApply{}
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

	req.UserId = user.GetUserId(c)
	err = s.Cancel(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改售后单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.CsNo, "关闭成功")
}

// Get 获取商城-售后中心
// @Summary 获取商城-售后中心
// @Description 获取商城-售后中心
// @Tags 商城-售后中心
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.CsApply} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/mall/cs-apply/{id} [get]
// @Security Bearer
func (e CsApply) Get(c *gin.Context) {
	req := dto.CsApplyGetReq{}
	s := service.CsApply{}
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
	var object dto.CsApplyGetPageGroupByOrderCsApply

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取售后信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建商城-售后中心
// @Summary 创建商城-售后中心
// @Description 创建商城-售后中心
// @Tags 商城-售后中心
// @Accept application/json
// @Product application/json
// @Param data body dto.CsApplyInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/oc/mall/cs-apply [post]
// @Security Bearer
func (e CsApply) Insert(c *gin.Context) {
	req := dto.CsApplyInsertReq{}
	s := service.CsApply{}
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
		e.Error(500, err, fmt.Sprintf("创建商城-售后中心失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改商城-售后中心
// @Summary 修改商城-售后中心
// @Description 修改商城-售后中心
// @Tags 商城-售后中心
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CsApplyUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/oc/mall/cs-apply/{id} [put]
// @Security Bearer
func (e CsApply) Update(c *gin.Context) {
	req := dto.CsApplyUpdateReq{}
	s := service.CsApply{}
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
		e.Error(500, err, fmt.Sprintf("修改商城-售后中心失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除商城-售后中心
// @Summary 删除商城-售后中心
// @Description 删除商城-售后中心
// @Tags 商城-售后中心
// @Param data body dto.CsApplyDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/oc/mall/cs-apply [delete]
// @Security Bearer
func (e CsApply) Delete(c *gin.Context) {
	s := service.CsApply{}
	req := dto.CsApplyDeleteReq{}
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
		e.Error(500, err, fmt.Sprintf("删除商城-售后中心失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
