package models

import (
	"fmt"
	"go-admin/common/global"
	"go-admin/common/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StockOutboundProductsSub struct {
	models.Model

	OutboundProductId   int `json:"outboundProductId" gorm:"type:int unsigned;comment:出库单商品id"`
	LocationQuantity    int `json:"locationQuantity" gorm:"type:int unsigned;comment:库位出库数量"`
	LocationActQuantity int `json:"locationActQuantity" gorm:"type:int unsigned;comment:实际库位出库数量"`
	StockLocationId     int `json:"stockLocationId" gorm:"type:int unsigned;comment:库位id"`
}

func (StockOutboundProductsSub) TableName() string {
	return "stock_outbound_products_sub"
}

func (e *StockOutboundProductsSub) GetId() interface{} {
	return e.Id
}

func (e *StockOutboundProductsSub) GetById(tx *gorm.DB, id int) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	return tx.Table(wcPrefix+"."+e.TableName()).Take(e, id).Error
}

// 扣减库位应出库数量
func (e *StockOutboundProductsSub) SubLocationQuantity(tx *gorm.DB, id int, quantity int) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	if err := e.GetById(tx, id); err != nil {
		return err
	}
	if e.LocationQuantity < quantity {
		return fmt.Errorf("出库SUB[%v]扣减库位应出库数量[%v] 小于扣减数量[%v]", e.Id, e.LocationQuantity, quantity)
	}

	// e.LocationQuantity -= quantity
	// res := tx.Table(wcPrefix+"."+e.TableName()).Where("location_quantity >= ?", quantity).Save(e)
	dincr := map[string]any{"location_quantity": gorm.Expr("location_quantity - ?", quantity)}
	res := tx.Debug().Table(wcPrefix+"."+e.TableName()).Where("id = ?", e.Id).Where("location_quantity = ?", e.LocationQuantity).Updates(dincr)
	err := res.Error
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		// return errors.New("StockOutboundProductsSub SubLocationQuantity Error")
		return fmt.Errorf("出库SUB[%v]扣减库位应出库数量,影响行数为0", e.Id)
	}
	return nil
}

func (e *StockOutboundProductsSub) DelById(tx *gorm.DB, id int) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	if err := tx.Table(wcPrefix+"."+e.TableName()).Delete(e, id).Error; err != nil {
		return err
	}
	/*if tx.RowsAffected == 0 {
		return errors.New("StockOutboundProductsSub DelById Error")
	}*/
	return nil
}

// 保存出库单产品子表并锁库位
func (e *StockOutboundProductsSub) SaveAndLockLocationGoods(tx *gorm.DB, goodsId int, outboundCode, outboundType string) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	tx.Table(wcPrefix + "." + e.TableName()).Omit(clause.Associations).Save(e)
	locationGoods := &StockLocationGoods{}
	if err := locationGoods.GetByGoodsIdAndLocationId(tx, goodsId, e.StockLocationId); err != nil {
		return err
	}
	// 锁定库位库存
	if err := locationGoods.Lock(tx, e.LocationQuantity, &StockLocationGoodsLockLog{
		DocketCode: outboundCode,
		FromType:   outboundType,
		Remark:     RemarkLocationProductsLockCreateOutbound,
	}); err != nil {
		return err
	}
	return nil
}

// 查询出库单下所有SUB列表
func (e *StockOutboundProductsSub) ListByOutProductIds(tx *gorm.DB, goodsIds []int) (*[]StockOutboundProductsSub, error) {
	list := &[]StockOutboundProductsSub{}
	err := tx.Table(e.TableName()).Where("outbound_product_id in ?", goodsIds).Find(list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}
