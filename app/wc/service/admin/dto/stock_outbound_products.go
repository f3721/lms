package dto

import (
	"go-admin/app/wc/models"
	common "go-admin/common/models"
)

type StockOutboundProductsReq struct {
	SkuCode  string `json:"skuCode" comment:"sku"`
	Quantity int    `json:"quantity" comment:"调拨数量"`
	GoodsId  int    `json:"goodsId"`
	HasSub   int    `json:"hasSub"`
	common.ControlBy
}

func (s *StockOutboundProductsReq) Generate(model *models.StockOutboundProducts) {
	model.SkuCode = s.SkuCode
	model.Quantity = s.Quantity
	model.GoodsId = s.GoodsId
	model.HasSub = s.HasSub
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *StockOutboundProductsReq) GenerateForUpdate(model *models.StockOutboundProducts) {

	model.SkuCode = s.SkuCode
	model.Quantity = s.Quantity
	model.GoodsId = s.GoodsId
	model.CreateBy = s.UpdateBy
	model.UpdateBy = s.UpdateBy
	model.CreateByName = s.UpdateByName
	model.UpdateByName = s.UpdateByName
}

type StockOutboundProductsConfirmReq struct {
	Id                  int    `json:"id"`                  // 入库单产品明细id
	SubId               int    `json:"subId"`               // 库位ID
	GoodsId             int    `json:"goodsId"`             // 商品ID
	SkuCode             string `json:"skuCode"`             // 商品SKU
	LocationActQuantity int    `json:"locationActQuantity"` // 实际出库数量
	StockLocationId     int    `json:"stockLocationId"`     // 出库库位ID
}

type StockOutboundProductsForOrderReq struct {
	SkuCode  string `json:"skuCode" comment:"sku"`
	Quantity int    `json:"quantity" comment:"数量"`
	GoodsId  int    `json:"goodsId"`
	VendorId int    `json:"vendorId" comment:"货主id"`
	common.ControlBy
}

func (s *StockOutboundProductsForOrderReq) Generate(model *models.StockOutboundProducts) {
	model.SkuCode = s.SkuCode
	model.Quantity = s.Quantity
	model.GoodsId = s.GoodsId
	model.HasSub = 1
	model.VendorId = s.VendorId
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

type StockOutboundProductsPartCancelForCsOrderReq struct {
	Quantity int `json:"quantity" comment:"取消数量"`
	GoodsId  int `json:"goodsId" comment:"商品ID"`
	common.ControlBy
}
