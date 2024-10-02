package models

import (
	"errors"
	"fmt"
	"go-admin/common/utils"
	"strconv"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"go-admin/common/models"
)

const (
	ControlStatus0  = "0"
	ControlStatus1  = "1"
	ControlStatus2  = "2"
	ControlStatus99 = "99"

	ControlType0 = "0"
	ControlType1 = "1"
	ControlType2 = "2"

	ControlVerifyStatus0  = "0"
	ControlVerifyStatus1  = "1"
	ControlVerifyStatus2  = "2"
	ControlVerifyStatus99 = "99"

	ControlLogModelName   = "stockControl"
	ControlLogModelInsert = "insert"
	ControlLogModelUpdate = "update"
	ControlLogModelDelete = "delete"
	ControlLogModelAudit  = "audit"
)

var (
	ControlStatusMap = map[string]string{
		ControlStatus0:  "已作废",
		ControlStatus1:  "已提交",
		ControlStatus2:  "已完成",
		ControlStatus99: "未提交",
	}

	ControlTypeMap = map[string]string{
		ControlType0: "盘盈",
		ControlType1: "盘亏",
		ControlType2: "无差异",
	}

	ControlVerifyStatusMap = map[string]string{
		ControlVerifyStatus0:  "待审核",
		ControlVerifyStatus1:  "审核通过",
		ControlVerifyStatus2:  "审核驳回",
		ControlVerifyStatus99: "未审核",
	}
)

type StockControl struct {
	models.Model

	ControlCode        string    `json:"controlCode" gorm:"type:varchar(20);comment:调整单编码"`
	Type               string    `json:"type" gorm:"type:tinyint;comment:调整类型: 0 调增 1 调减"`
	Status             string    `json:"status" gorm:"type:tinyint;comment:状态:0-已作废 1-创建 2-已完成 99未提交"`
	WarehouseCode      string    `json:"warehouseCode" gorm:"type:varchar(20);comment:实体仓code"`
	LogicWarehouseCode string    `json:"logicWarehouseCode" gorm:"type:varchar(20);comment:逻辑仓code"`
	Remark             string    `json:"remark" gorm:"type:varchar(255);comment:备注"`
	VerifyStatus       string    `json:"verifyStatus" gorm:"type:tinyint;comment:审核状态 0 待审核 1 审核通过 2 审核驳回 99初始化"`
	VerifyRemark       string    `json:"verifyRemark" gorm:"type:varchar(255);comment:审核描述"`
	VendorId           int       `json:"vendorId" gorm:"type:int unsigned;comment:货主id"`
	VerifyUid          int       `json:"verifyUid" gorm:"type:int unsigned;comment:审核人id"`
	VerifyTime         time.Time `json:"verifyTime" gorm:"type:datetime;comment:审核时间"`
	models.ModelTime
	models.ControlBy
	Warehouse            Warehouse              `json:"-" gorm:"foreignKey:WarehouseCode;references:WarehouseCode"`
	LogicWarehouse       LogicWarehouse         `json:"-" gorm:"foreignKey:LogicWarehouseCode;references:LogicWarehouseCode"`
	Vendor               Vendors                `json:"-" gorm:"foreignKey:VendorId"`
	StockControlProducts []StockControlProducts `json:"-" gorm:"foreignKey:ControlCode;references:ControlCode"`
}

func (StockControl) TableName() string {
	return "stock_control"
}

func (e *StockControl) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *StockControl) GetId() interface{} {
	return e.Id
}

// 新增调整单

func (e *StockControl) InsertControl(tx *gorm.DB, action string) error {
	if err := e.InserOrUpdateCheckAndSet(tx, action); err != nil {
		return err
	}
	return tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(e).Error; err != nil {
			return err
		}
		return nil
	})
}

func (e *StockControl) InsertImmediately(tx *gorm.DB, action string) error {
	e.SetTypeForControl(action)
	return tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(e).Error; err != nil {
			return err
		}
		return nil
	})
}

func (e *StockControl) InserOrUpdateCheckAndSet(tx *gorm.DB, action string) error {
	// todo 仓库权限 货主权限
	lwh := &LogicWarehouse{}
	if err := lwh.GetWhAndLwhInfo(tx, e.WarehouseCode, e.LogicWarehouseCode); err != nil {
		return err
	}
	if lwh.Warehouse.CheckIsVirtual() {
		return errors.New("实体仓不能是虚拟仓")
	}
	e.SetTypeForControl(action)
	e.SetProductsForControl()
	if err := e.CheckProductsForControl(tx); err != nil {
		return err
	}
	return nil
}

// 编辑调整单

func (e *StockControl) UpdateControl(tx *gorm.DB, action string) error {
	if err := e.InserOrUpdateCheckAndSet(tx, action); err != nil {
		return err
	}
	return tx.Transaction(func(tx *gorm.DB) error {
		stockControlProducts := StockControlProducts{}
		if err := stockControlProducts.DeleteByControlCode(tx, e.ControlCode); err != nil {
			return err
		}
		if err := tx.Save(e).Error; err != nil {
			return err
		}
		return nil
	})
}

// GenerateControlCode

func (e *StockControl) GenerateControlCode(tx *gorm.DB) (string, error) {
	var count int64
	start, end := utils.GetTodayTime()
	if err := tx.Model(&StockControl{}).Where("created_at BETWEEN ? AND ?", start, end).Count(&count).Error; err != nil {
		return "", err
	}
	code := "AO" + time.Now().Format("20060102150405") + fmt.Sprintf("%04d", count+1)
	e.ControlCode = code
	return code, nil
}

// 根据类型设置出库单类型

func (e *StockControl) SetTypeForControl(action string) {
	switch action {
	case "Commit":
		e.VerifyStatus = ControlVerifyStatus0
		e.Status = ControlStatus1
	default:
		e.VerifyStatus = ControlVerifyStatus99
		e.Status = ControlStatus99
	}
}

// 设置产品信息
func (e *StockControl) SetProductsForControl() {
	for index := range e.StockControlProducts {
		e.StockControlProducts[index].ControlCode = e.ControlCode
		e.StockControlProducts[index].VendorId = e.VendorId
		e.StockControlProducts[index].WarehouseCode = e.WarehouseCode
		e.StockControlProducts[index].LogicWarehouseCode = e.LogicWarehouseCode
	}
}

// 检查产品

func (e *StockControl) CheckProductsForControl(tx *gorm.DB) error {
	tempSku := []string{}
	skuSlice := []string{}
	tempGoodsIdQuantity := &map[int]int{}
	if len(e.StockControlProducts) == 0 {
		return errors.New("商品明细不能为空")
	}
	for _, item := range e.StockControlProducts {
		if item.SkuCode == "" {
			return errors.New("缺少SKU")
		}
		if item.StockLocationId <= 0 {
			return errors.New(item.SkuCode + ",缺少库位编号")
		}
		existKey := item.SkuCode + "," + strconv.Itoa(item.StockLocationId)
		if existFlag := lo.Contains(tempSku, existKey); existFlag {
			return errors.New(item.SkuCode + ",重复的SKU+库位编号")
		}
		tempSku = append(tempSku, existKey)
		skuSlice = append(skuSlice, item.SkuCode)
		if item.CurrentQuantity < 0 {
			return errors.New("盘后数量为大于等于0的整数")
		}

		//if item.Quantity <= 0 {
		//	return errors.New("调整数量为大于0的整数")
		//}
	}
	if err := e.GetGoodsInfoBySku(tx, lo.Uniq(skuSlice)); err != nil {
		return err
	}

	// 检查产品库位
	if err := e.CheckProductsLocationForControl(tx, tempGoodsIdQuantity); err != nil {
		return err
	}

	// 检查产品库位库存
	if err := e.CheckProductsLocationStockForControl(tx); err != nil {
		return err
	}

	// 检查产品库存
	if err := e.CheckProductsStockForControl(tx, tempGoodsIdQuantity); err != nil {
		return err
	}
	return nil
}

// CheckProductsLocationForControl 检查产品库位
func (e *StockControl) CheckProductsLocationForControl(tx *gorm.DB, tempGoodsIdQuantity *map[int]int) error {
	for index, item := range e.StockControlProducts {
		stockLocation := &StockLocation{}
		if !stockLocation.CheckStockLocation(tx, item.StockLocationId, item.LogicWarehouseCode) {
			return errors.New(item.SkuCode + ",所选的库位不在当前逻辑仓库中")
		}
		e.StockControlProducts[index].StockLocation = *stockLocation

		quantity := 0
		if item.Type == ControlType1 {
			quantity = 0 - item.Quantity
		} else {
			quantity = item.Quantity
		}

		if _, ok := (*tempGoodsIdQuantity)[item.GoodsId]; ok {
			(*tempGoodsIdQuantity)[item.GoodsId] += quantity
		} else {
			(*tempGoodsIdQuantity)[item.GoodsId] = quantity
		}
	}
	return nil
}

// 检查产品库位库存

func (e *StockControl) CheckProductsLocationStockForControl(tx *gorm.DB) error {
	for key, item := range e.StockControlProducts {
		locationGoods := &StockLocationGoods{}
		if !locationGoods.CheckStockByGoodsIdAndLocationId(tx, item) {
			return errors.New(item.SkuCode + ",库位:" + item.StockLocation.LocationCode + ",库位库存不足")
		}

		// 盘亏
		if locationGoods.LockStock+locationGoods.Stock > e.StockControlProducts[key].CurrentQuantity {
			e.StockControlProducts[key].Type = "1"
		} else if locationGoods.LockStock+locationGoods.Stock == e.StockControlProducts[key].CurrentQuantity {
			e.StockControlProducts[key].Type = "2"
		} else { // 盘盈
			e.StockControlProducts[key].Type = "0"
		}
		e.StockControlProducts[key].Quantity = utils.AbsToInt(locationGoods.LockStock + locationGoods.Stock - e.StockControlProducts[key].CurrentQuantity)
	}
	return nil
}

// CheckProductsStockForControl 检查产品库存
func (e *StockControl) CheckProductsStockForControl(tx *gorm.DB, tempGoodsIdQuantity *map[int]int) error {
	for goodsId, quantity := range *tempGoodsIdQuantity {
		if quantity < 0 {
			stockInfo := &StockInfo{}
			if !stockInfo.CheckQuantityByGoodsIdAndLogicWarehouseCode(tx, quantity, goodsId, e.LogicWarehouseCode) {
				return errors.New(e.GetGoodsInfoByGoodsIdFromProducts(goodsId).SkuCode + ",库存不足")
			}
		}
	}
	return nil
}

func (e *StockControl) GetGoodsInfoByGoodsIdFromProducts(goodsId int) StockControlProducts {
	for _, item := range e.StockControlProducts {
		if item.GoodsId == goodsId {
			return item
		}
	}
	return StockControlProducts{}
}

// 获取Goods信息

func (e *StockControl) GetGoodsInfoBySku(tx *gorm.DB, skuSlice []string) error {
	skuGoodsMap, skuGoodsSlice := GetGoodsInfoMapByThreeFromPcClient(tx, skuSlice, e.WarehouseCode, e.VendorId, 1)
	skuGoodsErrSlice := lo.Without(skuSlice, skuGoodsSlice...)
	if len(skuGoodsErrSlice) != 0 {
		return errors.New("商品关系存在异常")

	}
	for index, item := range e.StockControlProducts {
		e.StockControlProducts[index].GoodsId = skuGoodsMap[item.SkuCode].Id
	}
	return nil
}

// 调整单审核驳回

func (e *StockControl) AuditReject(tx *gorm.DB) error {
	e.SetControlStatus0()
	return tx.Omit(clause.Associations).Save(e).Error
}

func (e *StockControl) SetControlStatus0() {
	e.Status = ControlStatus0
}

func (e *StockControl) SetControlStatus2() {
	e.Status = ControlStatus2
}

// 调整单审核通过

func (e *StockControl) AuditPass(tx *gorm.DB) error {
	e.SetControlStatus2()
	return tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit(clause.Associations).Save(e).Error; err != nil {
			return err
		}
		if err := e.AuditPassForControl(tx); err != nil {
			return err
		}

		return nil
	})
}

// AuditPassForControl 调整单审核驳回库存、库位库存操作
func (e *StockControl) AuditPassForControl(tx *gorm.DB) error {
	tempGoodsIdQuantity := map[int]int{}
	for _, item := range e.StockControlProducts {
		quantity := 0
		if item.Type == ControlType1 {
			quantity = 0 - item.Quantity
		} else {
			quantity = item.Quantity
		}
		if _, ok := (tempGoodsIdQuantity)[item.GoodsId]; ok {
			(tempGoodsIdQuantity)[item.GoodsId] += quantity
		} else {
			(tempGoodsIdQuantity)[item.GoodsId] = quantity
		}
	}
	for goodsId, quantity := range tempGoodsIdQuantity {
		stockInfo := &StockInfo{}
		if quantity < 0 {
			quantity = utils.AbsToInt(quantity)
			if !stockInfo.CheckQuantityByGoodsIdAndLogicWarehouseCode(tx, quantity, goodsId, e.LogicWarehouseCode) {
				return errors.New(e.GetGoodsInfoByGoodsIdFromProducts(goodsId).SkuCode + ",库存不足")
			}
			if err := stockInfo.SubStock(tx, quantity, &StockLog{
				DocketCode:   e.ControlCode,
				FromType:     StockLogFromType2,
				CreateBy:     e.UpdateBy,
				CreateByName: e.UpdateByName,
			}); err != nil {
				return err
			}
		} else {
			for _, sitem := range e.StockControlProducts {
				if sitem.GoodsId == goodsId {
					if err := stockInfo.GetByGoodsIdAndLogicWarehouseCode(tx, sitem.GoodsId, sitem.LogicWarehouseCode); err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							stockInfo.SetBaseInfo(sitem.GoodsId, sitem.VendorId, e.UpdateBy, sitem.SkuCode, sitem.WarehouseCode, sitem.LogicWarehouseCode, e.UpdateByName)
						} else {
							return err
						}
					}
					break
				}
			}

			if stockInfo.GoodsId != 0 {
				if err := stockInfo.AddStock(tx, quantity, &StockLog{
					DocketCode:   e.ControlCode,
					FromType:     StockLogFromType2,
					CreateBy:     e.UpdateBy,
					CreateByName: e.UpdateByName,
				}); err != nil {
					return err
				}
			}
		}

	}

	for _, item := range e.StockControlProducts {
		// 库位库存相关
		locationGoods := &StockLocationGoods{}

		if item.Type == ControlType1 {
			if err := locationGoods.GetByGoodsIdAndLocationId(tx, item.GoodsId, item.StockLocationId); err != nil {
				return err
			}
			if !locationGoods.CheckStock(item.CurrentQuantity) {
				return errors.New(item.SkuCode + "库位:" + item.StockLocation.LocationCode + ",库位库存不足")
			}
			if err := locationGoods.SubStock(tx, item.Quantity, &StockLocationGoodsLog{
				DocketCode:   e.ControlCode,
				FromType:     StockLocationGoodsLogFromType2,
				CreateBy:     e.UpdateBy,
				CreateByName: e.UpdateByName,
			}); err != nil {
				return err
			}
		} else {
			if err := locationGoods.GetByGoodsIdAndLocationId(tx, item.GoodsId, item.StockLocationId); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					locationGoods.SetBaseInfo(item.GoodsId, item.StockLocationId, e.UpdateBy, e.UpdateByName)
				} else {
					return err
				}
			}
			if err := locationGoods.AddStock(tx, item.Quantity, &StockLocationGoodsLog{
				DocketCode:   e.ControlCode,
				FromType:     StockLocationGoodsLogFromType2,
				CreateBy:     e.UpdateBy,
				CreateByName: e.UpdateByName,
			}); err != nil {
				return err
			}
		}

	}
	return nil
}
