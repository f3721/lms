package dto

import (
	"go-admin/app/wc/models"
	common "go-admin/common/models"
)

type StockTransferProductsReq struct {
	Id       int    `json:"-" comment:"id"` // id
	SkuCode  string `json:"skuCode" comment:"sku"`
	Quantity int    `json:"quantity" comment:"调拨数量"`
	common.ControlBy
}

func (s *StockTransferProductsReq) Generate(model *models.StockTransferProducts) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.SkuCode = s.SkuCode
	model.Quantity = s.Quantity
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *StockTransferProductsReq) GenerateForUpdate(model *models.StockTransferProducts) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.SkuCode = s.SkuCode
	model.Quantity = s.Quantity
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName
}
