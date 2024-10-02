package dto

import (
	"go-admin/app/wc/models"
	common "go-admin/common/models"
)

type StockControlProductsReq struct {
	SkuCode            string `json:"skuCode" comment:"sku"`
	VendorId           int    `json:"vendorId" comment:"售后单所属货主id"`
	WarehouseCode      string `json:"warehouseCode" comment:"发货仓"`
	LogicWarehouseCode string `json:"logicWarehouseCode" comment:"逻辑仓"`
	GoodsId            int    `json:"goodsId" comment:"GoodsId"`
	Quantity           int    `json:"quantity" comment:"调整数量"`
	CurrentQuantity    int    `json:"currentQuantity" comment:"盘后数量"`
	Type               string `json:"type" comment:"调整类型: 0 盘盈 1 盘亏 2 无差异"`
	StockLocationId    int    `json:"stockLocationId"`
	common.ControlBy
}

func (s *StockControlProductsReq) Generate(model *models.StockControlProducts) {
	model.SkuCode = s.SkuCode
	model.VendorId = s.VendorId
	model.WarehouseCode = s.WarehouseCode
	model.LogicWarehouseCode = s.LogicWarehouseCode
	model.GoodsId = s.GoodsId
	model.Type = s.Type
	model.CurrentQuantity = s.CurrentQuantity
	model.Quantity = s.Quantity
	model.StockLocationId = s.StockLocationId
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *StockControlProductsReq) GenerateForUpdate(model *models.StockControlProducts) {
	model.SkuCode = s.SkuCode
	model.VendorId = s.VendorId
	model.WarehouseCode = s.WarehouseCode
	model.LogicWarehouseCode = s.LogicWarehouseCode
	model.GoodsId = s.GoodsId
	model.Type = s.Type
	model.CurrentQuantity = s.CurrentQuantity
	model.Quantity = s.Quantity
	model.StockLocationId = s.StockLocationId
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName
}

type StockControlProductsConfirmReq struct {
	SkuCode         string `json:"skuCode"`
	ActQuantity     int    `json:"actQuantity"`
	StockLocationId int    `json:"stockLocationId"`
}
