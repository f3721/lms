package models

import (
	"go-admin/common/models"
	"time"

	"gorm.io/gorm"
)

type StockEntryProductsSub struct {
	models.Model

	EntryCode        string    `json:"entryCode" gorm:"type:varchar(20);comment:入库单code"`
	EntryProductId   int       `json:"entryProductId" gorm:"type:int unsigned;comment:出库单产品ID"`
	ShouldQuantity   int       `json:"shouldQuantity" gorm:"type:int unsigned;comment:应入数量"`
	StockLocationId  int       `json:"stockLocationId" gorm:"type:int unsigned;comment:库位ID"`
	ActQuantity      int       `json:"actQuantity" gorm:"type:int unsigned;comment:实际入库数量"`
	StashLocationId  int       `json:"stashLocationId" gorm:"type:int unsigned;comment:暂存库位id"`
	StashActQuantity int       `json:"stashActQuantity" gorm:"type:int unsigned;comment:暂存数量"`
	EntryTime        time.Time `json:"entryTime" gorm:"type:datetime;comment:入库时间"`
}

func (StockEntryProductsSub) TableName() string {
	return "stock_entry_products_sub"
}

func (e *StockEntryProductsSub) GetId() interface{} {
	return e.Id
}

// 唯一主键 新增或更新数据
func (StockEntryProductsSub) CreateOrUpdate(tx *gorm.DB, data *StockEntryProductsSub) error {
	uniqueKey := StockEntryProductsSub{}
	err := tx.Where(uniqueKey).Assign(data).FirstOrCreate(data).Error
	return err
}

func (e *StockEntryProductsSub) Find(tx *gorm.DB, id int) (StockEntryProductsSub, error) {
	var subInfo StockEntryProductsSub
	err := tx.Model(e).Find(&subInfo, id).Error
	return subInfo, err
}
