package models

import (
	"go-admin/common/models"
	"time"

	"gorm.io/gorm"
)

type StockLocationGoods struct {
	models.Model

	GoodsId    int       `json:"goodsId" gorm:"type:int(11);comment:goods表ID"`
	LocationId int       `json:"locationId" gorm:"type:int(11);comment:库位表ID"`
	Stock      int       `json:"stock" gorm:"type:int(11);comment:库位商品数量"`
	LockStock  int       `json:"lockStock" gorm:"type:int(11);comment:锁定库存"`
	CreatedAt  time.Time `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`
	//StockLocation  *StockLocation `json:"stockLocation" gorm:"foreignkey:location_id"`
	models.ControlBy
}

func (StockLocationGoods) TableName() string {
	return "stock_location_goods"
}

func (e *StockLocationGoods) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *StockLocationGoods) GetId() interface{} {
	return e.Id
}

func (e *StockLocationGoods) GetTotalStockByLocationCode(tx *gorm.DB, locationCode string) (total int) {
	tx.Table(e.TableName()+" t").Select("sum(t.stock+t.lock_stock) total").Joins("left join stock_location sl on t.location_id = sl.id").Where("sl.location_code = ?", locationCode).Where("sl.status = 1").Scan(&total)
	return
}

// 批量查询商品库存数量
func (e *StockLocationGoods) ListByGoodsIdsAndLocationIds(tx *gorm.DB, goodIds []int, locationIds []int) (*[]StockLocationGoods, error) {
	list := &[]StockLocationGoods{}
	err := tx.Table(e.TableName()).Where("goods_id in ?", goodIds).Where("location_id in ?", locationIds).Find(list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}
