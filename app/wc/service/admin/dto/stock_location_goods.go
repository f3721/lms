package dto

import (
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type StockLocationGoodsGetPageReq struct {
	dto.Pagination     `search:"-"`
	LocationCode 	string `form:"locationCode"  search:"type:exact;column:location_code;table:sl"`
	WarehouseCode 	string `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:sl"`
	SkuCode         string `form:"skuCode"  search:"type:exact;column:sku_code;table:g"`
	ProductName     string `form:"productName"  search:"type:contains;column:name_zh;table:p"`
	SupplierSkuCode string `form:"supplierSkuCode"  search:"type:exact;column:supplier_sku_code;table:g"`
	ProductNo       string `form:"productNo"  search:"type:exact;column:product_no;table:g"`
	Ids 			[]int `form:"ids[]"  search:"-"`
    StockLocationGoodsOrder
}

type StockLocationGoodsOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:stock_location_goods"`
    
}

func (m *StockLocationGoodsGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type StockLocationGoodsInsertReq struct {
    Id int `json:"-" comment:""` // 
    GoodsId int `json:"goodsId" comment:"goods表ID"`
    LocationId int `json:"locationId" comment:"库位表ID"`
    Stock int `json:"stock" comment:"库位商品数量"`
    LockStock int `json:"lockStock" comment:"锁定库存"`
    common.ControlBy
}

func (s *StockLocationGoodsInsertReq) Generate(model *models.StockLocationGoods)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.GoodsId = s.GoodsId
    model.LocationId = s.LocationId
    model.Stock = s.Stock
    model.LockStock = s.LockStock
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *StockLocationGoodsInsertReq) GetId() interface{} {
	return s.Id
}

type StockLocationGoodsUpdateReq struct {
    Id int `uri:"id" comment:""` // 
    GoodsId int `json:"goodsId" comment:"goods表ID"`
    LocationId int `json:"locationId" comment:"库位表ID"`
    Stock int `json:"stock" comment:"库位商品数量"`
    LockStock int `json:"lockStock" comment:"锁定库存"`
    common.ControlBy
}

func (s *StockLocationGoodsUpdateReq) Generate(model *models.StockLocationGoods)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.GoodsId = s.GoodsId
    model.LocationId = s.LocationId
    model.Stock = s.Stock
    model.LockStock = s.LockStock
    model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *StockLocationGoodsUpdateReq) GetId() interface{} {
	return s.Id
}

// StockLocationGoodsGetReq 功能获取请求参数
type StockLocationGoodsGetReq struct {
     Id int `uri:"id"`
}
func (s *StockLocationGoodsGetReq) GetId() interface{} {
	return s.Id
}

// StockLocationGoodsDeleteReq 功能删除请求参数
type StockLocationGoodsDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *StockLocationGoodsDeleteReq) GetId() interface{} {
	return s.Ids
}

type StockLocationGoodsListResp struct {
	models.StockLocationGoods
	TotalStock int `json:"totalStock" comment:"在库数量"`
	LocationCode string `json:"locationCode" gorm:"type:varchar(20);comment:库位编码"`
	WarehouseName string `json:"warehouseName" gorm:"type:varchar(20);comment:仓库名称"`
	SkuCode           string  `json:"skuCode" gorm:"type:varchar(10);comment:商品SKU"`
	SupplierSkuCode   string  `json:"supplierSkuCode" gorm:"type:varchar(20);comment:货主SKU"`
	ProductNo         string  `json:"productNo" gorm:"type:varchar(20);comment:物料编码"`
	VendorName        string `json:"vendorName" gorm:"type:varchar(50);comment:货主名"`
	ProductName      string  `json:"productName" gorm:"type:varchar(512);comment:商品名"`
	MfgModel          string  `json:"mfgModel" gorm:"type:varchar(255);comment:制造厂型号"`
	BrandName         string `json:"brandName" gorm:"type:varchar(100);comment:品牌"`
}

type StockLocationGoodsExportListResp struct {
	SkuCode string `json:"skuCode" gorm:"type:varchar(10);comment:商品SKU"`
	ProductName string `json:"productName" gorm:"type:varchar(512);comment:商品名"`
	LocationCode string `json:"locationCode" gorm:"type:varchar(20);comment:库位编码"`
	WarehouseName string `json:"warehouseName" gorm:"type:varchar(20);comment:仓库名称"`
	TotalStock int `json:"totalStock" comment:"在库数量"`
	BrandName string `json:"brandName" gorm:"type:varchar(100);comment:品牌"`
	MfgModel   string  `json:"mfgModel" gorm:"type:varchar(255);comment:制造厂型号"`
	SupplierSkuCode string  `json:"supplierSkuCode" gorm:"type:varchar(20);comment:货主SKU"`
	VendorName string `json:"vendorName" gorm:"type:varchar(50);comment:货主名"`
	ProductNo string  `json:"productNo" gorm:"type:varchar(20);comment:物料编码"`
	UpdatedAt string `json:"updatedAt" gorm:"comment:最后更新时间"`
}

type SameLogicStockLocationReq struct {
	LocationCode string `form:"locationCode"  search:"type:exact;column:location_code;table:sl"`
	ExceptSelf string `form:"exceptSelf"`
}
type TransferStockReq struct {
	Id int `json:"id" comment:"原库位商品id" vd:"$>0"`
	GoodsId int `json:"goodsId" comment:"goods表ID" vd:"$>0"`
	LocationId int `json:"locationId" comment:"原库位ID" vd:"$>0"`
	TransferStock int `json:"transferStock" comment:"转移数量" vd:"$>0"`
	TargetLocationId int `json:"targetLocationId" comment:"目标库位ID" vd:"$>0"`
}
