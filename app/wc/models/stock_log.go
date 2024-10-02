package models

import (
	"go-admin/common/models"
	"time"
)

const (
	StockLogFromType0 = "0"
	StockLogFromType1 = "1"
	StockLogFromType2 = "2"

	StockLocationGoodsLogFromType0 = "0"
	StockLocationGoodsLogFromType1 = "1"
	StockLocationGoodsLogFromType2 = "2"
	StockLocationGoodsLogFromType3 = "3"
)

var (
	StockLogFromTypeMap = map[string]string{
		StockLogFromType0: "出库单",
		StockLogFromType1: "入库单",
		StockLogFromType2: "库存调整单",
	}

	StockLocationGoodsLogFromTypeMap = map[string]string{
		StockLocationGoodsLogFromType0: "出库单",
		StockLocationGoodsLogFromType1: "入库单",
		StockLocationGoodsLogFromType2: "库存调整单",
		StockLocationGoodsLogFromType3: "库位转移",
	}
)

type StockLog struct {
	models.Model

	StockInfoId             int            `json:"stockInfoId" gorm:"type:int unsigned;comment:stock_info表id"`
	WarehouseCode           string         `json:"warehouseCode" gorm:"type:varchar(20);comment:实体仓code"`
	LogicWarehouseCode      string         `json:"logicWarehouseCode" gorm:"type:varchar(20);comment:逻辑仓code"`
	SkuCode                 string         `json:"skuCode" gorm:"type:varchar(10);comment:sku code"`
	BeforeStock             int            `json:"beforeStock" gorm:"type:int unsigned;comment:期初在库数"`
	AfterStock              int            `json:"afterStock" gorm:"type:int unsigned;comment:期末在库数"`
	ChangeStock             int            `json:"changeStock" gorm:"type:int unsigned;comment:变动数量"`
	DocketCode              string         `json:"docketCode" gorm:"type:varchar(20);comment:单据编号"`
	FromType                string         `json:"fromType" gorm:"type:smallint;comment:0 出库单 1 入库单 2 库存调整单"`
	VendorId                int            `json:"vendorId" gorm:"type:int unsigned;comment:货主id"`
	SnapshotAfterStock      int            `json:"snapshotAfterStock" gorm:"type:int(11);comment:"`
	SnapshotAfterLockStock  int            `json:"snapshotAfterLockStock" gorm:"type:int(11);comment:"`
	SnapshotBeforeStock     int            `json:"snapshotBeforeStock" gorm:"type:int(11);comment:"`
	SnapshotBeforeLockStock int            `json:"snapshotBeforeLockStock" gorm:"type:int(11);comment:"`
	CreatedAt               time.Time      `json:"createdAt" gorm:"comment:创建时间"`
	CreateBy                int            `json:"createBy" gorm:"index;comment:创建者"`
	CreateByName            string         `json:"createByName" gorm:"type:varchar(20);comment:创建人姓名"`
	Warehouse               Warehouse      `json:"-" gorm:"foreignKey:WarehouseCode;references:WarehouseCode"`
	LogicWarehouse          LogicWarehouse `json:"-" gorm:"foreignKey:LogicWarehouseCode;references:LogicWarehouseCode"`
}

func (StockLog) TableName() string {
	return "stock_log"
}

func (e *StockLog) GetId() interface{} {
	return e.Id
}
