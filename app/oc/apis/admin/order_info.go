package admin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"go-admin/app/oc/models"
	service "go-admin/app/oc/service/admin"
	"go-admin/app/oc/service/admin/dto"
	"go-admin/common/actions"
	"strings"
)

type OrderInfo struct {
	api.Api
}

// GetPage 获取领用单列表
// @Summary 获取领用单列表
// @Description 获取领用单列表
// @Tags 领用单
// @Param orderIds query string false "订单号"
// @Param orderStatus query []int false "订单状态"
// @Param rmaStatus query int false "售后状态"
// @Param createFrom query string false "订单来源"
// @Param createBy query int false "用户ID"
// @Param startTime query string false "下单日期开始"
// @Param endTime query string false "下单日期结束"
// @Param isOverBudget query int false "是否超出预算"
// @Param actualOrderingCompanyId query int false "公司ID"
// @Param skuCode query string false "SKU"
// @Param productName query string false "商品名称"
// @Param warehouseCode query string false "发货仓库"
// @Param productNo query string false "物料编码"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.OrderInfoGetPageResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/order [get]
// @Security Bearer
func (e OrderInfo) GetPage(c *gin.Context) {
	req := dto.OrderInfoGetPageReq{}
	s := service.OrderInfo{}
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
	list := make([]dto.OrderInfoGetPageResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取领用单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取领用单
// @Summary 获取领用单
// @Description 获取领用单
// @Tags 领用单
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.OrderInfoGetResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/order/{id} [get]
// @Security Bearer
func (e OrderInfo) Get(c *gin.Context) {
	req := dto.OrderInfoGetReq{}
	s := service.OrderInfo{}
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
	//var object models.OrderInfo
	var object dto.OrderInfoGetResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取领用单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建领用单
// @Summary 创建领用单
// @Description 创建领用单
// @Tags 领用单
// @Accept application/json
// @Product application/json
// @Param data body dto.OrderInfoInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/oc/order [post]
// @Security Bearer
func (e OrderInfo) Insert(c *gin.Context) {
	req := dto.OrderInfoInsertReq{}
	s := service.OrderInfo{}
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
	req.CreateBy = user.GetUserId(c)
	req.CreateByName = user.GetUserName(c)

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建领用单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改领用单基本信息
// @Summary 修改领用单基本信息
// @Description 修改领用单基本信息
// @Tags 领用单
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.OrderInfoUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/oc/order/{id} [put]
// @Security Bearer
func (e OrderInfo) Update(c *gin.Context) {
	req := dto.OrderInfoUpdateReq{}
	s := service.OrderInfo{}
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

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改领用单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Receipt 签收领用单
// @Summary 签收领用单
// @Description 签收领用单
// @Tags 领用单
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.OrderInfoReceiptReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "签收成功"}"
// @Router /api/v1/oc/order/receipt/{id} [put]
// @Security Bearer
func (e OrderInfo) Receipt(c *gin.Context) {
	req := dto.OrderInfoReceiptReq{}
	s := service.OrderInfo{}
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

	err = s.Receipt(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("签收领用单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "签收成功")
}

// GetReceiptImage 订单-获取回单图片信息
// @Summary 订单-获取回单图片信息
// @Description 订单-获取回单图片信息
// @Tags 领用单
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.OrderInfoGetReceiptImageReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "签收成功"}"
// @Router /api/v1/oc/order/receipt-image/{orderId} [get]
// @Security Bearer
func (e OrderInfo) GetReceiptImage(c *gin.Context) {
	req := dto.OrderInfoGetReceiptImageReq{}
	res := dto.OrderInfoGetReceiptImageRes{}
	s := service.OrderInfo{}
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

	err = s.GetReceiptImage(&req, &res)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("签收领用单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(res, "成功")
}

// SaveReceiptImage 订单-回单上传（未签收时可多次上传）
// @Summary 订单-回单上传
// @Description 订单-回单上传
// @Tags 领用单
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.OrderInfoSaveReceiptImageReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "签收成功"}"
// @Router /api/v1/oc/order/receipt-image/{orderId} [put]
// @Security Bearer
func (e OrderInfo) SaveReceiptImage(c *gin.Context) {
	req := dto.OrderInfoSaveReceiptImageReq{}
	s := service.OrderInfo{}
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

	err = s.SaveReceiptImage(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("签收领用单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(nil, "上传成功")
}

// Cancel 取消领用单
// @Summary 取消领用单
// @Description 取消领用单
// @Tags 领用单
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.OrderInfoCancelReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "取消成功"}"
// @Router /api/v1/oc/order/cancel/{OrderId} [put]
// @Security Bearer
func (e OrderInfo) Cancel(c *gin.Context) {
	req := dto.OrderInfoCancelReq{}
	s := service.OrderInfo{}
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

	// 设置取消人
	req.CancelBy = user.GetUserId(c)
	req.CancelByName = user.GetUserName(c)

	err = s.Cancel(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("取消领用单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// UpdateShipping 修改领用单配送
// @Summary 修改领用单配送
// @Description 修改领用单配送
// @Tags 领用单
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.OrderInfoUpdateShippingReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/oc/order/shipping/{id} [put]
// @Security Bearer
func (e OrderInfo) UpdateShipping(c *gin.Context) {
	req := dto.OrderInfoUpdateShippingReq{}
	s := service.OrderInfo{}
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

	err = s.UpdateShipping(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改领用单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// UpdateProduct 修改领用单商品
// @Summary 修改领用单商品
// @Description 修改领用单商品
// @Tags 领用单
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.OrderInfoUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/oc/order/product/{id} [put]
// @Security Bearer
func (e OrderInfo) UpdateProduct(c *gin.Context) {
	req := dto.OrderInfoUpdateProductReq{}
	s := service.OrderInfo{}
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

	err = s.UpdateProduct(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改领用单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// GetByOrderId 获取领用单
// @Summary 获取领用单
// @Description 获取领用单
// @Tags 领用单
// @Param orderIds query string true "订单号"
// @Success 200 {object} response.Response{data=[]models.OrderInfo} "{"code": 200, "data": [...]}"
// @Router /inner/oc/order [get]
// @Security Bearer
func (e OrderInfo) GetByOrderId(c *gin.Context) {
	req := dto.OrderInfoByOrderIdReq{}
	s := service.OrderInfo{}
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
	//if len(req.OrderIds) <= 0 {
	//	e.Error(500, nil, "订单id必填")
	//	return
	//}

	var object []models.OrderInfo

	err = s.GetByOrderId(&req, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取领用单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	if len(object) == 1 {
		e.OK(object[0], "查询成功")
	} else {
		e.OK(object, "查询成功")
	}
}

// CheckExistUnCompletedOrder 判断是否有未完结订单
// @Summary 判断是否有未完结订单
// @Description 判断是否有未完结订单
// @Tags 领用单
// @Param orderIds query string true "订单号"
// @Success 200 {object} response.Response{data=[]models.OrderInfo} "{"code": 200, "data": [...]}"
// @Router /inner/oc/order [get]
// @Security Bearer
func (e OrderInfo) CheckExistUnCompletedOrder(c *gin.Context) {
	req := dto.OrderInfoCheckExistUnCompletedOrderReq{}
	s := service.OrderInfo{}
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

	var object []models.OrderInfo

	err = s.CheckExistUnCompletedOrder(&req, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取失败，\r\n失败信息 %s", err.Error()))
		return
	}

	if len(object) > 0 {
		e.OK(400, "存在未完结订单")
	} else {
		e.OK(200, "不存在未完结订单")
	}
}

// AddProduct 领用单添加商品
// @Summary 领用单添加商品
// @Description 领用单添加商品
// @Tags 领用单
// @Accept application/json
// @Product application/json
// @Param data body dto.OrderInfoAddProductReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/oc/order/addProduct [post]
// @Security Bearer
func (e OrderInfo) AddProduct(c *gin.Context) {
	req := dto.OrderInfoAddProductReq{}
	s := service.OrderInfo{}
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
	data, err := s.AddProduct(&req)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	} else {
		e.OK(data, "获取成功")
	}
}

// CheckIsOverBudget 领用单是否超预算
// @Summary 领用单是否超预算
// @Description 领用单是否超预算
// @Tags 领用单
// @Accept application/json
// @Product application/json
// @Param data body dto.OrderIdsReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "无订单超预算"}"
// @Router /api/v1/oc/order/check-over-budget [post]
// @Security Bearer
func (e OrderInfo) CheckIsOverBudget(c *gin.Context) {
	req := dto.OrderIdsReq{}
	s := service.OrderInfo{}
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
	data, err := s.CheckIsOverBudget(&req)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	} else if len(data) > 0 {
		err = errors.New(strings.Join(data, "\r\n"))
		e.OK("error", "部分订单超预算："+err.Error())
	} else {
		e.OK("success", "无订单超预算")
	}
}

// Confirm 领用单确认
// @Summary 领用单确认
// @Description 领用单确认
// @Tags 领用单
// @Accept application/json
// @Product application/json
// @Param data body dto.OrderIdsReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/oc/order/confirm [post]
// @Security Bearer
func (e OrderInfo) Confirm(c *gin.Context) {
	req := dto.OrderIdsReq{}
	s := service.OrderInfo{}
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
	data, err := s.Confirm(&req)
	if err != nil {
		e.Error(500, err, err.Error())
		return
	} else if len(data) > 0 {
		err = errors.New(strings.Join(data, "\r\n"))
		e.OK("error", "部分订单确认失败："+err.Error())
	} else {
		e.OK(nil, "确认订单成功")
	}
}

// AutoConfirm 订单确认的定时任务：(5分钟一次)
func (e OrderInfo) AutoConfirm(c *gin.Context) {
	s := service.OrderInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	echoMsg, err := s.AutoConfirm(c)
	if err != nil {
		c.String(500, err.Error())
	} else {
		c.String(200, echoMsg)
	}
}

// AutOfStock 缺货订单补货：（5分钟一次）
func (e OrderInfo) AutOfStock(c *gin.Context) {
	s := service.OrderInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	echoMsg, err := s.AutOfStock(c)
	if err != nil {
		c.String(500, err.Error())
	} else {
		c.String(200, echoMsg)
	}
}

// AutoSignFor 订单自动签收
func (e OrderInfo) AutoSignFor(c *gin.Context) {
	s := service.OrderInfo{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	echoMsg, err := s.AutoSignFor(c, nil)
	if err != nil {
		c.String(500, err.Error())
	} else {
		c.String(200, echoMsg)
	}
}

// ------------------------------------------------------INNER-----------------------------------------------------------

// GetOrderListByUserName 获取领用单
// @Summary 获取领用单
// @Description 获取领用单
// @Tags 领用单
// @Param userName query string true "领用人"
// @Success 200 {object} response.Response{data=[]dto.OrderListRes} "{"code": 200, "data": [...]}"
// @Router /inner/oc/order/list [get]
func (e OrderInfo) GetOrderListByUserName(c *gin.Context) {
	req := dto.OrderListReq{}
	s := service.OrderInfo{}
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
	var object []dto.OrderListResp

	err = s.GetOrderListByUserName(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取领用单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	if len(object) == 1 {
		e.OK(object[0], "查询成功")
	} else {
		e.OK(object, "查询成功")
	}
}
