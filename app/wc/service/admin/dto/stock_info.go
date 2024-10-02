package dto

import (
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	dtoStock "go-admin/common/dto/stock/dto"
)

type StockInfoGetPageReq struct {
	dto.Pagination `search:"-"`
	dtoStock.WarehousesSearch
	WarehouseCode      string `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:stock_info" comment:"实体仓code"`
	LogicWarehouseCode string `form:"logicWarehouseCode"  search:"type:exact;column:logic_warehouse_code;table:stock_info" comment:"逻辑仓code"`
	VendorSkuCode      string `form:"vendorSkuCode"  search:"type:exact;column:supplier_sku_code;table:goods"`
	VendorCode         string `form:"vendorCode"  search:"type:exact;column:code;table:vendors"`
	VendorName         string `form:"vendorName"  search:"type:contains;column:name_zh;table:vendors"`
	VendorShortName    string `form:"vendorShortName"  search:"type:contains;column:short_name;table:vendors"`
	IsVirtualWarehouse string `form:"isVirtualWarehouse"  search:"type:exact;column:is_virtual;table:warehouse"`
	Ids                string `form:"ids"  search:"-"`
	dtoStock.ProductSearch
}

func (m *StockInfoGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type OrderQuantity struct {
	Id                  int `json:"id"`
	OrderQuantity       int `json:"orderQuantity"`       // 订单购买数量
	OrderCancelQuantity int `json:"orderCancelQuantity"` // 订单取消数量
	OrderLockStock      int `json:"orderLockStock"`      // 订单占位库存
}

type StockInfoGetPageResp struct {
	models.StockInfo
	ProductName         string          `json:"productName"`
	MfgModel            string          `json:"mfgModel"`
	BrandName           string          `json:"brandName"`
	SalesUom            string          `json:"salesUom"`
	ProductNo           string          `json:"productNo"`
	WarehouseName       string          `json:"warehouseName"`
	LogicWarehouseName  string          `json:"logicWarehouseName"`
	VendorSkuCode       string          `json:"vendorSkuCode"`
	VendorCode          string          `json:"vendorCode"`
	VendorName          string          `json:"vendorName"`
	VendorShortName     string          `json:"vendorShortName"`
	TotalStock          int             `json:"totalStock"`
	LackStock           int             `json:"lackStock"`           // 实缺库存
	LocationStocks      []LocationStock `json:"locationStocks"`      // 可用库存数据
	LocationLockStocks  []LocationStock `json:"locationLockStocks"`  // 占用库存数据
	LocationTotalStocks []LocationStock `json:"locationTotalStocks"` // 在库库存数据
}

type LocationStock struct {
	StockInfoId        int    `json:"StockInfoId"`
	WarehouseCode      string `json:"warehouseCode"`
	LogicWarehouseCode string `json:"logicWarehouseCode"`
	LocationCode       string `json:"locationCode"`
	Stock              int    `json:"stock"`
	LockStock          int    `json:"lockStock"`
	TotalStock         int    `json:"totalStock"`
}

type InnerStockInfoGetByGoodsIdAndLwhCodeReq struct {
	Query []InnerStockInfoGetByGoodsIdAndLwhCodeReqInfo `json:"query"`
}

type InnerStockInfoGetByGoodsIdAndLwhCodeReqInfo struct {
	LogicWarehouseCode string `json:"logicWarehouseCode"  comment:"逻辑仓code"`
	GoodsId            int    `json:"goodsId" comment:"GoodsId"`
}

type InnerStockInfoGetByGoodsIdAndWarehouseCodeReq struct {
	WarehouseCode string `json:"warehouseCode"  search:"type:exact;column:warehouse_code;table:stock_info" comment:"实体仓code" vd:"@:len($)>0; msg:'WarehouseCode不能为空'"`
	GoodsIds      []int  `json:"goodsIds" search:"type:in;column:goods_id;table:stock_info" comment:"GoodsIds" vd:"@:len($)>0; msg:'GoodsIds不能为空'"`
}

func (m *InnerStockInfoGetByGoodsIdAndWarehouseCodeReq) GetNeedSearch() interface{} {
	return *m
}
