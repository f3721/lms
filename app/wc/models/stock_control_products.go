package models

import (
	"go-admin/common/models"
	"gorm.io/gorm"
)

type StockControlProducts struct {
	models.Model

	ControlCode        string        `json:"controlCode" gorm:"type:varchar(10);comment:调整单code"`
	SkuCode            string        `json:"skuCode" gorm:"type:varchar(10);comment:sku"`
	CurrentQuantity    int           `json:"currentQuantity" gorm:"type:int unsigned;comment:盘后数量"`
	Type               string        `json:"type" gorm:"type:tinyint;comment:调整类型: 0 盘盈 1 盘亏 2 无差异"`
	Quantity           int           `json:"quantity" gorm:"type:int unsigned;comment:差异数量"`
	VendorId           int           `json:"vendorId" gorm:"type:int unsigned;comment:货主id"`
	WarehouseCode      string        `json:"warehouseCode" gorm:"type:varchar(20);comment:实体仓code"`
	LogicWarehouseCode string        `json:"logicWarehouseCode" gorm:"type:varchar(20);comment:逻辑仓code"`
	GoodsId            int           `json:"goodsId" gorm:"type:int unsigned;comment:goods表id"`
	StockLocationId    int           `json:"stockLocationId" gorm:"type:int unsigned;comment:库位id"`
	StockLocation      StockLocation `json:"-" gorm:"foreignKey:StockLocationId;references:Id"`
	models.ModelTime
	models.ControlBy
}

func (StockControlProducts) TableName() string {
	return "stock_control_products"
}

func (e *StockControlProducts) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *StockControlProducts) GetId() interface{} {
	return e.Id
}

func (e *StockControlProducts) GetByControlCode(tx *gorm.DB, controlCode string) ([]StockControlProducts, error) {
	var rows []StockControlProducts
	if err := tx.Where("control_code = ?", controlCode).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

//根据调拨单code删除调拨单产品

func (e *StockControlProducts) DeleteByControlCode(tx *gorm.DB, ControlCode string) error {
	if err := tx.Where("control_code = ?", ControlCode).Delete(e).Error; err != nil {
		return err
	}
	return nil
}
