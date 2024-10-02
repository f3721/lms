package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/jinzhu/copier"
	"go-admin/common/excel"
	"go-admin/common/utils"

	service "go-admin/app/wc/service/admin"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type StockLocationGoods struct {
	api.Api
}

// GetPage 获取库位商品列表
// @Summary 获取库位商品列表
// @Description 获取库位商品列表
// @Tags 库位商品
// @Param locationCode query string false "库位编码"
// @Param skuCode query string false "SKU"
// @Param productName query string false "商品名称"
// @Param supplierSkuCode query string false "货主SKU"
// @Param productNo query string false "物料编码"
// @Param locationCode query string false "实体仓"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]dto.StockLocationGoodsListResp}} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/stock-location-goods [get]
// @Security Bearer
func (e StockLocationGoods) GetPage(c *gin.Context) {
	req := dto.StockLocationGoodsGetPageReq{}
	s := service.StockLocationGoods{}
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
	list := make([]dto.StockLocationGoodsListResp, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取库位商品失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取库位商品
// @Summary 获取库位商品
// @Description 获取库位商品
// @Tags 库位商品
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=dto.StockLocationGoodsListResp} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/stock-location-goods/{id} [get]
// @Security Bearer
func (e StockLocationGoods) Get(c *gin.Context) {
	req := dto.StockLocationGoodsGetReq{}
	s := service.StockLocationGoods{}
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
	var object dto.StockLocationGoodsListResp

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取库位商品失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// GetSameLogicStockLocationList 获取同一逻辑仓下的仓位
// @Summary 获取同一逻辑仓下的仓位
// @Description 获取同一逻辑仓下的仓位
// @Tags 库位商品
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=[]cDto.Option} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/stock-location-goods/location-list [get]
// @Security Bearer
func (e StockLocationGoods) GetSameLogicStockLocationList(c *gin.Context) {
	req := dto.SameLogicStockLocationReq{}
	s := service.StockLocationGoods{}
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
	list := make([]cDto.Option, 0)
	var count int64

	err = s.GetSameLogicStockLocationList(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取仓位列表失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(list, "查询成功")
}

// TransferStock 转移数量
// @Summary 转移数量
// @Description 转移数量
// @Tags 库位商品
// @Param data body dto.TransferStockReq true "data"
// @Success 200 {object} response.Response{data=[]cDto.Option} "{"code": 200, "data": [...]}"
// @Router /api/v1/wc/stock-location-goods/transfer-stock [post]
// @Security Bearer
func (e StockLocationGoods) TransferStock(c *gin.Context) {
	req := dto.TransferStockReq{}
	s := service.StockLocationGoods{}
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

	err = s.TransferStock(&req, p, c)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("转移失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(nil, "转移成功")
}

// Export 导出库位商品
// @Summary 导出库位商品
// @Description 导出库位商品
// @Tags 库位商品
// @Param filterIds query string false "勾选的id集"
// @Param locationCode query string false "库位编码"
// @Param warehouseCode query string false "实体仓"
// @Router /api/v1/wc/stock-location-goods/export [get]
// @Security Bearer
func (e StockLocationGoods) Export(c *gin.Context) {
	req := dto.StockLocationGoodsGetPageReq{}
	s := service.StockLocationGoods{}
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
	list := make([]dto.StockLocationGoodsListResp, 0)
	var count int64

	req.PageSize = -1
	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取库位商品失败，\r\n失败信息 %s", err.Error()))
		return
	}

	var exportData []interface{}

	for _, resp := range list {
		var tmp dto.StockLocationGoodsExportListResp
		_ = copier.Copy(&tmp, &resp)
		tmp.UpdatedAt = utils.TimeFormat(resp.UpdatedAt)
		exportData = append(exportData, tmp)
	}

	title := []map[string]string{
		{"skuCode": "SKU"},
		{"productName": "商品名称"},
		{"locationCode": "库位编码"},
		{"warehouseName": "实体仓"},
		{"totalStock": "在库数量"},
		{"brandName": "品牌"},
		{"mfgModel": "型号"},
		{"supplierSkuCode": "货主SKU"},
		{"vendorName": "货主"},
		{"productNo": "物料编码"},
		{"updatedAt": "修改时间"},
	}
	excelApp := excel.NewExcel()
	err = excelApp.ExportExcelByStruct(c, title, exportData, "库位商品管理-导出", "Sheet1")
	if err != nil {
		e.Error(500, err, fmt.Sprintf("库位商品管理导出失败： %s", err.Error()))
	}
}
