package dto

import (
	"go-admin/app/wc/models"
	common "go-admin/common/models"
)

type StockEntryProductsReq struct {
	SkuCode     string `json:"skuCode" comment:"sku"`
	Quantity    int    `json:"quantity" comment:"调拨数量"`
	GoodsId     int    `json:"goodsId"`
	IsDefective int    `json:"isDefective"`
	common.ControlBy
}

type AddStockEntryProductsReq struct {
	Id       int    `json:"id"`
	SkuCode  string `json:"skuCode" comment:"sku"`
	Quantity int    `json:"quantity" comment:"入库数量"`
	//ActQuantity int    `json:"actQuantity" comment:"实际入库数量"`
	GoodsId int `json:"goodsId"`
	//IsDefective int `json:"isDefective"`
	//VendorId int `json:"vendorId"`
	//WarehouseCode      string `json:"warehouseCode"`
	//LogicWarehouseCode string `json:"logicWarehouseCode"`
	StockLocationId int `json:"stockLocationId"`

	StockEntryProductsSub []StockEntryProductsSubReq `json:"stockEntryProductsSub"`

	common.ControlBy
}

func (s *StockEntryProductsReq) Generate(model *models.StockEntryProducts) {
	model.SkuCode = s.SkuCode
	model.Quantity = s.Quantity
	model.GoodsId = s.GoodsId
	model.IsDefective = s.IsDefective
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *AddStockEntryProductsReq) Generate(model *models.StockEntryProducts) {
	model.SkuCode = s.SkuCode
	model.Quantity = s.Quantity
	model.GoodsId = s.GoodsId
	//model.IsDefective = s.IsDefective
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *StockEntryProductsReq) GenerateForUpdate(model *models.StockEntryProducts) {

	model.SkuCode = s.SkuCode
	model.Quantity = s.Quantity
	model.GoodsId = s.GoodsId
	model.IsDefective = s.IsDefective
	model.CreateBy = s.UpdateBy
	model.UpdateBy = s.UpdateBy
	model.CreateByName = s.UpdateByName
	model.UpdateByName = s.UpdateByName
}

type StockEntryProductsConfirmReq struct {
	SkuCode         string `json:"skuCode"`
	ActQuantity     int    `json:"actQuantity"`
	StockLocationId int    `json:"stockLocationId"`
}
