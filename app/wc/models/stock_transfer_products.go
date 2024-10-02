package models

import (
	"go-admin/common/models"

	"gorm.io/gorm"
)

type StockTransferProducts struct {
	models.Model

	TransferCode         string `json:"transferCode" gorm:"type:varchar(20);comment:调拨单code"`
	SkuCode              string `json:"skuCode" gorm:"type:varchar(10);comment:sku"`
	Quantity             int    `json:"quantity" gorm:"type:int unsigned;comment:调拨数量"`
	VendorId             int    `json:"vendorId" gorm:"type:int unsigned;comment:货主id"`
	WarehouseCode        string `json:"warehouseCode" gorm:"type:varchar(20);comment:实体仓code"`
	LogicWarehouseCode   string `json:"logicWarehouseCode" gorm:"type:varchar(20);comment:逻辑仓code"`
	GoodsId              int    `json:"goodsId" gorm:"type:int unsigned;comment:goods表id"`
	ToWarehouseCode      string `json:"toWarehouseCode" gorm:"type:varchar(20);comment:入库实体仓code"`
	ToLogicWarehouseCode string `json:"toLogicWarehouseCode" gorm:"type:varchar(20);comment:入库逻辑仓code"`
	ToGoodsId            int    `json:"toGoodsId" gorm:"type:int unsigned;comment:入库goods表id"`
	models.ModelTime
	models.ControlBy
}

func (StockTransferProducts) TableName() string {
	return "stock_transfer_products"
}

func (e *StockTransferProducts) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *StockTransferProducts) GetId() interface{} {
	return e.Id
}

//根据调拨单code删除调拨单产品

func (e *StockTransferProducts) DeleteByTransferCode(tx *gorm.DB, transferCode string) error {
	if err := tx.Where("transfer_code = ?", transferCode).Delete(e).Error; err != nil {
		return err
	}
	return nil
}

//根据调拨单code获取调拨单产品

func (e *StockTransferProducts) GetByTransferCode(tx *gorm.DB, transferCode string) ([]StockTransferProducts, error) {
	var rows []StockTransferProducts
	if err := tx.Where("transfer_code = ?", transferCode).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
