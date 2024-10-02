package admin

import (
	"encoding/json"
	"fmt"
	"go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/excel"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/common/actions"
)

type StockLocationGoodsLog struct {
	api.Api
}

// GetPage 获取库位商品变动日志列表
// @Summary 获取库位商品变动日志列表
// @Description 获取库位商品变动日志列表
// @Tags 库位商品变动日志
// @Param logicWarehouseCode query string false "逻辑仓code"
// @Param warehouseCode query string false "实体仓code"
// @Param queryWarehouseCode query string false "实体仓code多个用','分割"
// @Param skuCode query string false "sku多个用','分割"
// @Param productName query string false "productName"
// @Param productNo query string false "productNo"
// @Param vendorSkuCode query string false "vendorSkuCode"
// @Param vendorCode query string false "vendorCode"
// @Param vendorName query string false "vendorName"
// @Param vendorShortName query string false "vendorShortName"
// @Param isVirtualWarehouse query string false "isVirtualWarehouse"
// @Param createdAtStart query string false "createdAtStart"
// @Param createdAtEnd query string false "createdAtEnd"
// @Param fromType query string false "fromType 0 出库单 1 入库单 2 库存调整单"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.StockLocationGoodsLog}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-location-goods-log [get]
// @Security Bearer
func (e StockLocationGoodsLog) GetPage(c *gin.Context) {
	req := dto.StockLocationGoodsLogGetPageReq{}
	s := admin.StockLocationGoodsLog{}
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
	list := make([]dto.StockLocationGoodsLogGetPageResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取库位商品变动日志失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Export 库位库存日志导出
// @Summary 库位库存日志导出
// @Description 库位库存日志导出
// @Tags 库位库存日志
// @Param logicWarehouseCode query string false "逻辑仓code"
// @Param warehouseCode query string false "实体仓code"
// @Param queryWarehouseCode query string false "实体仓code多个用','分割"
// @Param skuCode query string false "sku多个用','分割"
// @Param productName query string false "productName"
// @Param productNo query string false "productNo"
// @Param vendorSkuCode query string false "vendorSkuCode"
// @Param vendorCode query string false "vendorCode"
// @Param vendorName query string false "vendorName"
// @Param vendorShortName query string false "vendorShortName"
// @Param isVirtualWarehouse query string false "isVirtualWarehouse"
// @Param createdAtStart query string false "createdAtStart"
// @Param createdAtEnd query string false "createdAtEnd"
// @Param fromType query string false "fromType 0 出库单 1 入库单 2 库存调整单"
// @Param docketCode query string false "docketCode"
// @Param ids query string false "ids 多个用,分割 其他查询无效"
// @Router /api/v1/wc/admin/stock-location-goods-log/export [get]
// @Security Bearer
func (e StockLocationGoodsLog) Export(c *gin.Context) {
	req := dto.StockLocationGoodsLogGetPageReq{}
	s := admin.StockLocationGoodsLog{}
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
	outData := &[]dto.StockLocationGoodsLogGetPageResp{}
	err = s.Export(&req, p, outData)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("库存变动记录导出失败，\r\n失败信息 %s", err.Error()))
		return
	}

	title := []map[string]string{
		{"id": "ID"},
		{"skuCode": "SKU"},
		{"productName": "商品名称"},
		{"mfgModel": "型号"},
		{"brandName": "品牌"},
		{"salesUom": "单位"},
		{"productNo": "物料编码"},
		{"warehouseName": "实体仓"},
		{"warehouseCode": "实体仓编码"},
		{"logicWarehouseName": "逻辑仓"},
		{"logicWarehouseCode": "逻辑仓编码"},
		{"locationCode": "库位"},
		{"vendorName": "货主"},
		{"vendorCode": "货主编码"},
		{"vendorShortName": "货主简称"},
		{"vendorSkuCode": "货主SKU"},
		{"changeStock": "变动数量"},
		{"beforeStock": "期初在库数"},
		{"afterStock": "期末在库数"},
		{"createdTime": "变动时间"},
		{"docketCode": "来源单据编号"},
		{"fromTypeName": "来源单据类型"},
	}
	exportData := []map[string]interface{}{}
	byteJson, _ := json.Marshal(outData)
	_ = json.Unmarshal(byteJson, &exportData)
	excelApp := excel.NewExcel()
	err = excelApp.ExportExcelByMap(e.Context, title, exportData, "库位库存变动记录-导出", "A1")
	if err != nil {
		e.Error(500, err, fmt.Sprintf("库位库存变动记录导出失败，\r\n失败信息 %s", err.Error()))
		return
	}
}
