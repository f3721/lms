package models

import (
	"errors"
	"go-admin/common/global"
	"go-admin/common/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	InitVirtualWarehouseStock = 99999999
)

type StockInfo struct {
	models.Model

	GoodsId            int    `json:"goodsId" gorm:"type:int unsigned;comment:goods表ID"`
	LogicWarehouseCode string `json:"logicWarehouseCode" gorm:"type:varchar(20);comment:逻辑仓code"`
	VendorId           int    `json:"vendorId" gorm:"type:int unsigned;comment:货主id"`
	WarehouseCode      string `json:"warehouseCode" gorm:"type:varchar(20);comment:实体仓code"`
	SkuCode            string `json:"skuCode" gorm:"type:varchar(10);comment:sku"`
	Stock              int    `json:"stock" gorm:"type:int unsigned;comment:库存数量"`
	LockStock          int    `json:"lockStock" gorm:"type:int unsigned;comment:锁定库存"`
	models.ModelTime
	models.ControlBy
	Warehouse      Warehouse      `json:"-" gorm:"foreignKey:WarehouseCode;references:WarehouseCode"`
	LogicWarehouse LogicWarehouse `json:"-" gorm:"foreignKey:LogicWarehouseCode;references:LogicWarehouseCode"`
}

func (StockInfo) TableName() string {
	return "stock_info"
}

func (e *StockInfo) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *StockInfo) GetId() interface{} {
	return e.Id
}

// GetStockInfos 获取库存
func (e *StockInfo) GetStockInfos(tx *gorm.DB, vendorId int, logicWarehouseCode string, skuCodes []string) (stockInfos *[]StockInfo, err error) {
	stockInfos = &[]StockInfo{}
	err = tx.Where("sku_code = ?", skuCodes).Where("vendor_id = ?", vendorId).Where("logic_warehouse_code = ?", logicWarehouseCode).Find(stockInfos).Error
	return
}

// 检查库存
func (e *StockInfo) CheckStock(quantity int) bool {
	return e.Stock >= quantity
}

// 检查锁定库存
func (e *StockInfo) CheckLockStock(quantity int) bool {
	return e.LockStock >= quantity
}

// 批量查库存信息 | LwsCode+GoodIds
func (e *StockInfo) ListByLwsCodeAndGoogIds(tx *gorm.DB, logicWarehouseCode string, goodsIds []int) (*[]*StockInfo, error) {
	list := &[]*StockInfo{}
	err := tx.Table(e.TableName()).Where("logic_warehouse_code = ?", logicWarehouseCode).
		Where("goods_id in ?", goodsIds).Find(list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

// 检查库存
func (e *StockInfo) CheckStockByGoodsIdAndLogicWarehouseCode(tx *gorm.DB, quantity, goodsId int, logicWarehouseCode string) bool {
	if err := e.GetByGoodsIdAndLogicWarehouseCode(tx, goodsId, logicWarehouseCode); err != nil {
		return false
	}
	return e.CheckStock(quantity)
}

// CheckQuantityByGoodsIdAndLogicWarehouseCode 判断变动数量 是否在库存的允许范围内
func (e *StockInfo) CheckQuantityByGoodsIdAndLogicWarehouseCode(tx *gorm.DB, quantity, goodsId int, logicWarehouseCode string) bool {
	if err := e.GetByGoodsIdAndLogicWarehouseCode(tx, goodsId, logicWarehouseCode); err != nil {
		return false
	}
	if e.Stock+quantity < 0 {
		return false
	}
	return true
}

// 检查锁定库存
func (e *StockInfo) CheckLockStockByGoodsIdAndLogicWarehouseCode(tx *gorm.DB, quantity, goodsId int, logicWarehouseCode string) bool {
	if err := e.GetByGoodsIdAndLogicWarehouseCode(tx, goodsId, logicWarehouseCode); err != nil {
		return false
	}
	return e.CheckLockStock(quantity)
}

// 根据goodsId和逻辑仓code获取库存信息
func (e *StockInfo) GetByGoodsIdAndLogicWarehouseCode(tx *gorm.DB, goodsId int, logicWarehouseCode string) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	return tx.Table(wcPrefix+"."+e.TableName()).Where("goods_id = ?", goodsId).Where("logic_warehouse_code = ?", logicWarehouseCode).Take(e).Error
}

// 根据vendorId和逻辑仓code和skuCode获取库存信息
func (e *StockInfo) GetByVendorIdAndLogicWarehouseCodeAndSku(tx *gorm.DB, vendorId int, logicWarehouseCode, SkuCode string) error {
	return tx.Where("sku_code = ?", SkuCode).Where("vendor_id = ?", vendorId).Where("logic_warehouse_code = ?", logicWarehouseCode).Take(e).Error
}

// 插入库存数据时，初始化相关信息
func (e *StockInfo) SetBaseInfo(goodsId, vendorId, createBy int, skuCode, warehouseCode, logicWarehouseCode, createByName string) {
	e.GoodsId = goodsId
	e.WarehouseCode = warehouseCode
	e.LogicWarehouseCode = logicWarehouseCode
	e.VendorId = vendorId
	e.SkuCode = skuCode
	e.CreateBy = createBy
	e.CreateByName = createByName
}

// 解锁库存
func (e *StockInfo) Unlock(tx *gorm.DB, quantity int, log *StockLockLog) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	res := tx.Table(wcPrefix+"."+e.TableName()).Where("id = ?", e.Id).Where("lock_stock >= ?", quantity).Updates(map[string]interface{}{
		"stock":      gorm.Expr("stock + ?", quantity),
		"lock_stock": gorm.Expr("lock_stock - ?", quantity),
	})
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("StockInfo Unlock Error")
	}
	log.SetUnLockType()
	if err := e.SaveInfoToLockLog(tx, quantity, log); err != nil {
		return err
	}
	return nil
}

// 锁定库存
func (e *StockInfo) Lock(tx *gorm.DB, quantity int, log *StockLockLog) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	res := tx.Table(wcPrefix+"."+e.TableName()).Where("id = ?", e.Id).Where("stock >= ?", quantity).Updates(map[string]interface{}{
		"stock":      gorm.Expr("stock - ?", quantity),
		"lock_stock": gorm.Expr("lock_stock + ?", quantity),
	})
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("StockInfo Lock Error")
	}
	log.SetLockType()
	if err := e.SaveInfoToLockLog(tx, quantity, log); err != nil {
		return err
	}
	return nil
}

// 保存锁定库存日志
func (e *StockInfo) SaveInfoToLockLog(tx *gorm.DB, quantity int, log *StockLockLog) error {
	// 再次查询 防止并发导致数据异常
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	if err := tx.Table(wcPrefix + "." + e.TableName()).Take(e).Error; err != nil {
		return nil
	}
	log.AfterStock = e.Stock
	log.AfterLockStock = e.LockStock
	log.StockInfoId = e.Id
	log.LockQuantity = quantity
	log.GoodsId = e.GoodsId
	log.LogicWarehouseCode = e.LogicWarehouseCode
	return tx.Table(wcPrefix + "." + log.TableName()).Save(log).Error
}

// 扣减锁定库存
func (e *StockInfo) SubLockStock(tx *gorm.DB, quantity int, log *StockLog) error {
	res := tx.Model(e).Where("lock_stock >= ?", quantity).Updates(map[string]interface{}{
		"lock_stock":     gorm.Expr("lock_stock - ?", quantity),
		"update_by":      log.CreateBy,
		"update_by_name": log.CreateByName,
	})
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("StockInfo SubLockStock Error")
	}
	if err := e.SaveInfoToLog(tx, quantity*(-1), log, true); err != nil {
		return err
	}
	return nil
}

// 扣减库存
func (e *StockInfo) SubStock(tx *gorm.DB, quantity int, log *StockLog) error {
	res := tx.Model(e).Where("stock >= ?", quantity).Updates(map[string]interface{}{
		"stock":          gorm.Expr("stock - ?", quantity),
		"update_by":      log.CreateBy,
		"update_by_name": log.CreateByName,
	})
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("StockInfo SubStock Error")
	}
	if err := e.SaveInfoToLog(tx, quantity*(-1), log, false); err != nil {
		return err
	}
	return nil
}

// 扣减库存(虚拟仓)
func (e *StockInfo) SubStockForVirtualWarehouse(tx *gorm.DB, quantity int, log *StockLog) error {
	if e.Id == 0 {
		e.Stock = InitVirtualWarehouseStock - quantity
		if err := tx.Omit(clause.Associations).Create(e).Error; err != nil {
			return err
		}
		if err := e.SaveInfoToLog(tx, quantity*(-1), log, false); err != nil {
			return err
		}
	} else {
		if err := e.SubStock(tx, quantity, log); err != nil {
			return err
		}
	}
	return nil
}

// 增加库存
func (e *StockInfo) AddStock(tx *gorm.DB, quantity int, log *StockLog) error {
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

// 增加库存
func (e *StockInfo) AddStockForVirtualWarehouse(tx *gorm.DB, quantity int, log *StockLog) error {
	if e.Id == 0 {
		e.Stock = InitVirtualWarehouseStock + quantity
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

// 保存库存日志
func (e *StockInfo) SaveInfoToLog(tx *gorm.DB, changeStock int, log *StockLog, isLock bool) error {
	// 再次查询 防止并发导致数据异常
	if err := tx.Take(e).Error; err != nil {
		return err
	}
	log.StockInfoId = e.Id
	log.WarehouseCode = e.WarehouseCode
	log.LogicWarehouseCode = e.LogicWarehouseCode
	log.SkuCode = e.SkuCode
	log.VendorId = e.VendorId
	log.SnapshotAfterStock = e.Stock
	log.SnapshotAfterLockStock = e.LockStock
	log.BeforeStock = e.Stock + e.LockStock - changeStock
	log.AfterStock = e.Stock + e.LockStock
	log.ChangeStock = changeStock
	if isLock {
		log.SnapshotBeforeLockStock = e.LockStock - changeStock
		log.SnapshotBeforeStock = e.Stock
	} else {
		log.SnapshotBeforeLockStock = e.LockStock
		log.SnapshotBeforeStock = e.Stock - changeStock
	}
	return tx.Save(log).Error
}

func (e *StockInfo) GetByGoodsIdAndWarehouseCode(tx *gorm.DB, goodsId int, warehouseCode string) error {
	logicWarehouse := &LogicWarehouse{}
	if err := logicWarehouse.GetPassLogicWarehouseByWhCode(tx, warehouseCode); err != nil {
		return err
	}
	if err := e.GetByGoodsIdAndLogicWarehouseCode(tx, goodsId, logicWarehouse.LogicWarehouseCode); err != nil {
		return errors.New("库存信息不存在")
	}
	return nil
}

func LockStockInfoForOrder(tx *gorm.DB, quantity, goodsId int, warehouseCode, orderId, remark string) error {
	stock := &StockInfo{}
	if err := stock.GetByGoodsIdAndWarehouseCode(tx, goodsId, warehouseCode); err != nil {
		return err
	}
	if !stock.CheckStock(quantity) {
		return errors.New("库存不足")
	}
	if remark == "" {
		remark = RemarkStockLockLogOrderAboutLock
	}
	return stock.Lock(tx, quantity, &StockLockLog{
		DocketCode: orderId,
		FromType:   StockLockLogFromType1,
		Remark:     remark,
	})
}

func UnLockStockInfoForOrder(tx *gorm.DB, quantity, goodsId int, warehouseCode, orderId, remark string) error {
	stock := &StockInfo{}
	if err := stock.GetByGoodsIdAndWarehouseCode(tx, goodsId, warehouseCode); err != nil {
		return err
	}
	if !stock.CheckLockStock(quantity) {
		return errors.New("锁定库存不足")
	}
	if remark == "" {
		remark = RemarkStockLockLogOrderAboutUnLock
	}
	return stock.Unlock(tx, quantity, &StockLockLog{
		DocketCode: orderId,
		FromType:   StockLockLogFromType1,
		Remark:     remark,
	})
}
