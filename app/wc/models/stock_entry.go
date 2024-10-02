package models

import (
	"errors"
	"fmt"
	modelsCs "go-admin/app/oc/models"
	"go-admin/common/global"
	"go-admin/common/models"
	"go-admin/common/utils"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	CheckStatus0 = 0
	CheckStatus1 = "1"
	CheckStatus2 = "2"
	CheckStatus3 = "3"

	EntryStatus0 = "0"
	EntryStatus1 = "1"
	EntryStatus3 = "3"
	EntryStatus2 = "2"

	EntryType0 = "0"
	EntryType1 = "1"
	EntryType2 = "2"
	EntryType3 = "3"

	EntrySourceType0 = "0"
	EntrySourceType1 = "1"

	EntryLogModelName                       = "stockEntry"
	EntryLogModelTypeConfirm                = "confirmEntry"
	EntryLogModelTypeStash                  = "Stash"
	EntryLogModelTypeCancelForCancelCsOrder = "cancelForCancelCsOrder"
)

var (
	EntryStatusMap = map[string]string{
		EntryStatus0: "已作废",
		EntryStatus1: "未入库",
		EntryStatus3: "部分入库",
		EntryStatus2: "入库完成",
	}

	//审核状态：1已提交，2待审核，3
	EntryCheckStatusMap = map[string]string{
		CheckStatus1: "待审核",
		CheckStatus2: "审核通过",
		CheckStatus3: "审核驳回",
	}

	EntryTypeMap = map[string]string{
		EntryType0: "大货入库",
		EntryType1: "退货入库",
		EntryType3: "采购入库",
	}

	EntrySourceTypeMap = map[string]string{
		EntrySourceType0: "调拨单",
		EntrySourceType1: "售后单",
	}
)

type StockEntry struct {
	models.Model

	EntryCode          string               `json:"entryCode" gorm:"type:varchar(20);comment:入库单编码"`
	Type               string               `json:"type" gorm:"type:tinyint;comment:入库类型:  0 大货入库  1 退货入库 3采购入库"`
	Status             string               `json:"status" gorm:"type:tinyint;comment:状态:0-已作废 1-未入库 3-部分入库 2-已完成"`
	SourceCode         string               `json:"sourceCode" gorm:"type:varchar(32);comment:来源单据code"`
	SupplierId         int                  `json:"supplierId" gorm:"type:int;comment:供应商id"`
	CheckStatus        string               `json:"checkStatus" gorm:"type:varchar(32);comment:审核状态：1 待审核 2 审核通过 3 审核驳回"`
	CheckRemark        string               `json:"checkRemark" gorm:"type:varchar(500);comment:审核驳回内容"`
	Remark             string               `json:"remark" gorm:"type:varchar(255);comment:备注"`
	EntryTime          time.Time            `json:"entryTime" gorm:"type:datetime;comment:最早入库时间"`
	EntryEndTime       time.Time            `json:"entryEndTime" gorm:"type:datetime;comment:入库完成时间"`
	WarehouseCode      string               `json:"warehouseCode" gorm:"type:varchar(20);comment:实体仓code"`
	LogicWarehouseCode string               `json:"logicWarehouseCode" gorm:"type:varchar(20);comment:逻辑仓code"`
	VendorId           int                  `json:"vendorId" gorm:"type:int unsigned;comment:货主id"`
	Warehouse          Warehouse            `json:"-" gorm:"foreignKey:WarehouseCode;references:WarehouseCode"`
	LogicWarehouse     LogicWarehouse       `json:"-" gorm:"foreignKey:LogicWarehouseCode;references:LogicWarehouseCode"`
	Supplier           Supplier             `json:"-" gorm:"foreignKey:SupplierId;references:id"`
	StockEntryProducts []StockEntryProducts `json:"-" gorm:"foreignKey:EntryCode;references:EntryCode"`
	models.ModelTime
	models.ControlBy
}

func (StockEntry) TableName() string {
	return "stock_entry"
}

func (e *StockEntry) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *StockEntry) GetId() interface{} {
	return e.Id
}

func (e *StockEntry) IsCsOrderEntry() bool {
	return e.Type == EntryType1
}

func (e *StockEntry) IsTransferEntry() bool {
	return e.Type == EntryType0
}

// 新增出库单
func (e *StockEntry) InsertEntry(tx *gorm.DB, Type string) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	if _, err := e.GenerateEntryCode(tx); err != nil {
		return err
	}
	e.SetTypeForEntry(Type)
	e.SetProductsForEntry()
	//err := tx.Create(e).Error
	err := tx.Table(wcPrefix + "." + e.TableName()).Omit(clause.Associations).Create(e).Error
	if err != nil {
		return err
	}
	err = tx.Table(wcPrefix + ".stock_entry_products").Omit(clause.Associations).Create(e.StockEntryProducts).Error
	if err != nil {
		return err
	}
	return nil
}

// 新增出库单
func (e *StockEntry) InsertStockEntry(tx *gorm.DB, EntryCode string) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	err := tx.Table(wcPrefix+"."+e.TableName()).Omit(clause.Associations, "id").Create(e).Error
	if err != nil {
		return err
	}

	for _, product := range e.StockEntryProducts {
		err = tx.Table(wcPrefix+".stock_entry_products").Omit(clause.Associations, "id").Create(&product).Error
		if err != nil {
			return err
		}

		//sub信息
		for _, sub := range product.StockEntryProductsSub {
			sub.EntryProductId = product.Id
			err = tx.Table(wcPrefix+".stock_entry_products_sub").Omit(clause.Associations, "id").Create(&sub).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 根据类型设置出库单类型
func (e *StockEntry) SetTypeForEntry(Type string) {
	e.Type = Type
}

// 设置产品信息
func (e *StockEntry) SetProductsForEntry() {
	for index := range e.StockEntryProducts {
		e.StockEntryProducts[index].EntryCode = e.EntryCode
		e.StockEntryProducts[index].VendorId = e.VendorId
		e.StockEntryProducts[index].WarehouseCode = e.WarehouseCode
		e.StockEntryProducts[index].LogicWarehouseCode = e.LogicWarehouseCode
	}
}

// GenerateEntryCode
func (e *StockEntry) GenerateEntryCode(tx *gorm.DB) (string, error) {
	var count int64
	start, end := utils.GetTodayTime()
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)

	if err := tx.Table(wcPrefix+"."+e.TableName()).Where("created_at BETWEEN ? AND ?", start, end).Count(&count).Error; err != nil {
		return "", err
	}
	code := "IR" + time.Now().Format("20060102150405") + fmt.Sprintf("%04d", count+1)
	e.EntryCode = code
	return code, nil
}

// 根据入库单编号获取入库单
func (e *StockEntry) GetByEntryCode(tx *gorm.DB, code string) error {
	err := tx.Where("entry_code = ?", code).Take(e).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("入库单不存在")
	}
	return err
}

// 根据入库单来源编号获取入库单
func (e *StockEntry) GetByEntrySourceCode(tx *gorm.DB, sourceCode string) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	err := tx.Table(wcPrefix+"."+e.TableName()).Where("source_code = ?", sourceCode).Take(e).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("入库单不存在")
	}
	return err
}

// CheckBeforeConfirmEntry 确认出库前check
func (e *StockEntry) CheckBeforeConfirmEntry(tx *gorm.DB) error {
	if e.Status != EntryStatus1 {
		return errors.New("入库单状态不正确")
	}
	if e.Warehouse.CheckIsVirtual() {
		return errors.New("虚拟仓无需手动入库")
	}
	// 调拨单需先出库再入库
	if e.IsTransferEntry() {
		transfer := StockTransfer{}
		if err := transfer.GetByTransferCode(tx, e.SourceCode); err != nil {
			return err
		}
		if transfer.Status != TransferStatus2 {
			return errors.New("调拨类型的入库单需先出库再入库")
		}
	}
	return nil
}

// 根据入库单来源编号获取入库单options
func (e *StockEntry) GetByEntrySourceCodeWithOptions(tx *gorm.DB, sourceCode string, options func(tx *gorm.DB) *gorm.DB) error {
	if options != nil {
		tx = options(tx)
	}
	return e.GetByEntrySourceCode(tx, sourceCode)
}

// 手动入库
func (e *StockEntry) ConfirmEntryManual(tx *gorm.DB) error {
	stockEntryDefectiveProducts := make([]*StockEntryDefectiveProduct, 0)
	// 更改入库单状态
	if err := tx.Omit(clause.Associations).Save(e).Error; err != nil {
		return err
	}

	for _, item := range e.StockEntryProducts {
		// 次品退货入库商品 这里的库位是次品仓库位 需要自动选择入库库位
		if item.CheckIsDefective() {
			defectiveProduct, err := item.DefectiveHandle(tx)
			if err != nil {
				return err
			}
			stockEntryDefectiveProducts = append(stockEntryDefectiveProducts, defectiveProduct)
		}

		// 库位库存相关
		locationGoods := &StockLocationGoods{}
		if err := locationGoods.GetByGoodsIdAndLocationId(tx, item.GoodsId, item.StockLocationId); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				locationGoods.SetBaseInfo(item.GoodsId, item.StockLocationId, e.UpdateBy, e.UpdateByName)
			} else {
				return err
			}
		}
		if err := locationGoods.AddStock(tx, item.ActQuantity, &StockLocationGoodsLog{
			DocketCode:   e.EntryCode,
			FromType:     StockLocationGoodsLogFromType1,
			CreateBy:     e.UpdateBy,
			CreateByName: e.UpdateByName,
		}); err != nil {
			return err
		}

		// 库存相关
		stockInfo := &StockInfo{}
		if err := stockInfo.GetByGoodsIdAndLogicWarehouseCode(tx, item.GoodsId, item.LogicWarehouseCode); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				stockInfo.SetBaseInfo(item.GoodsId, item.VendorId, e.UpdateBy, item.SkuCode, item.WarehouseCode, item.LogicWarehouseCode, e.UpdateByName)
			} else {
				return err
			}
		}
		if err := stockInfo.AddStock(tx, item.ActQuantity, &StockLog{
			DocketCode:   e.EntryCode,
			FromType:     StockLogFromType1,
			CreateBy:     e.UpdateBy,
			CreateByName: e.UpdateByName,
		}); err != nil {
			return err
		}

		// 更新入库商品表实际数量(落库)
		if err := item.Save(tx); err != nil {
			return err
		}
	}

	// 创建调拨单For次品入库
	if len(stockEntryDefectiveProducts) > 0 {
		stockTransfer := &StockTransfer{}
		if err := stockTransfer.InsertTransferForDefectiveEntry(tx, e, stockEntryDefectiveProducts); err != nil {
			return err
		}
	}

	if e.IsTransferEntry() {
		transfer := &StockTransfer{}
		if err := transfer.GetByTransferCode(tx, e.SourceCode); err != nil {
			return err
		}
		if err := transfer.SetEntryCompleteStatus(tx); err != nil {
			return err
		}
	} else if e.IsCsOrderEntry() {
		// 通知售后出库完成 todo
		csApplyModel := modelsCs.CsApply{}
		if err := csApplyModel.ReturnCompleted(tx, e.SourceCode); err != nil {
			return err
		}
	}
	return nil
}

// 自动入库（虚拟仓）
func (e *StockEntry) ConfirmEntryAutoForVirtual(tx *gorm.DB) error {

	// 更改入库单状态
	if err := tx.Omit(clause.Associations).Save(e).Error; err != nil {
		return err
	}

	// 实际入库数量等于应入库数量,时间,状态
	e.SetProductsActQuantityComplete()

	for _, item := range e.StockEntryProducts {
		// 更新入库商品表实际数量
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
		if err := stockInfo.AddStockForVirtualWarehouse(tx, item.ActQuantity, &StockLog{
			DocketCode:   e.EntryCode,
			FromType:     StockLogFromType1,
			CreateBy:     e.UpdateBy,
			CreateByName: e.UpdateByName,
		}); err != nil {
			return err
		}
	}
	return nil
}

// 自动入库 设置实际入库数量为入库数量,时间,状态
func (e *StockEntry) SetProductsActQuantityComplete() {
	currTime := time.Now()
	for index, item := range e.StockEntryProducts {
		e.StockEntryProducts[index].EntryTime = currTime
		e.StockEntryProducts[index].EntryEndTime = currTime
		e.StockEntryProducts[index].Status = 2
		e.StockEntryProducts[index].ActQuantity = item.Quantity
	}
}

// 自动入库-入库单入库状态，时间设置
func (e *StockEntry) SetConfirmEntryStatus(updateBy int, UpdateByName string) {
	currTime := time.Now()
	e.Status = EntryStatus2
	e.EntryTime = currTime
	e.EntryEndTime = currTime
	e.UpdateBy = updateBy
	e.UpdateByName = UpdateByName
	for index := range e.StockEntryProducts {
		e.StockEntryProducts[index].UpdateBy = updateBy
		e.StockEntryProducts[index].UpdateByName = UpdateByName
	}
}

// 自动出库(次品入库创建的入库单)
func (e *StockEntry) ConfirmEntryForDefectiveEntry(tx *gorm.DB) error {
	for _, item := range e.StockEntryProducts {
		// 库位库存相关
		locationGoods := &StockLocationGoods{}
		if err := locationGoods.GetByGoodsIdAndLocationId(tx, item.GoodsId, item.StockLocationId); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				locationGoods.SetBaseInfo(item.GoodsId, item.StockLocationId, e.CreateBy, e.CreateByName)
			} else {
				return err
			}
		}
		if err := locationGoods.AddStock(tx, item.ActQuantity, &StockLocationGoodsLog{
			DocketCode:   e.EntryCode,
			FromType:     StockLocationGoodsLogFromType1,
			CreateBy:     e.CreateBy,
			CreateByName: e.CreateByName,
		}); err != nil {
			return err
		}
		// 库存相关
		stockInfo := &StockInfo{}
		if err := stockInfo.GetByGoodsIdAndLogicWarehouseCode(tx, item.GoodsId, item.LogicWarehouseCode); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				stockInfo.SetBaseInfo(item.GoodsId, item.VendorId, e.CreateBy, item.SkuCode, item.WarehouseCode, item.LogicWarehouseCode, e.CreateByName)
			} else {
				return err
			}
		}
		if err := stockInfo.AddStock(tx, item.ActQuantity, &StockLog{
			DocketCode:   e.EntryCode,
			FromType:     StockLogFromType1,
			CreateBy:     e.CreateBy,
			CreateByName: e.CreateByName,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (e *StockEntry) Stash(tx *gorm.DB) error {
	for _, item := range e.StockEntryProducts {
		if err := tx.Model(&item).Updates(map[string]interface{}{
			"stash_location_id":  item.StockLocationId,
			"stash_act_quantity": item.ActQuantity,
			"update_by":          item.UpdateBy,
			"update_by_name":     item.UpdateByName,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}
