package models

import (
	"go-admin/common/models"
)

const (
	StockLockLogFromType0 = "0"
	StockLockLogFromType1 = "1"
	StockLockLogFromType2 = "2"

	RemarkManualConfirmOutboundProductsNotEqual = "手动出库-实际出库数量不等于应出库数量 解锁"
	RemarkStockLockLogTransferAuditReject       = "调拨单审批驳回 解锁"
	RemarkStockLockLogTransferCommit            = "调拨单提交 加锁"
	RemarkStockLockLogOrderAboutLock            = "领用单相关 加锁"
	RemarkStockLockLogOrderAboutUnLock          = "领用单相关 解锁"
	RemarkStockLockLogCsOrderCancelAll          = "售后单取消全部商品 解锁"
	RemarkStockLockLogCsOrderCancelPart         = "售后单取消部分商品 解锁"
)

type StockLockLog struct {
	models.Model

	StockInfoId        int    `json:"stockInfoId" gorm:"type:int unsigned;comment:stock_info表id"`
	AfterLockStock     int    `json:"afterLockStock" gorm:"type:int unsigned;comment:期末锁库数"`
	AfterStock         int    `json:"afterStock" gorm:"type:int unsigned;comment:期末在库数"`
	LockQuantity       int    `json:"lockQuantity" gorm:"type:int unsigned;comment:锁定库数"`
	DocketCode         string `json:"docketCode" gorm:"type:varchar(20);comment:单据编号"`
	FromType           string `json:"fromType" gorm:"type:tinyint;comment:0 调拨单 1 领用单 2 售后单"`
	Type               string `json:"type" gorm:"type:tinyint(1);comment:0 库存加锁 1库存解锁"`
	Remark             string `json:"remark" gorm:"type:varchar(255);comment:remark"`
	GoodsId            int    `json:"goodsId" gorm:"type:int unsigned;comment:goods表ID"`
	LogicWarehouseCode string `json:"logicWarehouseCode" gorm:"type:varchar(20);comment:逻辑仓code"`
}

func (StockLockLog) TableName() string {
	return "stock_lock_log"
}

func (e *StockLockLog) GetId() interface{} {
	return e.Id
}

func (e *StockLockLog) SetUnLockType() {
	e.Type = "1"
}

func (e *StockLockLog) SetLockType() {
	e.Type = "0"
}
