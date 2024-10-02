package mall

import (
	"fmt"
	"go-admin/app/oc/models"
	service "go-admin/app/oc/service/mall"
	"go-admin/app/oc/service/mall/dto"
	"go-admin/common/actions"
	"go-admin/common/excel"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
)

type OrderInfo struct {
	api.Api
}

// GetPage 获取领用单列表
// @Summary 获取领用单列表
// @Description 获取领用单列表
// @Tags 商城领用单
// @Param keyword query string false "关键词"
// @Param orderType query int false "订单类型"
// @Param startTime query string false "下单日期开始"
// @Param endTime query string false "下单日期结束"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.OrderInfoGetPageResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/mall/order [get]
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
		e.Error(500, err, fmt.Sprintf("获取领用单列表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetPageCount 获取领用单列表各状态数量
// @Summary 获取领用单列表各状态数量
// @Description 获取领用单列表各状态数量
// @Tags 商城领用单
// @Success 200 {object} response.Response{data=map[string]int} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/mall/order/count [get]
// @Security Bearer
func (e OrderInfo) GetPageCount(c *gin.Context) {
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

	err = s.GetPageCount(&req, p, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取领用单列表统计失败，\r\n失败信息 %s", err.Error()))
		return
	}
	resp := map[string]int{
		"orderType5": 0,
		"orderType6": 0,
		"orderType1": 0,
		"orderType7": 0,
		"orderType9": 0,
	}
	for _, v := range list {
		// 小程序-仓库处理中 mall-待发货
		if v.OrderStatus == 5 || v.OrderStatus == 6 || v.OrderStatus == 0 {
			resp["orderType5"]++
		} else if v.OrderStatus == 0 {
			resp["orderType6"]++
			//resp["orderType1"]++
		} else if v.OrderStatus == 1 || v.OrderStatus == 11 { // 待收货
			resp["orderType1"]++
		} else if v.OrderStatus == 7 { // 已完成
			resp["orderType7"]++
		} else if v.OrderStatus == 9 || v.OrderStatus == 10 { // 已取消
			resp["orderType9"]++
		} else if v.OrderStatus == 11 {
			resp["orderType1"]++
			resp["orderType6"]++
		}
	}

	e.OK(resp, "查询成功")
}

// Get 获取领用单
// @Summary 获取领用单
// @Description 获取领用单
// @Tags 商城领用单
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.OrderInfoGetResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/mall/order/{id} [get]
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

// UpdatePo 修改Po单号
// @Summary 修改Po单号
// @Description 修改Po单号
// @Tags 商城领用单
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.OrderInfoReceiptReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改Po单号成功"}"
// @Router /api/v1/oc/mall/order/update-po/{id} [put]
// @Security Bearer
func (e OrderInfo) UpdatePo(c *gin.Context) {
	req := dto.OrderInfoUpdatePoReq{}
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

	err = s.UpdatePo(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改Po单号失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改Po单号成功")
}

// Cancel 申请取消
// @Summary 申请取消
// @Description 申请取消
// @Tags 商城领用单
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.OrderInfoCancelReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "取消成功"}"
// @Router /api/v1/oc/mall/order/cancel/{id} [put]
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
		e.Error(500, err, fmt.Sprintf("取消失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "取消成功")
}

// BuyAgain 再次购买
// @Summary 再次购买
// @Description 再次购买
// @Tags 商城领用单
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.OrderInfoBuyAgainReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "取消成功"}"
// @Router /api/v1/oc/mall/order/buy-again/{id} [put]
// @Security Bearer
func (e OrderInfo) BuyAgain(c *gin.Context) {
	req := dto.OrderInfoBuyAgainReq{}
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

	errList, err := s.BuyAgain(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("加入失败，\r\n失败信息 %s", err.Error()))
		return
	} else if len(errList) > 0 {
		e.OK("warning", "部分商品加入失败："+strings.Join(errList, "\r\n"))
	} else {
		e.OK("success", "加入成功")
	}
}

// Export 导出订单明细报表
// @Summary 导出订单明细报表
// @Description 导出订单明细报表
// @Tags 商城领用单
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.OrderInfoGetPageResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/oc/mall/order/export [get]
// @Security Bearer
func (e OrderInfo) Export(c *gin.Context) {
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
	list := make([]dto.OrderInfoGetExportData, 0)

	err = s.GetExport(&req, p, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取订单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	var exportData []interface{}

	for _, resp := range list {
		exportData = append(exportData, resp)
	}

	title := []map[string]string{
		{"orderId": "订单编号"},
		{"externalOrderNo": "审批单编号"},
		{"contractNo": "审批编号（客户方）"},
		{"userCompanyName": "下单公司"},
		{"userName": "下单用户"},
		{"createdTime": "下单时间"},
		{"warehouseName": "发货仓库"},
		{"consignee": "收货人"},
		{"sendStatusText": "订单发货状态"},
		{"receiveStatusText": "订单收货状态"},
		{"productName": "产品名称"},
		{"userProductRemark": "SKU备注"},
		//{"skuCode": "行号"},
		{"skuCode": "SKU编号"},
		{"vendorName": "货主"},
		{"supplierSkuCode": "货主SKU号"},
		{"productNo": "客户物料编码"},
		{"brandName": "品牌"},
		{"productModel": "型号"},
		{"catName": "一级产线"},
		{"catName2": "二级产线"},
		{"unit": "包装单位"},
		{"taxPrice": "含税单价"},
		{"quantity": "购买数量"},
		{"totalTaxPrice": "含税总价"},
		{"cancelQuantity": "取消数量"},
		{"returnQuantity": "退货数量"},
	}
	excelApp := excel.NewExcel()
	err = excelApp.ExportExcelByStruct(c, title, exportData, "我的订单-个人报表", "Sheet1")
	if err != nil {
		e.Error(500, err, fmt.Sprintf("订单明细导出失败： %s", err.Error()))
	}
}

// Delete 删除订单
// @Summary 删除订单
// @Description 删除订单
// @Tags 商城领用单
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.OrderInfoDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/oc/mall/order/delete/{id} [put]
// @Security Bearer
func (e OrderInfo) Delete(c *gin.Context) {
	req := dto.OrderInfoDeleteReq{}
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

	err = s.Delete(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// Insert 创建领用单
// @Summary 创建领用单
// @Description 创建领用单
// @Tags 商城领用单
// @Accept application/json
// @Product application/json
// @Param data body dto.PresoSubmitApprovalReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/oc/mall/order [post]
// @Security Bearer
func (e OrderInfo) Insert(c *gin.Context) {
	req := dto.PresoSubmitApprovalReq{}
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
	//// 设置创建人
	//req.CreateBy = user.GetUserId(c)
	//req.CreateByName = user.GetUserName(c)

	var order models.OrderInfo
	err = s.Insert(&req, &order)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("领用失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(order, "领用成功")
}

// Receipt 签收领用单
// @Summary 签收领用单
// @Description 签收领用单
// @Tags 商城领用单
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
