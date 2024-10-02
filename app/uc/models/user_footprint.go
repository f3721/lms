package models

import (
	"go-admin/common/models"
)

type UserFootprint struct {
	models.Model

	GoodsId       int    `json:"goodsId" gorm:"type:int unsigned;comment:GoodsId"`     //
	SkuCode       string `json:"skuCode" gorm:"type:varchar(20);comment:SkuCode"`      //
	UserId        int    `json:"userId" gorm:"type:int unsigned;comment:UserId"`       //
	WarehouseCode string `json:"warehouseCode" gorm:"type:varchar(20);comment:仓库code"` // 仓库code
	models.ModelTime
}

func (UserFootprint) TableName() string {
	return "user_footprint"
}

func (e *UserFootprint) Generate() *UserFootprint {
	o := *e
	return &o
}

func (e *UserFootprint) GetId() interface{} {
	return e.Id
}
