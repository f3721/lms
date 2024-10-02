package models

import (
	"go-admin/common/models"
	"time"

	"gorm.io/gorm"
)

type StockOutboundProductsSubLog struct {
	models.Model

	SubId               int       `json:"subId" gorm:"type:int unsigned;comment:出库库位ID"`
	OutboundTime        time.Time `json:"outboundTime" gorm:"type:datetime;comment:库位出库时间"`
	LocationActQuantity int       `json:"locationActQuantity" gorm:"type:int unsigned;comment:实际库位出库数量"`
}

func (StockOutboundProductsSubLog) TableName() string {
	return "stock_outbound_products_sub_log"
}

func (e *StockOutboundProductsSubLog) GetId() interface{} {
	return e.Id
}

// 新增记录
func (e *StockOutboundProductsSubLog) Create(tx *gorm.DB) error {
	err := tx.Save(e).Error
	return err
}

// 根据ID查询所有
func (e *StockOutboundProductsSubLog) ListbySubIds(tx *gorm.DB, subIds []int) ([]*StockOutboundProductsSubLog, error) {
	list := []*StockOutboundProductsSubLog{}
	err := tx.Model(e).Where("sub_id in ?", subIds).Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return list, err
}
