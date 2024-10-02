package models

import (
	"errors"
	"go-admin/common/global"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 检查库位库存

func (e *StockLocationGoods) CheckStock(currentQuantity int) bool {
	return e.LockStock <= currentQuantity
}

// 检查库位锁定库存
func (e *StockLocationGoods) CheckLockStock(quantity int) bool {
	return e.LockStock >= quantity
}

func (e *StockLocationGoods) GetByGoodsIdAndLocationId(tx *gorm.DB, goodsId, locationId int) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	return tx.Table(wcPrefix+"."+e.TableName()).Where("goods_id = ?", goodsId).Where("location_id = ?", locationId).Take(e).Error
}

func (e *StockLocationGoods) GetByGoodsId(tx *gorm.DB, goodsIds []int) (stockLocationGoods *[]StockLocationGoods, err error) {
	stockLocationGoods = &[]StockLocationGoods{}
	err = tx.Where("goods_id = ?", goodsIds).Find(stockLocationGoods).Error
	return
}

func (e *StockLocationGoods) CheckLockStockByGoodsIdAndLocationId(tx *gorm.DB, quantity, goodsId, locationId int) bool {
	if err := e.GetByGoodsIdAndLocationId(tx, goodsId, locationId); err != nil {
		return false
	}
	return e.CheckLockStock(quantity)
}

func (e *StockLocationGoods) CheckStockByGoodsIdAndLocationId(tx *gorm.DB, stockControlProducts StockControlProducts) bool {
	if err := e.GetByGoodsIdAndLocationId(tx, stockControlProducts.GoodsId, stockControlProducts.StockLocationId); err != nil {
		return true
	}
	return e.CheckStock(stockControlProducts.CurrentQuantity)
}

// 解锁库位库存
func (e *StockLocationGoods) Unlock(tx *gorm.DB, quantity int, log *StockLocationGoodsLockLog) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	res := tx.Debug().Table(wcPrefix+"."+e.TableName()).Where("id = ?", e.Id).Where("lock_stock >= ?", quantity).Updates(map[string]interface{}{
		"stock":      gorm.Expr("stock + ?", quantity),
		"lock_stock": gorm.Expr("lock_stock - ?", quantity),
	})
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("StockLocationGoods Unlock Error")
	}
	log.SetUnLockType()
	if err := e.SaveInfoToLockLog(tx, quantity, log); err != nil {
		return err
	}
	return nil
}

// 锁定库位库存
func (e *StockLocationGoods) Lock(tx *gorm.DB, quantity int, log *StockLocationGoodsLockLog) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	res := tx.Table(wcPrefix+"."+e.TableName()).Where("id = ?", e.Id).Where("stock >= ?", quantity).Updates(map[string]interface{}{
		"stock":      gorm.Expr("stock - ?", quantity),
		"lock_stock": gorm.Expr("lock_stock + ?", quantity),
	})
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("StockLocationGoods Lock Error")
	}
	log.SetLockType()
	if err := e.SaveInfoToLockLog(tx, quantity, log); err != nil {
		return err
	}
	return nil
}

// 保存锁定库位库存日志
func (e *StockLocationGoods) SaveInfoToLockLog(tx *gorm.DB, quantity int, log *StockLocationGoodsLockLog) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	// 再次查询 防止并发导致数据异常
	if err := tx.Table(wcPrefix + "." + e.TableName()).Take(e).Error; err != nil {
		return nil
	}
	log.AfterStock = e.Stock
	log.AfterLockStock = e.LockStock
	log.StockLocationGoodsId = e.Id
	log.LockQuantity = quantity
	log.GoodsId = e.GoodsId
	log.LocationId = e.LocationId
	return tx.Table(wcPrefix + "." + log.TableName()).Save(log).Error
}

// 扣减锁定库位库存
func (e *StockLocationGoods) SubLockStock(tx *gorm.DB, quantity int, log *StockLocationGoodsLog) error {
	res := tx.Model(e).Where("lock_stock >= ?", quantity).Updates(map[string]interface{}{
		"lock_stock":     gorm.Expr("lock_stock - ?", quantity),
		"update_by":      log.CreateBy,
		"update_by_name": log.CreateByName,
	})
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("StockLocationGoods SubLockStock Error")
	}
	if err := e.SaveInfoToLog(tx, quantity*(-1), log, true); err != nil {
		return err
	}
	return nil
}

// 扣减库位库存

func (e *StockLocationGoods) SubStock(tx *gorm.DB, quantity int, log *StockLocationGoodsLog) error {
	res := tx.Model(e).Where("stock >= ?", quantity).Updates(map[string]interface{}{
		"stock":          gorm.Expr("stock - ?", quantity),
		"update_by":      log.CreateBy,
		"update_by_name": log.CreateByName,
	})
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("StockLocationGoods SubStock Error")
	}
	if err := e.SaveInfoToLog(tx, quantity*(-1), log, false); err != nil {
		return err
	}
	return nil
}

// 增加库位库存
func (e *StockLocationGoods) AddStock(tx *gorm.DB, quantity int, log *StockLocationGoodsLog) error {
	if e.Id == 0 {
		e.Stock = quantity
		if err := tx.Omit(clause.Associations).Create(e).Error; err != nil {
			return err
		}
	} else {
		if err := tx.Model(e).Updates(map[string]interface{}{
			"stock":          gorm.Expr("stock + ?", quantity),
			"update_by":      log.CreateBy,
			"update_by_name": log.CreateByName,
		}).Error; err != nil {
			return err
		}
	}

	if err := e.SaveInfoToLog(tx, quantity, log, false); err != nil {
		return err
	}
	return nil
}

// 保存库位库存日志
func (e *StockLocationGoods) SaveInfoToLog(tx *gorm.DB, changeStock int, log *StockLocationGoodsLog, isLock bool) error {
	// 再次查询 防止并发导致数据异常
	if err := tx.Take(e).Error; err != nil {
		return err
	}
	log.StockLocationGoodsId = e.Id
	log.BeforeStock = e.Stock + e.LockStock - changeStock
	log.AfterStock = e.Stock + e.LockStock
	log.ChangeStock = changeStock

	log.SnapshotAfterStock = e.Stock
	log.SnapshotAfterLockStock = e.LockStock
	if isLock {
		log.SnapshotBeforeLockStock = e.LockStock - changeStock
		log.SnapshotBeforeStock = e.Stock
	} else {
		log.SnapshotBeforeLockStock = e.LockStock
		log.SnapshotBeforeStock = e.Stock - changeStock
	}

	return tx.Save(log).Error
}

// 插入库位库存数据时，初始化相关信息
func (e *StockLocationGoods) SetBaseInfo(goodsId, locationId, createBy int, createByName string) {
	e.GoodsId = goodsId
	e.LocationId = locationId
	e.CreateBy = createBy
	e.CreateByName = createByName
}
