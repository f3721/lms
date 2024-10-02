package admin

import (
	"encoding/json"
	"fmt"
	"go-admin/common/excel"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
)

type StockOutbound struct {
	api.Api
}

// GetPage 获取出库单列表
// @Summary 获取出库单列表
// @Description 获取出库单列表
// @Tags 出库单
// @Param outboundCode query string false "出库单编码"
// @Param type query string false "出库类型:  0 大货出库  1 订单出库  2 其他"
// @Param status query string false "状态:0-已作废 1-创建 2-已完成"
// @Param sourceCode query string false "来源单据code"
// @Param outboundTime query time.Time false "出库时间"
// @Param warehouseCode query string false "实体仓code"
// @Param logicWarehouseCode query string false "逻辑仓code"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.StockOutboundGetPageResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-outbound [get]
// @Security Bearer
func (e StockOutbound) GetPage(c *gin.Context) {
	req := dto.StockOutboundGetPageReq{}
	s := service.StockOutbound{}
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
	list := make([]dto.StockOutboundGetPageResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取出库单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

func (e StockOutbound) Export(c *gin.Context) {
	req := dto.StockOutboundGetPageReq{}
	s := service.StockOutbound{}
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

	outData, err := s.Export(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("产品维护导出失败： %s", err.Error()))
		return
	}
	title := []map[string]string{
		{"typeName": "出库类型"},
		{"statusName": "单据状态"},
		{"warehouseName": "实体仓"},
		{"logicWarehouseName": "逻辑仓"},
		{"recipient": "领用人"},
		{"outboundCode": "单据编号"},
		{"sourceCode": "来源单据编号"},
		{"sourceTypeName": "来源单据类型"},
		{"remark": "备注"},
		{"skuCode": "SKU"},
		{"productName": "商品名称"},
		{"mfgModel": "型号"},
		{"brandName": "品牌"},
		{"salesUom": "单位"},
		{"productNo": "物料编码"},
		{"vendorName": "货主名称"},
		{"vendorSkuCode": "货主SKU"},
		{"locationQuantity": "应出库数量"},
		{"locationActQuantity": "实际出库数量"},
		{"diffNum": "差异数量"},
		{"locationCode": "库位编号"},
		{"outboundTime": "行出库时间"},
	}
	exportData := []map[string]interface{}{}
	byteJson, _ := json.Marshal(outData)
	_ = json.Unmarshal(byteJson, &exportData)
	excelApp := excel.NewExcel()
	err = excelApp.ExportExcelByMap(e.Context, title, exportData, "出库单-导出", "出库单明细")
	if err != nil {
		e.Error(500, err, fmt.Sprintf("出库单导出失败： %s", err.Error()))
	}
}

// Get 获取出库单
// @Summary 获取出库单
// @Description 获取出库单
// @Tags 出库单
// @Param id path int false "id"
// @Param type path string false "eg: confirm"
// @Success 200 {object} response.Response{data=dto.StockOutboundGetResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-outbound/{id} [get]
// @Security Bearer
func (e StockOutbound) Get(c *gin.Context) {
	req := dto.StockOutboundGetReq{}
	s := service.StockOutbound{}
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
	var object dto.StockOutboundGetResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取出库单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// ConfirmOutbound 确认出库
// @Summary 确认出库
// @Description 确认出库
// @Tags 出库单
// @Accept application/json
// @Product application/json
// @Param data body dto.StockOutboundConfirmReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "出库成功"}"
// @Router /api/v1/wc/admin/stock-outbound/confirm-outbound [post]
// @Security Bearer
func (e StockOutbound) ConfirmOutbound(c *gin.Context) {
	req := dto.StockOutboundConfirmReq{}
	s := service.StockOutbound{}
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

	err = s.ConfirmOutbound(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("出库失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "出库成功")
}

// ConfirmOutbound 部分出库
// @Summary 部分出库
// @Description 部分出库
// @Tags 出库单
// @Accept application/json
// @Product application/json
// @Param data body dto.StockPartOutboundReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "操作成功"}"
// @Router /api/v1/wc/admin/stock-outbound/confirm-outbound [post]
// @Security Bearer
func (e StockOutbound) PartOutbound(c *gin.Context) {
	req := dto.StockPartOutboundReq{}
	s := service.StockOutbound{}
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

	err = s.PartOutbound(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("出库失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.Id, "操作成功")
}

// PrintOutbound 打印出库单
// @Summary 打印出库单
// @Description 打印出库单
// @Tags 出库单
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.StockOutboundPrintResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-outbound/print-outbound/{id} [get]
// @Security Bearer
func (e StockOutbound) PrintOutbound(c *gin.Context) {
	req := dto.StockOutboundPrintReq{}
	s := service.StockOutbound{}
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
	var object dto.StockOutboundPrintResp

	p := actions.GetPermissionFromContext(c)
	err = s.PrintOutbound(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("打印出库单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(object, "成功")
}

// PrintPicking 打印拣货单
// @Summary 打印拣货单
// @Description 打印拣货单
// @Tags 出库单
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.StockOutboundPrintResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-outbound/print-picking/{id} [get]
// @Security Bearer
func (e StockOutbound) PrintPicking(c *gin.Context) {
	req := dto.StockOutboundPrintReq{}
	s := service.StockOutbound{}
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
	var object dto.StockOutboundPrintResp

	p := actions.GetPermissionFromContext(c)
	err = s.PrintPicking(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("打印拣货单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(object, "成功")
}
