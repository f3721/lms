package dto

import (
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	dtoStock "go-admin/common/dto/stock/dto"
	"time"
)

type StockLocationGoodsLogGetPageResp struct {
	models.StockLocationGoodsLog
	ProductName        string `json:"productName"`
	MfgModel           string `json:"mfgModel"`
	BrandName          string `json:"brandName"`
	SalesUom           string `json:"salesUom"`
	ProductNo          string `json:"productNo"`
	WarehouseCode      string `json:"warehouseCode"`
	LogicWarehouseCode string `json:"logicWarehouseCode"`
	SkuCode            string `json:"skuCode"`
	WarehouseName      string `json:"warehouseName"`
	LogicWarehouseName string `json:"logicWarehouseName"`
	LocationId         int    `json:"locationId"`
	LocationCode       string `json:"locationCode"`
	VendorSkuCode      string `json:"vendorSkuCode"`
	VendorCode         string `json:"vendorCode"`
	VendorName         string `json:"vendorName"`
	VendorShortName    string `json:"vendorShortName"`
	FromTypeName       string `json:"fromTypeName"`
	CreatedTime        string `json:"createdTime"`
}

type StockLocationGoodsLogGetPageReq struct {
	dto.Pagination `search:"-"`
	dtoStock.WarehousesSearch
	WarehouseCode      string    `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:stock_location" comment:"实体仓code"`
	LogicWarehouseCode string    `form:"logicWarehouseCode"  search:"type:exact;column:logic_warehouse_code;table:stock_location" comment:"逻辑仓code"`
	LocationCode       string    `form:"locationCode"  search:"type:exact;column:location_code;table:stock_location" comment:"库位编号"`
	VendorSkuCode      string    `form:"vendorSkuCode"  search:"type:exact;column:supplier_sku_code;table:product"`
	VendorCode         string    `form:"vendorCode"  search:"type:exact;column:code;table:vendors"`
	VendorName         string    `form:"vendorName"  search:"type:contains;column:name_zh;table:vendors"`
	VendorShortName    string    `form:"vendorShortName"  search:"type:contains;column:short_name;table:vendors"`
	IsVirtualWarehouse string    `form:"isVirtualWarehouse"  search:"type:exact;column:is_virtual;table:warehouse"`
	CreatedAtStart     time.Time `form:"createdAtStart"  search:"-" comment:"" time_format:"2006-01-02 15:04:05"`
	CreatedAtEnd       time.Time `form:"createdAtEnd"  search:"-" comment:"" time_format:"2006-01-02 15:04:05"`
	FromType           int       `form:"fromType"  search:"type:exact;column:from_type;table:stock_location_goods_log"`
	DocketCode         string    `form:"docketCode"  search:"type:contains;column:docket_code;table:stock_location_goods_log"`
	Ids                string    `form:"ids"  search:"-"`
	dtoStock.ProductSearch
}

func (m *StockLocationGoodsLogGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type StockLocationLogGetPageResp struct {
	models.StockLocationGoodsLog
	ProductName        string `json:"productName"`
	MfgModel           string `json:"mfgModel"`
	BrandName          string `json:"brandName"`
	SalesUom           string `json:"salesUom"`
	ProductNo          string `json:"productNo"`
	WarehouseName      string `json:"warehouseName"`
	LogicWarehouseName string `json:"logicWarehouseName"`
	VendorSkuCode      string `json:"vendorSkuCode"`
	VendorCode         string `json:"vendorCode"`
	VendorName         string `json:"vendorName"`
	VendorShortName    string `json:"vendorShortName"`
	FromTypeName       string `json:"fromTypeName"`
	CreatedTime        string `json:"createdTime"`
}

// StockLocationGoodsLogGetReq 功能获取请求参数
type StockLocationGoodsLogGetReq struct {
	Id int `uri:"id"`
}

func (s *StockLocationGoodsLogGetReq) GetId() interface{} {
	return s.Id
}

// StockLocationGoodsLogDeleteReq 功能删除请求参数
type StockLocationGoodsLogDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *StockLocationGoodsLogDeleteReq) GetId() interface{} {
	return s.Ids
}
