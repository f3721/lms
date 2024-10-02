package models

import (
	"errors"
	"fmt"
	modelsOc "go-admin/app/oc/models"
	modelsPc "go-admin/app/pc/models"
	dtoPc "go-admin/app/pc/service/admin/dto"
	dtoUc "go-admin/app/uc/service/admin/dto"
	ocClient "go-admin/common/client/oc"
	pcClient "go-admin/common/client/pc"
	ucClient "go-admin/common/client/uc"
	"go-admin/common/utils"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"go-admin/common/models"
)

const (
	TransferStatus0  = "0"
	TransferStatus1  = "1"
	TransferStatus2  = "2"
	TransferStatus3  = "3"
	TransferStatus4  = "4"
	TransferStatus5  = "5"
	TransferStatus6  = "6"
	TransferStatus98 = "98"
	TransferStatus99 = "99"

	TransferType0 = "0"
	TransferType1 = "1"
	TransferType2 = "2"
	TransferType3 = "3"

	TransferVerifyStatus0  = "0"
	TransferVerifyStatus1  = "1"
	TransferVerifyStatus2  = "2"
	TransferVerifyStatus99 = "99"

	TransferLogModelName   = "stockTransfer"
	TransferLogModelInsert = "insert"
	TransferLogModelUpdate = "update"
	TransferLogModelDelete = "delete"
	TransferLogModelAudit  = "audit"

	TransferRemarkDefectiveEntry = "次品退货入库生成的自动调拨单相关"
)

var (
	TransferStatusMap = map[string]string{
		TransferStatus0:  "已作废",
		TransferStatus1:  "未出库未入库",
		TransferStatus2:  "已出库未入库",
		TransferStatus3:  "出库完成入库完成",
		TransferStatus4:  "部分出库未入库",
		TransferStatus5:  "部分出库部分入库",
		TransferStatus6:  "全部出库部分入库",
		TransferStatus98: "已提交",
		TransferStatus99: "待提交",
	}

	TransferTypeMap = map[string]string{
		TransferType0: "正品调拨",
		TransferType1: "次品调拨",
		TransferType2: "正转次调拨",
		TransferType3: "次转正调拨",
	}

	TransferVerifyStatusMap = map[string]string{
		TransferVerifyStatus0:  "待审核",
		TransferVerifyStatus1:  "审核通过",
		TransferVerifyStatus2:  "审核驳回",
		TransferVerifyStatus99: "无",
	}
	TransferAutoInMap = map[int]string{
		0: "否",
		1: "是",
	}
	TransferAutoOutMap = map[int]string{
		0: "否",
		1: "是",
	}
)

type StockTransfer struct {
	models.Model

	TransferCode           string    `json:"transferCode" gorm:"type:varchar(32);comment:调拨单编码"`
	Type                   string    `json:"type" gorm:"type:tinyint;comment:调拨类型: 0 -正常调拨 1-次品调拨 2-正转次调拨 3-次转正调拨"`
	Status                 string    `json:"status" gorm:"type:tinyint;comment:状态:0-已作废 1-未出库未入库 2-已出库未入库 4-部分出库未入库 5-部分出库部分入库 6-全部出库部分入库  3-出库完成入库完成 98已提交 99未提交"`
	FromWarehouseCode      string    `json:"fromWarehouseCode" gorm:"type:varchar(20);comment:出库实体仓code"`
	FromLogicWarehouseCode string    `json:"fromLogicWarehouseCode" gorm:"type:varchar(20);comment:出库逻辑仓code"`
	ToWarehouseCode        string    `json:"toWarehouseCode" gorm:"type:varchar(20);comment:入库实体仓code"`
	ToLogicWarehouseCode   string    `json:"toLogicWarehouseCode" gorm:"type:varchar(20);comment:入库逻辑仓code"`
	Remark                 string    `json:"remark" gorm:"type:varchar(255);comment:备注"`
	LogisticsRemark        string    `json:"logisticsRemark" gorm:"type:varchar(255);comment:物流备注"`
	Mobile                 string    `json:"mobile" gorm:"type:varchar(20);comment:Mobile"`
	Linkman                string    `json:"linkman" gorm:"type:varchar(20);comment:联系人"`
	Address                string    `json:"address" gorm:"type:varchar(100);comment:详细地址"`
	District               int       `json:"district" gorm:"type:int;comment:区"`
	City                   int       `json:"city" gorm:"type:int;comment:市"`
	Province               int       `json:"province" gorm:"type:int;comment:省"`
	SourceCode             string    `json:"sourceCode" gorm:"type:varchar(100);comment:来源单号"`
	AutoIn                 int       `json:"autoIn" gorm:"type:tinyint(1);comment:自动入库 0 否  1是"`
	AutoOut                int       `json:"autoOut" gorm:"type:tinyint(1);comment:自动出库 0 否 1 是"`
	VerifyStatus           string    `json:"verifyStatus" gorm:"type:tinyint;comment:审核状态 0 待审核 1 审核通过 2 审核驳回 99初始化"`
	VerifyRemark           string    `json:"verifyRemark" gorm:"type:varchar(255);comment:审核描述"`
	VendorId               int       `json:"vendorId" gorm:"type:int unsigned;comment:货主id"`
	VerifyUid              int       `json:"verifyUid" gorm:"type:int unsigned;comment:审核人id"`
	VerifyTime             time.Time `json:"verifyTime" gorm:"type:datetime;comment:审核时间"`
	models.RegionName
	models.ModelTime
	models.ControlBy

	FromWarehouse         Warehouse               `json:"-" gorm:"foreignKey:FromWarehouseCode;references:WarehouseCode"`
	ToWarehouse           Warehouse               `json:"-" gorm:"foreignKey:ToWarehouseCode;references:WarehouseCode"`
	FromLogicWarehouse    LogicWarehouse          `json:"-" gorm:"foreignKey:FromLogicWarehouseCode;references:LogicWarehouseCode"`
	ToLogicWarehouse      LogicWarehouse          `json:"-" gorm:"foreignKey:ToLogicWarehouseCode;references:LogicWarehouseCode"`
	Vendor                Vendors                 `json:"-" gorm:"foreignKey:VendorId"`
	StockTransferProducts []StockTransferProducts `json:"-" gorm:"foreignKey:TransferCode;references:TransferCode"`
}

func (StockTransfer) TableName() string {
	return "stock_transfer"
}

func (e *StockTransfer) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *StockTransfer) GetId() interface{} {
	return e.Id
}

// 生成调拨单编码
func (e *StockTransfer) GenerateTransferCode(tx *gorm.DB) (string, error) {
	var count int64
	start, end := utils.GetTodayTime()
	if err := tx.Model(&StockTransfer{}).Where("created_at BETWEEN ? AND ?", start, end).Count(&count).Error; err != nil {
		return "", err
	}
	code := "TF" + time.Now().Format("20060102150405") + fmt.Sprintf("%04d", count+1)
	e.TransferCode = code
	return code, nil
}

// 新增调拨单
func (e *StockTransfer) InsertTransfer(tx *gorm.DB, action string) error {
	formLwh, _, err := e.InserOrUpdateCheckAndSet(tx, action)
	if err != nil {
		return err
	}
	//开启事务
	return tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(e).Error; err != nil {
			return err
		}
		if err := e.CommitTransferLockStock(tx, formLwh); err != nil {
			return err
		}
		return nil
	})
}

// 调拨单提交时锁库
func (e *StockTransfer) CommitTransferLockStock(tx *gorm.DB, formLwh *LogicWarehouse) error {
	if !formLwh.Warehouse.CheckIsVirtual() && e.Status == TransferStatus98 {
		for _, item := range e.StockTransferProducts {
			stockInfo := &StockInfo{}
			if err := stockInfo.GetByGoodsIdAndLogicWarehouseCode(tx, item.GoodsId, item.LogicWarehouseCode); err != nil {
				return err
			}
			if err := stockInfo.Lock(tx, item.Quantity, &StockLockLog{
				DocketCode: e.TransferCode,
				FromType:   StockLockLogFromType0,
				Remark:     RemarkStockLockLogTransferCommit,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

// 新增、编辑调拨单公共check
func (e *StockTransfer) InserOrUpdateCheckAndSet(tx *gorm.DB, action string) (formLwh, toLwh *LogicWarehouse, err error) {
	//判断货主仓库权限 todo
	formLwh = &LogicWarehouse{}
	if err = formLwh.GetWhAndLwhInfo(tx, e.FromWarehouseCode, e.FromLogicWarehouseCode); err != nil {
		return
	}
	toLwh = &LogicWarehouse{}
	if err = toLwh.GetWhAndLwhInfo(tx, e.ToWarehouseCode, e.ToLogicWarehouseCode); err != nil {
		return
	}
	e.SetStatusForTransfer(action)
	e.SetProductsForTransfer()
	if err = e.CheckWarehouseByTypeForTransfer(formLwh, toLwh); err != nil {
		return
	}
	if err = e.CheckProductsForTransfer(tx, formLwh, toLwh); err != nil {
		return
	}

	return
}

// 编辑调拨单
func (e *StockTransfer) UpdateTransfer(tx *gorm.DB, action string) error {
	formLwh, _, err := e.InserOrUpdateCheckAndSet(tx, action)
	if err != nil {
		return err
	}
	//开启事务
	return tx.Transaction(func(tx *gorm.DB) error {
		stockTransferProducts := StockTransferProducts{}
		if err := stockTransferProducts.DeleteByTransferCode(tx, e.TransferCode); err != nil {
			return err
		}
		if err := tx.Save(e).Error; err != nil {
			return err
		}
		if err := e.CommitTransferLockStock(tx, formLwh); err != nil {
			return err
		}
		return nil
	})
}

// 设置产品信息
func (e *StockTransfer) SetProductsForTransfer() {
	for index := range e.StockTransferProducts {
		e.StockTransferProducts[index].TransferCode = e.TransferCode
		e.StockTransferProducts[index].VendorId = e.VendorId
		e.StockTransferProducts[index].WarehouseCode = e.FromWarehouseCode
		e.StockTransferProducts[index].LogicWarehouseCode = e.FromLogicWarehouseCode
		e.StockTransferProducts[index].ToWarehouseCode = e.ToWarehouseCode
		e.StockTransferProducts[index].ToLogicWarehouseCode = e.ToLogicWarehouseCode
	}
}

func (e *StockTransfer) SetStatusForTransfer(action string) {

	switch action {
	case "Commit":
		e.VerifyStatus = TransferVerifyStatus0
		e.Status = TransferStatus98
	case "external":
		e.VerifyStatus = TransferVerifyStatus1
		e.Status = TransferStatus2
	default:
		e.VerifyStatus = TransferVerifyStatus99
		e.Status = TransferStatus99
	}
}

// 获取Goods信息
func (e *StockTransfer) SetGoodsIdByFromToMap(fromMap map[string]modelsPc.Goods, toMap map[string]modelsPc.Goods) error {
	for index, item := range e.StockTransferProducts {
		goodsInfo, ok := fromMap[item.SkuCode]
		goodsInfoTo, okTo := toMap[item.SkuCode]
		fromGoodsId := 0
		toGoodsId := 0
		if !ok && !okTo {
			return errors.New("SetGoodsIdByFromToMap Error")
		}
		if ok && !okTo {
			fromGoodsId = goodsInfo.Id
			toGoodsId = fromGoodsId
		}
		if !ok && okTo {
			toGoodsId = goodsInfoTo.Id
			fromGoodsId = toGoodsId
		}
		if ok && okTo {
			fromGoodsId = goodsInfo.Id
			toGoodsId = goodsInfoTo.Id
		}
		e.StockTransferProducts[index].GoodsId = fromGoodsId
		e.StockTransferProducts[index].ToGoodsId = toGoodsId
	}
	return nil
}

// 为调拨单检查产品
func (e *StockTransfer) CheckProductsForTransfer(tx *gorm.DB, formLwh, toLwh *LogicWarehouse) error {
	tempSku := []string{}
	skuGoodsFromMap := map[string]modelsPc.Goods{}
	skuGoodsFromSlice := []string{}
	skuGoodsToMap := map[string]modelsPc.Goods{}
	skuGoodsToSlice := []string{}

	if len(e.StockTransferProducts) == 0 {
		return errors.New("商品明细不能为空")
	}
	for _, item := range e.StockTransferProducts {
		if item.SkuCode == "" {
			return errors.New("缺少SKU")
		}
		if existSkuFlag := lo.Contains(tempSku, item.SkuCode); existSkuFlag {
			return errors.New(item.SkuCode + "，重复的SKU")
		}
		tempSku = append(tempSku, item.SkuCode)
		if item.Quantity <= 0 {
			return errors.New("调拨数量为大于0的整数")
		}
	}
	if !formLwh.Warehouse.CheckIsVirtual() {
		skuGoodsFromMap, skuGoodsFromSlice = GetGoodsInfoMapByThreeFromPcClient(tx, tempSku, e.FromWarehouseCode, e.VendorId, 1)
		var skuGoodsFromErrSlice []string = lo.Without(tempSku, skuGoodsFromSlice...)
		if len(skuGoodsFromErrSlice) != 0 {
			return errors.New("出库仓的商品关系存在异常")

		}
	} else {
		e.AutoOut = 1
	}
	if !toLwh.Warehouse.CheckIsVirtual() {
		skuGoodsToMap, skuGoodsToSlice = GetGoodsInfoMapByThreeFromPcClient(tx, tempSku, e.ToWarehouseCode, e.VendorId, 1)
		var skuGoodsToErrSlice []string = lo.Without(tempSku, skuGoodsToSlice...)
		if len(skuGoodsToErrSlice) != 0 {
			return errors.New("入库仓的商品关系存在异常")

		}
	} else {
		e.AutoIn = 1
	}
	_ = e.SetGoodsIdByFromToMap(skuGoodsFromMap, skuGoodsToMap)

	// 检查产品库存
	if err := e.CheckProductsStockForTransfer(tx, formLwh); err != nil {
		return err
	}
	return nil
}

// 检查产品库存
func (e *StockTransfer) CheckProductsStockForTransfer(tx *gorm.DB, formLwh *LogicWarehouse) error {
	if !formLwh.Warehouse.CheckIsVirtual() {
		for _, item := range e.StockTransferProducts {
			stockInfo := &StockInfo{}
			if !stockInfo.CheckStockByGoodsIdAndLogicWarehouseCode(tx, item.Quantity, item.GoodsId, item.LogicWarehouseCode) {
				return errors.New(item.SkuCode + ",库存不足")
			}
		}
	}
	return nil
}

// 新增调拨单时检查仓库信息
func (e *StockTransfer) CheckWarehouseByTypeForTransfer(formLwh, toLwh *LogicWarehouse) error {
	if formLwh.Warehouse.CheckIsVirtual() && toLwh.Warehouse.CheckIsVirtual() {
		return errors.New("出库实体仓和入库实体仓不能同时为虚拟仓")
	}
	switch e.Type {
	case TransferType0:
		if formLwh.Warehouse.WarehouseCode == toLwh.Warehouse.WarehouseCode {
			return errors.New("出库实体仓和入库实体仓不能一样")
		}
		if formLwh.Type != "0" {
			return errors.New("出库逻辑仓不是正品仓")
		}
		if toLwh.Type != "0" {
			return errors.New("入库逻辑仓不是正品仓")
		}
	case TransferType1:
		if formLwh.Warehouse.WarehouseCode == toLwh.Warehouse.WarehouseCode {
			return errors.New("出库实体仓和入库实体仓不能一样")
		}
		if formLwh.Type != "1" {
			return errors.New("出库逻辑仓不是次品仓")
		}
		if toLwh.Type != "1" {
			return errors.New("入库逻辑仓不是次品仓")
		}
	case TransferType2:
		if formLwh.Warehouse.WarehouseCode != toLwh.Warehouse.WarehouseCode {
			return errors.New("出库实体仓和入库实体仓不一样")
		}
		if formLwh.Type != "0" {
			return errors.New("出库逻辑仓不是正品仓")
		}
		if toLwh.Type != "1" {
			return errors.New("入库逻辑仓不是次品仓")
		}
	case TransferType3:
		if formLwh.Warehouse.WarehouseCode != toLwh.Warehouse.WarehouseCode {
			return errors.New("出库实体仓和入库实体仓不一样")
		}
		if formLwh.Type != "1" {
			return errors.New("出库逻辑仓不是次品仓")
		}
		if toLwh.Type != "0" {
			return errors.New("入库逻辑仓不是正品仓")
		}
	}
	return nil
}

// 调拨单审核驳回
func (e *StockTransfer) AuditReject(tx *gorm.DB) error {
	e.SetOutboundStatus0()
	return tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit(clause.Associations).Save(e).Error; err != nil {
			return err
		}

		//非虚拟仓解锁库存
		if !e.FromWarehouse.CheckIsVirtual() {
			for _, item := range e.StockTransferProducts {
				stockInfo := &StockInfo{}
				if err := stockInfo.GetByGoodsIdAndLogicWarehouseCode(tx, item.GoodsId, item.LogicWarehouseCode); err != nil {
					return err
				}
				if err := stockInfo.Unlock(tx, item.Quantity, &StockLockLog{
					DocketCode: e.TransferCode,
					FromType:   StockLockLogFromType0,
					Remark:     RemarkStockLockLogTransferAuditReject,
				}); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// 通过code获取调拨单
func (e *StockTransfer) GetByTransferCode(tx *gorm.DB, code string) error {
	err := tx.Where("transfer_code = ?", code).Take(e).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("调拨单不存在")
	}
	return err
}

// 通过code获取调拨单带关联查询
func (e *StockTransfer) GetByTransferCodeWithOptions(tx *gorm.DB, code string, options func(tx *gorm.DB) *gorm.DB) error {
	if options != nil {
		tx = options(tx)
	}
	return e.GetByTransferCode(tx, code)
}

// 调拨单废弃状态设置
func (e *StockTransfer) SetOutboundStatus0() {
	e.Status = TransferStatus0
}

// 调拨单审核通过状态设置
func (e *StockTransfer) SetOutboundStatus1() {
	e.Status = TransferStatus1
}

// 调拨单出库完成状态设置
func (e *StockTransfer) SetOutboundCompleteStatus(tx *gorm.DB) error {
	e.Status = TransferStatus2
	return tx.Omit(clause.Associations).Save(e).Error
}

// 调拨单入库完成状态设置
func (e *StockTransfer) SetEntryCompleteStatus(tx *gorm.DB) error {
	e.Status = TransferStatus3
	return tx.Omit(clause.Associations).Save(e).Error
}

// 新增调拨单For次品入库
func (e *StockTransfer) InsertTransferForDefectiveEntry(tx *gorm.DB, stockEntry *StockEntry, defectiveProducts []*StockEntryDefectiveProduct) error {
	stockOutbound := &StockOutbound{}
	newStockEntry := &StockEntry{}

	// 获取入库单逻辑仓对应次品仓
	defectiveLogicWarehouse := &LogicWarehouse{}
	if err := defectiveLogicWarehouse.GetDefectiveOrPassedLogicWarehouse(tx, stockEntry.LogicWarehouseCode, LwhType1); err != nil {
		return errors.New("入库单逻辑仓对应次品仓获取失败")
	}
	if err := e.initTransferDataForDefectiveEntry(tx, stockEntry, defectiveLogicWarehouse); err != nil {
		return err
	}
	e.initOutboundDataForDefectiveEntry(stockOutbound)
	e.initEntryDataForDefectiveEntry(newStockEntry)

	curTime := time.Now()
	for _, item := range defectiveProducts {
		e.StockTransferProducts = append(e.StockTransferProducts, StockTransferProducts{
			TransferCode:         e.TransferCode,
			SkuCode:              item.SkuCode,
			Quantity:             item.ActQuantity,
			VendorId:             item.VendorId,
			WarehouseCode:        item.WarehouseCode,
			LogicWarehouseCode:   item.LogicWarehouseCode,
			GoodsId:              item.GoodsId,
			ToWarehouseCode:      item.WarehouseCode,
			ToLogicWarehouseCode: defectiveLogicWarehouse.LogicWarehouseCode,
			ToGoodsId:            item.GoodsId,
			ControlBy: models.ControlBy{
				CreateBy:     stockEntry.UpdateBy,
				CreateByName: stockEntry.UpdateByName,
			},
		})
		stockOutbound.StockOutboundProducts = append(stockOutbound.StockOutboundProducts, StockOutboundProducts{
			SkuCode:         item.SkuCode,
			Quantity:        item.ActQuantity,
			ActQuantity:     item.ActQuantity,
			GoodsId:         item.GoodsId,
			StockLocationId: item.PassedStockLocationId,
			ControlBy: models.ControlBy{
				CreateBy:     stockEntry.UpdateBy,
				CreateByName: stockEntry.UpdateByName,
			},
		})
		newStockEntry.StockEntryProducts = append(newStockEntry.StockEntryProducts, StockEntryProducts{
			SkuCode:         item.SkuCode,
			Quantity:        item.ActQuantity,
			ActQuantity:     item.ActQuantity,
			GoodsId:         item.GoodsId,
			IsDefective:     item.IsDefective,
			EntryTime:       curTime,
			EntryEndTime:    curTime,
			StockLocationId: item.DefectiveStockLocationId,
			ControlBy: models.ControlBy{
				CreateBy:     stockEntry.UpdateBy,
				CreateByName: stockEntry.UpdateByName,
			},
		})
	}

	// 落库调拨单
	if err := tx.Create(e).Error; err != nil {
		return err
	}
	// 落库出库单
	if err := stockOutbound.InsertOutbound(tx, OutboundType0); err != nil {
		return err
	}
	// 落库入库单
	if err := newStockEntry.InsertEntry(tx, EntryType0); err != nil {
		return err
	}
	// 自动出库
	if err := stockOutbound.ConfirmOutboundForDefectiveEntry(tx); err != nil {
		return err
	}
	// 自动入库
	if err := newStockEntry.ConfirmEntryForDefectiveEntry(tx); err != nil {
		return err
	}
	return nil
}

// 初始化调拨单数据For次品入库
func (e *StockTransfer) initTransferDataForDefectiveEntry(tx *gorm.DB, stockEntry *StockEntry, defectiveLogicWarehouse *LogicWarehouse) error {
	if _, err := e.GenerateTransferCode(tx); err != nil {
		return err
	}
	e.Type = TransferType2
	e.Status = TransferStatus3
	e.FromWarehouseCode = stockEntry.WarehouseCode
	e.FromLogicWarehouseCode = stockEntry.LogicWarehouseCode
	e.ToWarehouseCode = stockEntry.WarehouseCode
	e.ToLogicWarehouseCode = defectiveLogicWarehouse.LogicWarehouseCode
	e.Remark = stockEntry.EntryCode + "," + TransferRemarkDefectiveEntry
	e.Mobile = stockEntry.Warehouse.Mobile
	e.Linkman = stockEntry.Warehouse.Linkman
	e.Address = stockEntry.Warehouse.Address
	e.District = stockEntry.Warehouse.District
	e.DistrictName = stockEntry.Warehouse.DistrictName
	e.City = stockEntry.Warehouse.City
	e.CityName = stockEntry.Warehouse.CityName
	e.Province = stockEntry.Warehouse.Province
	e.ProvinceName = stockEntry.Warehouse.ProvinceName
	e.SourceCode = stockEntry.EntryCode
	e.AutoIn = 1
	e.AutoOut = 1
	e.VerifyStatus = TransferVerifyStatus1
	e.VendorId = stockEntry.VendorId
	e.VerifyTime = time.Now()
	e.CreateBy = stockEntry.UpdateBy
	e.CreateByName = stockEntry.UpdateByName
	return nil
}

// 初始化出库单数据For次品入库
func (e *StockTransfer) initOutboundDataForDefectiveEntry(stockOutbound *StockOutbound) {
	currTime := time.Now()
	stockOutbound.Status = OutboundStatus2
	stockOutbound.SourceCode = e.TransferCode
	stockOutbound.Remark = TransferRemarkDefectiveEntry
	stockOutbound.OutboundTime = currTime
	stockOutbound.OutboundEndTime = currTime
	stockOutbound.WarehouseCode = e.FromWarehouseCode
	stockOutbound.LogicWarehouseCode = e.FromLogicWarehouseCode
	stockOutbound.VendorId = e.VendorId
	stockOutbound.CreateBy = e.CreateBy
	stockOutbound.CreateByName = e.CreateByName
}

// 初始化入库单数据For次品入库
func (e *StockTransfer) initEntryDataForDefectiveEntry(stockEntry *StockEntry) {
	currTime := time.Now()
	stockEntry.Status = EntryStatus2
	stockEntry.SourceCode = e.TransferCode
	stockEntry.Remark = TransferRemarkDefectiveEntry
	stockEntry.EntryTime = currTime
	stockEntry.EntryEndTime = currTime
	stockEntry.WarehouseCode = e.ToWarehouseCode
	stockEntry.LogicWarehouseCode = e.ToLogicWarehouseCode
	stockEntry.VendorId = e.VendorId
	stockEntry.CreateBy = e.CreateBy
	stockEntry.CreateByName = e.CreateByName
}

// 通过sku从pc获取sku信息
func GetSkuMapFromPcClient(tx *gorm.DB, skuSlice []string) (skuMap map[string]dtoPc.InnerGetProductBySkuResp, skuSliceRes []string) {
	result := pcClient.ApiByDbContext(tx).GetProductBySku(skuSlice)
	resultInfo := &struct {
		response.Response
		Data []dtoPc.InnerGetProductBySkuResp
	}{}
	result.Scan(resultInfo)
	skuSliceRes = lo.Map(resultInfo.Data, func(f dtoPc.InnerGetProductBySkuResp, _ int) string {
		return f.SkuCode
	})
	skuMap = lo.Associate(resultInfo.Data, func(f dtoPc.InnerGetProductBySkuResp) (string, dtoPc.InnerGetProductBySkuResp) {
		return f.SkuCode, f
	})
	return
}

// 通过仓库、货主、sku从pc获取goods信息
func GetGoodsInfoMapByThreeFromPcClient(tx *gorm.DB, skuSlice []string, whCode string, vendorId, approveStatus int) (skuGoodsMap map[string]modelsPc.Goods, skuGoodsSlice []string) {
	result := pcClient.ApiByDbContext(tx).GetGoodsBySkuAndVendorAndWarehouse(skuSlice, whCode, vendorId, approveStatus)
	resultInfo := &struct {
		response.Response
		Data []modelsPc.Goods
	}{}
	result.Scan(resultInfo)
	skuGoodsSlice = lo.Map(resultInfo.Data, func(f modelsPc.Goods, _ int) string {
		return f.SkuCode
	})
	skuGoodsMap = lo.Associate(resultInfo.Data, func(f modelsPc.Goods) (string, modelsPc.Goods) {
		return f.SkuCode, f
	})
	return
}

// 获取goods、product信息通过pc
func GetGoodsProductInfoFromPc(tx *gorm.DB, goodsIdSlice []int) map[int]modelsPc.Goods {
	result := pcClient.ApiByDbContext(tx).GetGoodsById(goodsIdSlice)
	resultInfo := &struct {
		response.Response
		Data []modelsPc.Goods
	}{}
	result.Scan(resultInfo)

	return lo.Associate(resultInfo.Data, func(f modelsPc.Goods) (int, modelsPc.Goods) {
		return f.Id, f
	})
}

// 获取订单信息通过oc
func GeOrderInfoFromOc(tx *gorm.DB, orderId string) modelsOc.OrderInfo {
	result := ocClient.ApiByDbContext(tx).GetByOrderId(orderId)
	resultInfo := &struct {
		response.Response
		Data modelsOc.OrderInfo
	}{}
	result.Scan(resultInfo)
	return resultInfo.Data
}

func GeOrderInfoFromOcByCsNo(tx *gorm.DB, csNo string) modelsOc.OrderInfo {
	result := ocClient.ApiByDbContext(tx).GetByCsNo(csNo)
	resultInfo := &struct {
		response.Response
		Data modelsOc.OrderInfo
	}{}
	result.Scan(resultInfo)
	return resultInfo.Data
}

func CheckOrderIsAfterPendingFromOc(tx *gorm.DB, orderId string) bool {
	result := ocClient.ApiByDbContext(tx).CheckOrderIsAfterPending(orderId)
	resultInfo := &struct {
		response.Response
		Data bool
	}{}
	result.Scan(resultInfo)
	return resultInfo.Data
}

func GetUserInfoByIdFromUc(tx *gorm.DB, userId int) dtoUc.UserInfoGetRes {
	result := ucClient.ApiByDbContext(tx).GetUserInfoById(userId)
	resultInfo := &struct {
		response.Response
		Data dtoUc.UserInfoGetRes
	}{}
	result.Scan(resultInfo)
	return resultInfo.Data
}
