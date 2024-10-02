package dto

import (
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	dtoStock "go-admin/common/dto/stock/dto"
	common "go-admin/common/models"
	"time"
)

type StockLogGetPageReq struct {
	dto.Pagination `search:"-"`
	dtoStock.WarehousesSearch
	WarehouseCode      string    `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:stock_log" comment:"实体仓code"`
	LogicWarehouseCode string    `form:"logicWarehouseCode"  search:"type:exact;column:logic_warehouse_code;table:stock_log" comment:"逻辑仓code"`
	VendorSkuCode      string    `form:"vendorSkuCode"  search:"type:exact;column:supplier_sku_code;table:product"`
	VendorCode         string    `form:"vendorCode"  search:"type:exact;column:code;table:vendors"`
	VendorName         string    `form:"vendorName"  search:"type:contains;column:name_zh;table:vendors"`
	VendorShortName    string    `form:"vendorShortName"  search:"type:contains;column:short_name;table:vendors"`
	IsVirtualWarehouse string    `form:"isVirtualWarehouse"  search:"type:exact;column:is_virtual;table:warehouse"`
	CreatedAtStart     time.Time `form:"createdAtStart"  search:"-" comment:"" time_format:"2006-01-02 15:04:05"`
	CreatedAtEnd       time.Time `form:"createdAtEnd"  search:"-" comment:"" time_format:"2006-01-02 15:04:05"`
	FromType           string    `form:"fromType"  search:"type:exact;column:from_type;table:stock_log"`
	DocketCode         string    `form:"docketCode"  search:"type:contains;column:docket_code;table:stock_log"`
	Ids                string    `form:"ids"  search:"-"`
	dtoStock.ProductSearch
}

func (m *StockLogGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type StockLogGetPageResp struct {
	models.StockLog
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

type StockLogInsertReq struct {
	Id                 int    `json:"-" comment:""` //
	StockInfoId        int    `json:"stockInfoId" comment:"stock_info表id"`
	WarehouseCode      string `json:"warehouseCode" comment:"实体仓code"`
	LogicWarehouseCode string `json:"logicWarehouseCode" comment:"逻辑仓code"`
	SkuCode            string `json:"skuCode" comment:"sku code"`
	BeforeStock        int    `json:"beforeStock" comment:"期初在库数"`
	AfterStock         int    `json:"afterStock" comment:"期末在库数"`
	ChangeStock        int    `json:"changeStock" comment:"期末在库数"`
	DocketCode         string `json:"docketCode" comment:"单据编号"`
	FromType           string `json:"fromType" comment:"0 出库单 1 入库单 2 库存调整单"`
	VendorId           int    `json:"vendorId" comment:"货主id"`
	common.ControlBy
}

func (s *StockLogInsertReq) Generate(model *models.StockLog) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StockInfoId = s.StockInfoId
	model.WarehouseCode = s.WarehouseCode
	model.LogicWarehouseCode = s.LogicWarehouseCode
	model.SkuCode = s.SkuCode
	model.BeforeStock = s.BeforeStock
	model.AfterStock = s.AfterStock
	model.ChangeStock = s.ChangeStock
	model.DocketCode = s.DocketCode
	model.FromType = s.FromType
	model.VendorId = s.VendorId
}

func (s *StockLogInsertReq) GetId() interface{} {
	return s.Id
}

type StockLogUpdateReq struct {
	Id                 int    `uri:"id" comment:""` //
	StockInfoId        int    `json:"stockInfoId" comment:"stock_info表id"`
	WarehouseCode      string `json:"warehouseCode" comment:"实体仓code"`
	LogicWarehouseCode string `json:"logicWarehouseCode" comment:"逻辑仓code"`
	SkuCode            string `json:"skuCode" comment:"sku code"`
	BeforeStock        int    `json:"beforeStock" comment:"期初在库数"`
	AfterStock         int    `json:"afterStock" comment:"期末在库数"`
	ChangeStock        int    `json:"changeStock" comment:"期末在库数"`
	DocketCode         string `json:"docketCode" comment:"单据编号"`
	FromType           string `json:"fromType" comment:"0 出库单 1 入库单 2 库存调整单"`
	VendorId           int    `json:"vendorId" comment:"货主id"`
	common.ControlBy
}

func (s *StockLogUpdateReq) Generate(model *models.StockLog) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.StockInfoId = s.StockInfoId
	model.WarehouseCode = s.WarehouseCode
	model.LogicWarehouseCode = s.LogicWarehouseCode
	model.SkuCode = s.SkuCode
	model.BeforeStock = s.BeforeStock
	model.AfterStock = s.AfterStock
	model.ChangeStock = s.ChangeStock
	model.DocketCode = s.DocketCode
	model.FromType = s.FromType
	model.VendorId = s.VendorId
}

func (s *StockLogUpdateReq) GetId() interface{} {
	return s.Id
}

// StockLogGetReq 功能获取请求参数
type StockLogGetReq struct {
	Id int `uri:"id"`
}

func (s *StockLogGetReq) GetId() interface{} {
	return s.Id
}

// StockLogDeleteReq 功能删除请求参数
type StockLogDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *StockLogDeleteReq) GetId() interface{} {
	return s.Ids
}
