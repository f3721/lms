package models

import (
	"go-admin/common/models"
	"time"
)

type StockLocationGoodsLog struct {
	models.Model

	StockLocationGoodsId    int       `json:"stockLocationGoodsId" gorm:"type:int(11) unsigned;comment:库位商品表id"`
	BeforeStock             int       `json:"beforeStock" gorm:"type:int(11) unsigned;comment:期初在库数"`
	AfterStock              int       `json:"afterStock" gorm:"type:int(11) unsigned;comment:期末在库数"`
	ChangeStock             int       `json:"changeStock" gorm:"type:int(11);comment:变动数量"`
	SnapshotAfterStock      int       `json:"snapshotAfterStock" gorm:"type:int(11);comment:"`
	SnapshotAfterLockStock  int       `json:"snapshotAfterLockStock" gorm:"type:int(11);comment:"`
	SnapshotBeforeStock     int       `json:"snapshotBeforeStock" gorm:"type:int(11);comment:"`
	SnapshotBeforeLockStock int       `json:"snapshotBeforeLockStock" gorm:"type:int(11);comment:"`
	DocketCode              string    `json:"docketCode" gorm:"type:varchar(20);comment:单据编号"`
	FromType                string    `json:"fromType" gorm:"type:tinyint(1) unsigned;comment:0 出库单 1 入库单 2 库存调整单 3 库位转移"`
	CreatedAt               time.Time `json:"createdAt" gorm:"comment:创建时间"`
	CreateBy                int       `json:"createBy" gorm:"index;comment:创建者"`
	CreateByName            string    `json:"createByName" gorm:"type:varchar(20);comment:创建人姓名"`
}

func (StockLocationGoodsLog) TableName() string {
	return "stock_location_goods_log"
}

func (e *StockLocationGoodsLog) GetId() interface{} {
	return e.Id
}
