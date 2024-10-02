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

type StockEntry struct {
	api.Api
}

// GetPage 获取入库单列表
// @Summary 获取入库单列表
// @Description 获取入库单列表
// @Tags 入库单
// @Param entryCode query string false "入库单编码"
// @Param type query string false "入库类型:  0 大货入库  1 退货入库"
// @Param status query string false "状态:0-已作废 1-创建 2-已完成"
// @Param sourceCode query string false "来源单据code"
// @Param entryTime query time.Time false "入库时间"
// @Param warehouseCode query string false "实体仓code"
// @Param logicWarehouseCode query string false "逻辑仓code"
// @Param vendorId query int false "货主id"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.StockEntryGetPageResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-entry [get]
// @Security Bearer
func (e StockEntry) GetPage(c *gin.Context) {
	req := dto.StockEntryGetPageReq{}
	s := service.StockEntry{}
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
	list := make([]dto.StockEntryGetPageResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取入库单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取入库单
// @Summary 获取入库单
// @Description 获取入库单
// @Tags 入库单
// @Param id path int false "id"
// @Param type path string false "eg: confirm"
// @Success 200 {object} response.Response{data=dto.StockEntryGetResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-entry/{id} [get]
// @Security Bearer
func (e StockEntry) Get(c *gin.Context) {
	req := dto.StockEntryGetReq{}
	s := service.StockEntry{}
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
	var object dto.StockEntryGetResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取入库单失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Get 审核
// @Summary 审核
// @Description 审核
// @Tags 入库单
// @Param id path int false "id"
// @Param type path string false "eg: confirm"
// @Success 200 {object} response.Response{data=dto.StockEntryGetResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-entry/check-entry{id} [put]
// @Security Bearer
func (e StockEntry) CheckEntry(c *gin.Context) {
	req := dto.CheckStockEntryReq{}
	s := service.StockEntry{}
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
	//var object dto.StockEntryGetResp

	p := actions.GetPermissionFromContext(c)
	err = s.CheckStockEntry(c, &req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("入库单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "审核成功")
}

// Add 采购入库新建
// @Summary 采购入库新建
// @Description 采购入库新建
// @Tags 采购入库新建
// @Accept application/json
// @Product application/json
// @Param data body dto.AddStockEntryReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "入库成功"}"
// @Router /api/v1/wc/admin/stock-entry/add [post]
// @Security Bearer
func (e StockEntry) Add(c *gin.Context) {
	req := dto.AddStockEntryReq{}
	s := service.StockEntry{}
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

	req.SetCreateBy(user.GetUserId(c))
	req.SetCreateByName(user.GetUserName(c))
	p := actions.GetPermissionFromContext(c)
	err = s.AddStockEntry(e.Orm, &req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("入库失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(0, "入库成功")
}

// Add 采购入库编辑
// @Summary 采购入库编辑
// @Description 采购入库编辑
// @Tags 采购入库编辑
// @Accept application/json
// @Product application/json
// @Param data body dto.AddStockEntryReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "入库成功"}"
// @Router /api/v1/wc/admin/stock-entry/add [post]
// @Security Bearer
func (e StockEntry) Edit(c *gin.Context) {
	req := dto.AddStockEntryReq{}
	s := service.StockEntry{}
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

	err = s.EditStockEntry(e.Orm, &req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(0, "保存成功")
}

// ConfirmEntry 确认入库
// @Summary 确认入库
// @Description 确认入库
// @Tags 入库单
// @Accept application/json
// @Product application/json
// @Param data body dto.StockEntryConfirmReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "入库成功"}"
// @Router /api/v1/wc/admin/stock-entry/confirm-entry [post]
// @Security Bearer
func (e StockEntry) ConfirmEntry(c *gin.Context) {
	req := dto.StockEntryConfirmReq{}
	s := service.StockEntry{}
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

	err = s.ConfirmEntry(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("入库失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "入库成功")
}

// PartEntry 部分入库
// @Summary 部分入库
// @Description 部分入库
// @Tags 入库单
// @Accept application/json
// @Product application/json
// @Param data body dto.StockEntryPartReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "入库成功"}"
// @Router /api/v1/wc/admin/stock-entry/part-entry [post]
// @Security Bearer
func (e StockEntry) PartEntry(c *gin.Context) {
	req := dto.StockEntryPartReq{}
	s := service.StockEntry{}
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

	html, err := s.PartEntry(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("入库失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(html, "操作成功")
}

// Stash 暂存（保存）入库单
// @Summary 暂存（保存）入库单
// @Description 暂存（保存）入库单
// @Tags 入库单
// @Accept application/json
// @Product application/json
// @Param data body dto.StockEntryStashReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "保存成功"}"
// @Router /api/v1/wc/admin/stock-entry/stash [post]
// @Security Bearer
func (e StockEntry) Stash(c *gin.Context) {
	req := dto.StockEntryStashReq{}
	s := service.StockEntry{}
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

	err = s.Stash(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("保存入库单失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "保存成功")
}

// StockLocationsSelect 下拉获取库位列表
// @Summary 下拉获取库位列表
// @Description 下拉获取库位列表
// @Tags 入库单
// @Param logicWarehouseCode query string false "逻辑仓code"
// @Param isDefective query int false "是否次品"
// @Success 200 {object} response.Response{data=[]dto.StockEntryProductsLocationGetResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-entry/stock-locations [get]
// @Security Bearer
func (e StockEntry) StockLocationsSelect(c *gin.Context) {
	req := dto.StockEntryProductsLocationGetReq{}
	s := service.StockEntry{}
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
	list := make([]dto.StockEntryProductsLocationGetResp, 0)

	err = s.StockLocationsSelect(&req, p, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("下拉获取库位列表，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(list, "成功")
}

// GetPrintHtml 获取打印html
// @Summary 获取打印html
// @Description 获取打印html
// @Tags 入库单
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.StockEntryPrintHtmlResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-entry/print-html/{id} [get]
// @Security Bearer
func (e StockEntry) GetPrintHtml(c *gin.Context) {
	req := dto.StockEntryPrintHtmlReq{}
	s := service.StockEntry{}
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
	data := dto.StockEntryPrintHtmlResp{}

	err = s.GetPrintHtml(&req, p, &data)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取入库单打印html，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(data, "成功")
}

// GetPrintSkuInfo 获取打印sku信息
// @Summary 获取打印sku信息
// @Description 获取打印sku信息
// @Tags 入库单
// @Param goodsId path int false "goodsId"
// @Success 200 {object} response.Response{data=dto.StockEntryPrintSkuResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-entry/print-sku/{goodsId} [get]
// @Security Bearer
func (e StockEntry) GetPrintSkuInfo(c *gin.Context) {
	req := dto.StockEntryPrintSkuReq{}
	s := service.StockEntry{}
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
	data := dto.StockEntryPrintSkuResp{}

	err = s.GetPrintSkuInfo(&req, p, &data)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取打印sku信息，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(data, "成功")
}

// GetPrintSkuInfo 批量打印sku信息
// @Summary 批量打印sku信息
// @Description 批量打印sku信息
// @Tags 入库单
// @Param goodsId path int false "goodsId"
// @Success 200 {object} response.Response{data=dto.StockEntryPrintSkusReq} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-entry/print-sku/{goodsId} [get]
// @Security Bearer
func (e StockEntry) GetPrintSkuInfos(c *gin.Context) {
	req := dto.StockEntryPrintSkusReq{}
	s := service.StockEntry{}
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

	data, err := s.GetPrintSkuInfos(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取打印sku信息，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(data, "成功")
}

// GetSkuInfoByGoodsId 通过goodsId获取sku信息
// @Summary 通过goodsId获取sku信息
// @Description 通过goodsId获取sku信息
// @Tags 入库单
// @Param goodsId path int false "goodsId"
// @Success 200 {object} response.Response{data=dto.StockEntryScanSkuBaseInfoResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-entry/sku-info/{goodsId} [get]
// @Security Bearer
func (e StockEntry) GetSkuInfoByGoodsId(c *gin.Context) {
	req := dto.StockEntryPrintSkuReq{}
	s := service.StockEntry{}
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

	//p := actions.GetPermissionFromContext(c)
	data := dto.StockEntryScanSkuBaseInfoResp{}

	err = s.GetSkuInfoByGoodsId(&req, &data)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取打印sku信息，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(data, "成功")
}

func (e StockEntry) Export(c *gin.Context) {
	//req := dto.StockEntryGetPageReq{}
	//s := service.StockEntry{}

	req := dto.StockEntryGetPageReq{}
	s := service.StockEntry{}
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
	//list := make([]dto.StockEntryGetPageResp, 0)
	//var count int64

	outData, err := s.Export(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("产品维护导出失败： %s", err.Error()))
		return
	}
	title := []map[string]string{
		{"typeName": "入库类型"},
		{"statusName": "单据状态"},
		{"warehouseName": "实体仓"},
		{"logicWarehouseName": "逻辑仓"},
		{"vendorName": "货主"},
		{"entryCode": "单据编号"},
		{"sourceCode": "来源单据编号"},
		{"sourceTypeName": "来源单据类型"},
		{"remark": "备注"},
		{"skuCode": "SKU"},
		{"productName": "商品名称"},
		{"mfgModel": "型号"},
		{"brandName": "品牌"},
		{"salesUom": "单位"},
		{"productNo": "物料编码"},
		{"vendorSkuCode": "货主SKU"},
		{"Quantity": "应入库数量"},
		{"ActQuantity": "实际入库数量"},
		{"diffNum": "差异数量"},
		{"locationCode": "库位编号"},
		{"isDefective": "是否次品"},
		{"entryTime": "行入库时间"},
	}
	exportData := []map[string]interface{}{}
	byteJson, _ := json.Marshal(outData)
	_ = json.Unmarshal(byteJson, &exportData)
	excelApp := excel.NewExcel()
	err = excelApp.ExportExcelByMap(e.Context, title, exportData, "入库单-导出", "入库单明细")
	if err != nil {
		e.Error(500, err, fmt.Sprintf("入库单导出失败： %s", err.Error()))
	}
}

// ValidateSkus 入库单验证skus
// @Summary 入库单验证skus
// @Description 入库单验证skus
// @Tags 入库单
// @Param skuCodes query string false "多个用,分割"
// @Param vendorId query int false "货主id"
// @Param WarehouseCode query string false "实体仓code"
// @Param logicWarehouseCode query string false "逻辑仓code"
// @Success 200 {object} response.Response{data=[]dto.StockTransferValidateSkusResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-entry/validate-skus [get]
// @Security Bearer
func (e StockEntry) ValidateSkus(c *gin.Context) {
	req := dto.StockEntryferValidateSkusReq{}
	s := service.StockEntry{}
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
	var object []dto.StockTransferValidateSkusResp

	p := actions.GetPermissionFromContext(c)
	err = s.ValidateSkus(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf(err.Error()))
		return
	}
	e.OK(object, "查询成功")
}
