package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	service "go-admin/app/oc/service/admin"
	"go-admin/app/oc/service/admin/dto"
	"go-admin/common/actions"
)

type CsApply struct {
	api.Api
}

// GetPage 获取售后单列表
// @Summary 获取售后单列表
// @Description 获取售后单列表
// @Tags 售后中心
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
// @Param filterSkuCode query string false "商品sku"
// @Param filterProductNo query string false "物料编码"
// @Param filterProductName query string false "商品名"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.CsApplyData}} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/admin/cs-apply [get]
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
	list := make([]*dto.CsApplyListData, 0)
	var count int64

	err = s.GetListPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取售后单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取售后单
// @Summary 获取售后单
// @Description 获取售后单
// @Tags 售后中心
// @Param csNo path string false "csNo"
// @Success 200 {object} response.Response{data=dto.CsApplyData} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/admin/cs-apply/{csNo} [get]
// @Security Bearer
func (e CsApply) Get(c *gin.Context) {
	req := dto.CsApplyGetReq{}
	s := service.CsApply{}
	csApplyDetailService := service.CsApplyDetail{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		MakeService(&csApplyDetailService.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object dto.CsApplyInfoData

	p := actions.GetPermissionFromContext(c)
	err = s.GetInfo(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取售后单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建售后单
// @Summary 创建售后单
// @Description 创建售后单
// @Tags 售后中心
// @Accept application/json
// @Product application/json
// @Param data body dto.CsApplyInsertRequest true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/oc/admin/cs-apply [post]
// @Security Bearer
func (e CsApply) Insert(c *gin.Context) {
	req := dto.CsApplyInsertRequest{}
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

	p := actions.GetPermissionFromContext(c)
	err = s.Add(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建售后单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(nil, "创建成功")
}

// Undo 关闭售后单
// @Summary 关闭售后单
// @Description 关闭售后单
// @Tags 售后中心
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CsApplyCancelReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/oc/admin/cs-apply/un-do [put]
// @Security Bearer
func (e CsApply) Undo(c *gin.Context) {
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

	err = s.Undo(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改售后单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.CsNo, "作废成功")
}

// Cancel 拒绝售后单
// @Summary 拒绝售后单
// @Description 拒绝售后单
// @Tags 售后中心
// @Accept application/json
// @Product application/json
// @Param data body dto.CsApplyCancelReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/oc/admin/cs-apply/cancel [put]
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

	err = s.Cancel(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改售后单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.CsNo, "关闭成功")
}

// Confirm 确认审核售后单
// @Summary 确认审核售后单
// @Description 确认审核售后单
// @Tags 售后中心
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CsApplyConfirmReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/oc/admin/cs-apply/confirm [put]
// @Security Bearer
func (e CsApply) Confirm(c *gin.Context) {
	req := dto.CsApplyConfirmReq{}
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

	err, callStockChangeEventErr := s.Confirm(e.Orm, req.CsNo, p)
	if callStockChangeEventErr {
		e.Error(500, err, fmt.Sprintf("审核售后单失败，\n失败信息 库存中心出库入库变更失败 %s", err.Error()))
		return
	}
	if err != nil {
		e.Error(500, err, fmt.Sprintf("审核售后单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.CsNo, "售后审核成功")
}

// AllAudit 批量审核售后丹
// @Summary 批量审核售后丹
// @Description 批量审核售后丹
// @Tags 售后中心
// @Accept application/json
// @Product application/json
// @Param data body dto.CsApplyAllAuditReq true "body"
// @Success 200 {object} response.Response{data=dto.CsApplyAllAuditRes}	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/oc/admin/cs-apply/all-audit/{id} [put]
// @Security Bearer
func (e CsApply) AllAudit(c *gin.Context) {
	req := dto.CsApplyAllAuditReq{}
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

	warning := ""
	msgType := "success"

	// 作废
	if req.OperationType == "undo" {
		for _, csNo := range req.CsNos {
			// 作废操作
			err = s.Undo(&dto.CsApplyCancelReq{
				CsNo:        csNo,
				AuditReason: req.AuditReason,
			}, p)
			if err != nil {
				msgType = "warning"
				warning += csNo + "作废失败提示：" + err.Error() + "<br />"
			} else {
				warning += csNo + "作废成功！<br />"
			}
		}
	} else if req.OperationType == "confirm" { // 确认
		for _, csNo := range req.CsNos {
			// 确认操作
			err, callStockChangeEventErr := s.Confirm(nil, csNo, p)
			if callStockChangeEventErr {
				msgType = "warning"
				warning += csNo + "售后单库存中心出库入库变更失败：" + err.Error() + "<br />"
			} else if err != nil {
				msgType = "warning"
				warning += csNo + "确认失败提示：" + err.Error() + "<br />"
			} else {
				warning += csNo + "确认成功！<br />"
			}
		}
	}

	res := dto.CsApplyAllAuditRes{
		Msg:     warning,
		MsgType: msgType,
	}
	e.OK(res, warning)
}

// Update 修改售后单
// @Summary 修改售后单
// @Description 修改售后单
// @Tags 售后中心
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.CsApplyUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/oc/admin/cs-apply/{id} [put]
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
		e.Error(500, err, fmt.Sprintf("修改售后单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除售后单
// @Summary 删除售后单
// @Description 删除售后单
// @Tags 售后中心
// @Param data body dto.CsApplyDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/oc/admin/cs-apply [delete]
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
		e.Error(500, err, fmt.Sprintf("删除售后单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// GetSaleProducts 申请售后时-获取订单商品信息
// @Summary 申请售后时-获取订单商品信息
// @Description 申请售后时-获取订单商品信息
// @Tags 售后中心
// @Param orderId path string false "orderId"
// @Param csType query int false "csType"
// @Success 200 {object} response.Response{data=dto.CsApplyGetSaleProductsRes} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/admin/cs-apply/sale-products/{orderId} [get]
// @Security Bearer
func (e CsApply) GetSaleProducts(c *gin.Context) {
	req := dto.CsApplyGetSaleProductsReq{}
	s := service.CsApply{}
	csApplyDetailService := service.CsApplyDetail{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		MakeService(&csApplyDetailService.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	object, err := s.GetSaleProducts(req.OrderId, req.CsType)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取售后单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// GetOrderInfoByCsNo 根据售后单号获取订单信息
// @Summary 根据售后单号获取订单信息
// @Description 根据售后单号获取订单信息
// @Tags 售后中心
// @Param csNo path string false "csNo"
// @Success 200 {object} response.Response{data=models.OrderInfo} "{"code": 200, "data": [...]}"
// @Router /inner/oc/order/by-cs-no/{csNo} [get]
// @Security Bearer
func (e CsApply) GetOrderInfoByCsNo(c *gin.Context) {
	req := dto.GetOrderInfoByCsReq{}
	s := service.CsApply{}
	csApplyDetailService := service.CsApplyDetail{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		MakeService(&csApplyDetailService.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	object, err := s.GetOrderInfoByCsNo(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取售后单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// IsOrderInAfterPendingReview 根据订单号查询是否有订单再售后审核中inner
// @Summary 根据订单号查询是否有订单再售后中inner
// @Description 根据订单号查询是否有订单再售后中inner
// @Tags 售后中心
// @Param orderId path string false "orderId"
// @Success 200 {object} response.Response{data=dto.CsApplyIsOrderInAfterPendingReviewRes} "{"code": 200, "data": [...]}"
// @Router /inner/oc/order/is-order-in-after-pending-review/{orderId} [get]
// @Security Bearer
func (e CsApply) IsOrderInAfterPendingReview(c *gin.Context) {
	req := dto.CsApplyIsOrderInAfterPendingReview{}
	s := service.CsApply{}
	csApplyDetailService := service.CsApplyDetail{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		MakeService(&csApplyDetailService.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	isPendingReview, err := s.IsOrderInAfterPendingReview(req.OrderId)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取订单是否再售后审核中失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(isPendingReview, "查询成功")
}
