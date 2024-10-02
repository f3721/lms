package admin

import (
	"encoding/json"
	"fmt"
	"go-admin/app/wc/models"
	"go-admin/common/excel"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
)

type StockInfo struct {
	api.Api
}

// GetPage 获取库存信息列表
// @Summary 获取库存信息列表
// @Description 获取库存信息列表
// @Tags 库存信息
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
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.StockInfoGetPageResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/admin/stock-info [get]
// @Security Bearer
func (e StockInfo) GetPage(c *gin.Context) {
	req := dto.StockInfoGetPageReq{}
	s := service.StockInfo{}
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
	list := make([]dto.StockInfoGetPageResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取库存信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Export 库存信息导出
// @Summary 库存信息导出
// @Description 库存信息导出
// @Tags 库存信息
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
// @Param ids query string false "ids 多个用,分割 其他查询无效"
// @Router /api/v1/wc/admin/stock-info/export [get]
// @Security Bearer
func (e StockInfo) Export(c *gin.Context) {
	req := dto.StockInfoGetPageReq{}
	s := service.StockInfo{}
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
	outData := &[]dto.StockInfoGetPageResp{}
	err = s.Export(&req, p, outData)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("库存信息导出失败，\r\n失败信息 %s", err.Error()))
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
		{"vendorName": "货主"},
		{"vendorCode": "货主编码"},
		{"vendorShortName": "货主简称"},
		{"vendorSkuCode": "货主SKU"},
		{"stock": "可用库存"},
		{"lockStock": "占用库存"},
		{"totalStock": "在库库存"},
		{"lackStock": "实缺库存"},
	}
	exportData := []map[string]interface{}{}
	byteJson, _ := json.Marshal(outData)
	_ = json.Unmarshal(byteJson, &exportData)
	excelApp := excel.NewExcel()
	err = excelApp.ExportExcelByMap(e.Context, title, exportData, "库存明细-导出", "A1")
	if err != nil {
		e.Error(500, err, fmt.Sprintf("库存信息导出失败，\r\n失败信息 %s", err.Error()))
		return
	}
}

// InnerGetByGoodsIdAndLwhCOde Inner批量获取库存信息
// @Summary Inner批量获取库存信息
// @Description Inner批量获取库存信息
// @Tags Inner库存信息
// @Accept application/json
// @Product application/json
// @Param data body dto.InnerStockInfoGetByGoodsIdAndLwhCodeReq true "body"
// @Success 200 {object} response.Response{data=[]models.Warehouse} "{"code": 200, "data": [...]}"
// @Router /inner/wc/admin/stock-info/get [post]
func (e StockInfo) InnerGetByGoodsIdAndLwhCode(c *gin.Context) {
	req := dto.InnerStockInfoGetByGoodsIdAndLwhCodeReq{}
	s := service.StockInfo{}
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

	list := make([]models.StockInfo, 0)

	err = s.InnerGetByGoodsIdAndLwhCode(&req, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("Inner批量获取库存信息失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(list, "成功")
}

// InnerGetByGoodsIdAndWarehouseCode Inner批量获取库存信息通过实体仓和goodsId
// @Summary Inner批量获取库存信息通过实体仓和goodsId
// @Description Inner批量获取库存信息通过实体仓和goodsId
// @Tags Inner库存信息
// @Accept application/json
// @Product application/json
// @Param data body dto.InnerStockInfoGetByGoodsIdAndWarehouseCodeReq true "body"
// @Success 200 {object} response.Response{data=[]models.Warehouse} "{"code": 200, "data": [...]}"
// @Router /inner/wc/admin/stock-info [post]
func (e StockInfo) InnerGetByGoodsIdAndWarehouseCode(c *gin.Context) {
	req := dto.InnerStockInfoGetByGoodsIdAndWarehouseCodeReq{}
	s := service.StockInfo{}
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

	list := make([]models.StockInfo, 0)

	err = s.InnerGetByGoodsIdAndWarehouseCode(&req, &list)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("Inner批量获取库存信息通过实体仓和goodsId失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(list, "成功")
}
