package models

import (
	"errors"
	"fmt"
	modelsOc "go-admin/app/oc/models"
	"go-admin/common/global"
	"go-admin/common/models"
	"go-admin/common/utils"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	OutboundStatus0 = "0"
	OutboundStatus1 = "1"
	OutboundStatus3 = "3"
	OutboundStatus2 = "2"

	OutboundType0 = "0"
	OutboundType1 = "1"

	OutboundSourceType0 = "0"
	OutboundSourceType1 = "1"

	OutboundLogModelName                     = "stockOutbound"
	OutboundLogModelTypeConfirm              = "confirmOutbound"
	OutboundLogModelTypeCancelForCsOrder     = "cancelForCsOrder"
	OutboundLogModelTypePartCancelForCsOrder = "partCancelForCsOrder"
)

var (
	OutboundStatusMap = map[string]string{
		OutboundStatus0: "已作废",
		OutboundStatus1: "未出库",
		OutboundStatus3: "部分出库",
		OutboundStatus2: "出库完成",
	}

	OutboundTypeMap = map[string]string{
		OutboundType0: "大货出库",
		OutboundType1: "订单出库",
	}

	OutboundSourceTypeMap = map[string]string{
		OutboundSourceType0: "调拨单",
		OutboundSourceType1: "领用单",
	}
)

type StockOutbound struct {
	models.Model

	OutboundCode          string                  `json:"outboundCode" gorm:"type:varchar(20);comment:出库单编码"`
	Type                  string                  `json:"type" gorm:"type:tinyint;comment:出库类型:  0 大货出库  1 订单出库  2 其他"`
	Status                string                  `json:"status" gorm:"type:tinyint;comment:状态:0-已作废 1-创建 3-部分出库 2-出库完成"`
	SourceCode            string                  `json:"sourceCode" gorm:"type:varchar(32);comment:来源单据code"`
	Remark                string                  `json:"remark" gorm:"type:varchar(255);comment:备注"`
	OutboundTime          time.Time               `json:"outboundTime" gorm:"type:datetime;comment:首次出库时间"`
	OutboundEndTime       time.Time               `json:"outboundEndTime" gorm:"type:datetime;comment:出库完成时间"`
	WarehouseCode         string                  `json:"warehouseCode" gorm:"type:varchar(20);comment:实体仓code"`
	LogicWarehouseCode    string                  `json:"logicWarehouseCode" gorm:"type:varchar(20);comment:逻辑仓code"`
	VendorId              int                     `json:"vendorId" gorm:"type:int unsigned;comment:货主id"`
	Warehouse             Warehouse               `json:"-" gorm:"foreignKey:WarehouseCode;references:WarehouseCode"`
	LogicWarehouse        LogicWarehouse          `json:"-" gorm:"foreignKey:LogicWarehouseCode;references:LogicWarehouseCode"`
	StockOutboundProducts []StockOutboundProducts `json:"-" gorm:"foreignKey:OutboundCode;references:OutboundCode"`
	models.ModelTime
	models.ControlBy
}

func (StockOutbound) TableName() string {
	return "stock_outbound"
}

func (e *StockOutbound) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *StockOutbound) GetId() interface{} {
	return e.Id
}

//是否领用单出库

func (e *StockOutbound) IsOrderOutbound() bool {
	return e.Type == OutboundType1
}

//是否调拨单出库

func (e *StockOutbound) IsTransferOutbound() bool {
	return e.Type == OutboundType0
}

// 新增出库单
func (e *StockOutbound) InsertOutbound(tx *gorm.DB, Type string) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	if _, err := e.GenerateOutboundCode(tx); err != nil {
		return err
	}
	e.SetTypeForOutbound(Type)

	if e.IsOrderOutbound() {
		e.SetProductsForOutboundForOrder()
	} else {
		e.SetProductsForOutbound()
	}

	err := tx.Table(wcPrefix + ".stock_outbound").Omit(clause.Associations).Create(e).Error
	if err != nil {
		return err
	}

	err = tx.Table(wcPrefix + ".stock_outbound_products").Omit(clause.Associations).Create(e.StockOutboundProducts).Error
	if err != nil {
		return err
	}
	return nil
}

// 根据类型设置出库单类型
func (e *StockOutbound) SetTypeForOutbound(Type string) {
	e.Type = Type
}

// 设置产品信息
func (e *StockOutbound) SetProductsForOutbound() {
	for index := range e.StockOutboundProducts {
		e.StockOutboundProducts[index].OutboundCode = e.OutboundCode
		e.StockOutboundProducts[index].VendorId = e.VendorId
		e.StockOutboundProducts[index].WarehouseCode = e.WarehouseCode
		e.StockOutboundProducts[index].LogicWarehouseCode = e.LogicWarehouseCode
		e.StockOutboundProducts[index].OriQuantity = e.StockOutboundProducts[index].Quantity
	}
}

func (e *StockOutbound) SetProductsForOutboundForOrder() {
	for index := range e.StockOutboundProducts {
		e.StockOutboundProducts[index].OutboundCode = e.OutboundCode
		e.StockOutboundProducts[index].WarehouseCode = e.WarehouseCode
		e.StockOutboundProducts[index].LogicWarehouseCode = e.LogicWarehouseCode
		e.StockOutboundProducts[index].OriQuantity = e.StockOutboundProducts[index].Quantity
	}
}

// GenerateOutboundCode
func (e *StockOutbound) GenerateOutboundCode(tx *gorm.DB) (string, error) {
	var count int64
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)

	start, end := utils.GetTodayTime()
	if err := tx.Table(wcPrefix+"."+e.TableName()).Where("created_at BETWEEN ? AND ?", start, end).Count(&count).Error; err != nil {
		return "", err
	}
	code := "OR" + time.Now().Format("20060102150405") + fmt.Sprintf("%04d", count+1)
	e.OutboundCode = code
	return code, nil
}

func (e *StockOutbound) GetBySourceCode(tx *gorm.DB, sourceCode string) {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	tx.Table(wcPrefix+"."+e.TableName()).Where("source_code = ?", sourceCode).Take(e)
}

// 手动出库
func (e *StockOutbound) ConfirmOutboundManual(tx *gorm.DB, list []OutboundProductSubCustom) error {
	e.SetProductsActQuantity(list)
	// 更改出库单状态
	if err := tx.Omit(clause.Associations).Save(e).Error; err != nil {
		return err
	}

	for _, item := range list {
		// 更新出库商品子表实际数量
		if err := item.SetActQuantityForProductSub(tx); err != nil {
			return err
		}
		// 检查库位库存 - start
		locationGoods := &StockLocationGoods{}
		if err := locationGoods.GetByGoodsIdAndLocationId(tx, item.GoodsId, item.StockLocationId); err != nil {
			return err
		}
		if !locationGoods.CheckLockStock(item.LocationActQuantity) {
			return errors.New(item.SkuCode + "库位:" + item.StockLocation.LocationCode + ",锁定库位库存不足")
		}
		// 检查库位库存 - end

		// 实际出库数量不等于应出库数量(这里只可能是调拨单类型手动出库)，释放多锁住库位库存
		LocationDiffQuantity := item.LocationQuantity - item.LocationActQuantity
		if LocationDiffQuantity > 0 {
			if err := locationGoods.Unlock(tx, LocationDiffQuantity, &StockLocationGoodsLockLog{
				DocketCode: e.OutboundCode,
				FromType:   e.Type,
				Remark:     RemarkManualConfirmOutboundLocationProductsNotEqual,
			}); err != nil {
				return err
			}
		}
		// 扣减库位库存
		if err := locationGoods.SubLockStock(tx, item.LocationActQuantity, &StockLocationGoodsLog{
			DocketCode:   e.OutboundCode,
			FromType:     StockLocationGoodsLogFromType0,
			CreateBy:     e.UpdateBy,
			CreateByName: e.UpdateByName,
		}); err != nil {
			return err
		}
	}

	for _, item := range e.StockOutboundProducts {
		// 更新出库商品表实际数量
		if err := item.Save(tx); err != nil {
			return err
		}
		// 检查库存 - start
		stockInfo := &StockInfo{}
		if !stockInfo.CheckLockStockByGoodsIdAndLogicWarehouseCode(tx, item.ActQuantity, item.GoodsId, item.LogicWarehouseCode) {
			return errors.New(item.SkuCode + "逻辑仓:" + item.LogicWarehouseCode + ",锁定库存不足")
		}
		// 检查库存 - end

		// 实际出库数量不等于应出库数量(这里只可能是调拨单类型手动出库)，释放多锁住库存
		diffQuantity := item.Quantity - item.ActQuantity
		if diffQuantity > 0 {
			if err := stockInfo.Unlock(tx, diffQuantity, &StockLockLog{
				DocketCode: e.OutboundCode,
				FromType:   e.Type,
				Remark:     RemarkManualConfirmOutboundProductsNotEqual,
			}); err != nil {
				return err
			}
		}
		// 扣减库存
		if err := stockInfo.SubLockStock(tx, item.ActQuantity, &StockLog{
			DocketCode:   e.OutboundCode,
			FromType:     StockLogFromType0,
			CreateBy:     e.UpdateBy,
			CreateByName: e.UpdateByName,
		}); err != nil {
			return err
		}
	}

	// 判断是否来源于调拨单，并且出库仓为虚拟仓
	// 更新订单状态 更新订单发货数量
	if e.IsTransferOutbound() {
		transfer := &StockTransfer{}
		if err := transfer.GetByTransferCodeWithOptions(tx, e.SourceCode, func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("ToWarehouse")
		}); err != nil {
			return err
		}
		if err := transfer.SetOutboundCompleteStatus(tx); err != nil {
			return err
		}
		// 虚拟出库仓自动出
		if transfer.ToWarehouse.CheckIsVirtual() {
			stockEntry := &StockEntry{}
			if err := stockEntry.GetByEntrySourceCodeWithOptions(tx, e.SourceCode, func(tx *gorm.DB) *gorm.DB {
				return tx.Preload("StockEntryProducts")
			}); err != nil {
				return err
			}
			stockEntry.SetConfirmEntryStatus(e.UpdateBy, e.UpdateByName)
			if err := stockEntry.ConfirmEntryAutoForVirtual(tx); err != nil {
				return err
			}
			if err := transfer.SetEntryCompleteStatus(tx); err != nil {
				return err
			}
		}
	} else if e.IsOrderOutbound() {
		goodsQuantityMap := lo.Associate(e.StockOutboundProducts, func(p StockOutboundProducts) (int, int) {
			return p.GoodsId, p.ActQuantity
		})
		modelOc := &modelsOc.OrderInfo{}
		if err := modelOc.ShipmentUpdateOrder(tx, e.SourceCode, goodsQuantityMap); err != nil {
			return err
		}
	}
	return nil
}

// 获取聚合实际出库数量
func (e *StockOutbound) SetProductsActQuantity(list []OutboundProductSubCustom) {
	for index, item := range e.StockOutboundProducts {
		actQuantity := 0
		for _, itemReq := range list {
			if item.GoodsId == itemReq.GoodsId {
				actQuantity += itemReq.LocationActQuantity
			}
		}
		e.StockOutboundProducts[index].ActQuantity = actQuantity
	}
}

// 自动出库（虚拟仓）
func (e *StockOutbound) ConfirmOutboundAutoForVirtual(tx *gorm.DB) error {
	// 实际出库数量等于应出库数量,时间，状态
	e.SetProductsActQuantityComplete()
	// 更改出库单状态
	if err := tx.Omit(clause.Associations).Save(e).Error; err != nil {
		return err
	}

	for _, item := range e.StockOutboundProducts {
		if err := item.Save(tx); err != nil {
			return err
		}
		stockInfo := &StockInfo{}
		if err := stockInfo.GetByVendorIdAndLogicWarehouseCodeAndSku(tx, item.VendorId, item.LogicWarehouseCode, item.SkuCode); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				stockInfo.SetBaseInfo(item.GoodsId, item.VendorId, e.UpdateBy, item.SkuCode, item.WarehouseCode, item.LogicWarehouseCode, e.UpdateByName)
			} else {
				return err
			}
		}
		// 扣减库存
		if err := stockInfo.SubStockForVirtualWarehouse(tx, item.ActQuantity, &StockLog{
			DocketCode:   e.OutboundCode,
			FromType:     StockLogFromType0,
			CreateBy:     e.UpdateBy,
			CreateByName: e.UpdateByName,
		}); err != nil {
			return err
		}
	}
	return nil
}

// 自动出库(次品入库创建的出库单)
func (e *StockOutbound) ConfirmOutboundForDefectiveEntry(tx *gorm.DB) error {
	for _, item := range e.StockOutboundProducts {
		locationGoods := &StockLocationGoods{}
		if err := locationGoods.GetByGoodsIdAndLocationId(tx, item.GoodsId, item.StockLocationId); err != nil {
			return err
		}
		// 扣减库位库存
		if err := locationGoods.SubStock(tx, item.ActQuantity, &StockLocationGoodsLog{
			DocketCode:   e.OutboundCode,
			FromType:     StockLocationGoodsLogFromType0,
			CreateBy:     e.CreateBy,
			CreateByName: e.CreateByName,
		}); err != nil {
			return err
		}
		// 检查库存 - start
		stockInfo := &StockInfo{}
		if err := stockInfo.GetByGoodsIdAndLogicWarehouseCode(tx, item.GoodsId, item.LogicWarehouseCode); err != nil {
			return err
		}
		// 扣减库存
		if err := stockInfo.SubStock(tx, item.ActQuantity, &StockLog{
			DocketCode:   e.OutboundCode,
			FromType:     StockLogFromType0,
			CreateBy:     e.CreateBy,
			CreateByName: e.CreateByName,
		}); err != nil {
			return err
		}
	}
	return nil
}

// 自动出库设置实际出库商品数量为出库数量
func (e *StockOutbound) SetProductsActQuantityComplete() {
	for index, item := range e.StockOutboundProducts {
		currTime := time.Now()
		e.StockOutboundProducts[index].OutboundTime = currTime
		e.StockOutboundProducts[index].OutboundEndTime = currTime
		e.StockOutboundProducts[index].Status = 2
		// 时间出库数量 = 应出库数量
		e.StockOutboundProducts[index].ActQuantity = item.Quantity
	}
}

// 自动出库出库单出库状态设置
func (e *StockOutbound) SetConfirmOutboundStatus(updateBy int, UpdateByName string) {
	currTime := time.Now()
	e.Status = OutboundStatus2
	e.OutboundTime = currTime
	e.OutboundEndTime = currTime
	e.UpdateBy = updateBy
	e.UpdateByName = UpdateByName
	for index := range e.StockOutboundProducts {
		e.StockOutboundProducts[index].UpdateBy = updateBy
		e.StockOutboundProducts[index].UpdateByName = UpdateByName
	}
}

func (e *StockLocationGoodsLog) GenerateStockLocationGoodsLogCode(tx *gorm.DB) (string, error) {
	var count int64
	start, end := utils.GetTodayTime()
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)

	if err := tx.Table(wcPrefix+"."+e.TableName()).
		Where("created_at BETWEEN ? AND ?", start, end).
		Where("from_type = ?", "3").
		Count(&count).Error; err != nil {
		return "", err
	}
	code := "TR" + time.Now().Format("20060102150405") + fmt.Sprintf("%04d", count/2+1)
	e.DocketCode = code
	return code, nil
}
