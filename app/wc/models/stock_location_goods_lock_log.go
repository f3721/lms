package models

import (
	"go-admin/common/models"
)

const (
	StockLocationGoodsLockLogFromType0 = "0"
	StockLocationGoodsLockLogFromType1 = "1"
	StockLocationGoodsLockLogFromType2 = "2"

	RemarkManualConfirmOutboundLocationProductsNotEqual = "手动出库-库位实际出库数量不等于库位应出库数量 解锁库位"
	RemarkLocationProductsLockCreateOutbound            = "出库单（调拨单or领用单）创建 锁定库位"
	RemarkLocationProductsLockCsOrderCancelAll          = "售后单取消全部商品 解锁库位"
	RemarkLocationProductsLockCsOrderCancelPart         = "售后单取消部分商品 解锁库位"
)

type StockLocationGoodsLockLog struct {
	models.Model

	StockLocationGoodsId int    `json:"stockLocationGoodsId" gorm:"type:int unsigned;"`
	AfterLockStock       int    `json:"afterLockStock" gorm:"type:int unsigned;comment:期末锁库数"`
	AfterStock           int    `json:"afterStock" gorm:"type:int unsigned;comment:期末在库数"`
	LockQuantity         int    `json:"lockQuantity" gorm:"type:int unsigned;comment:锁定库数"`
	DocketCode           string `json:"docketCode" gorm:"type:varchar(20);comment:单据编号"`
	FromType             string `json:"fromType" gorm:"type:tinyint;comment:0 调拨单 1 领用单 2 售后单"`
	Type                 string `json:"type" gorm:"type:tinyint(1);comment:0 库存加锁 1库存解锁"`
	Remark               string `json:"remark" gorm:"type:varchar(255);comment:remark"`
	GoodsId              int    `json:"goodsId" gorm:"type:int unsigned;comment:goods表ID"`
	LocationId           int    `json:"locationId" gorm:"type:int;"`
}

func (StockLocationGoodsLockLog) TableName() string {
	return "stock_location_goods_lock_log"
}

func (e *StockLocationGoodsLockLog) GetId() interface{} {
	return e.Id
}

func (e *StockLocationGoodsLockLog) SetLockType() {
	e.Type = "0"
}

func (e *StockLocationGoodsLockLog) SetUnLockType() {
	e.Type = "1"
}
