package models

import (
	"errors"
	//"go-admin/app/wc/service/admin/dto"
	"go-admin/common/global"
	"go-admin/common/models"
	"go-admin/common/utils"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	EntryProductIsDefective0 = "0"
	EntryProductIsDefective1 = "1"
)

var (
	EntryProductIsDefectiveMap = map[string]string{
		EntryProductIsDefective0: "否",
		EntryProductIsDefective1: "是",
	}
)

type StockEntryDefectiveProduct struct {
	StockEntryProducts           // 入库商品明细
	DefectiveStockLocationId int // 用户选择的次品库位ID
	PassedStockLocationId    int // 系统选择的正品库位ID
}

type StockEntryProducts struct {
	models.Model

	EntryCode             string                  `json:"entryCode" gorm:"type:varchar(20);comment:入库单code"`
	SkuCode               string                  `json:"skuCode" gorm:"type:varchar(20);comment:sku"`
	Quantity              int                     `json:"quantity" gorm:"type:int unsigned;comment:应入库数量"`
	ActQuantity           int                     `json:"actQuantity" gorm:"type:int unsigned;comment:实际入库数量"`
	VendorId              int                     `json:"vendorId" gorm:"type:int unsigned;comment:货主id"`
	WarehouseCode         string                  `json:"warehouseCode" gorm:"type:varchar(20);comment:实体仓code"`
	LogicWarehouseCode    string                  `json:"logicWarehouseCode" gorm:"type:varchar(20);comment:逻辑仓code"`
	GoodsId               int                     `json:"goodsId" gorm:"type:int unsigned;comment:goods表id"`
	StockLocationId       int                     `json:"stockLocationId" gorm:"type:int unsigned;comment:库位id"`
	IsDefective           int                     `json:"isDefective" gorm:"type:tinyint(1);comment:是否为次品退货入库 0 否 1 是"`
	Status                int                     `json:"status" gorm:"type:tinyint(1);comment:入库状态:0-未入库,1-部分入库,2-入库完成"`
	EntryTime             time.Time               `json:"entryTime" gorm:"type:datetime;comment:最早入库时间"`
	EntryEndTime          time.Time               `json:"entryEndTime" gorm:"type:datetime;comment:入库完成时间"`
	StashLocationId       int                     `json:"stashLocationId" gorm:"type:int unsigned;comment:暂存库位id"`
	StashActQuantity      int                     `json:"stashActQuantity" gorm:"type:int unsigned;comment:暂存数量"`
	StockLocation         StockLocation           `json:"-" gorm:"foreignKey:StockLocationId;references:Id"`
	StockEntryProductsSub []StockEntryProductsSub `json:"stockEntryProductsSub" gorm:"foreignKey:EntryProductId;references:Id"`
	models.ModelTime
	models.ControlBy
}

type StockEntryProductsSubCustom struct {
	StockEntryProducts
	EntryCode           string        `json:"entryCode"`
	SkuCode             string        `json:"skuCode"`
	Status              int           `json:"status"`
	EntryTime           time.Time     `json:"entryTime"`
	EntryEndTime        time.Time     `json:"entryEndTime"`
	Quantity            int           `json:"quantity"`
	ActQuantity         int           `json:"actQuantity"`
	WarehouseCode       string        `json:"warehouseCode"`
	LogicWarehouseCode  string        `json:"logicWarehouseCode"`
	GoodsId             int           `json:"goodsId"`
	VendorId            int           `json:"vendorId"`
	HasSub              int           `json:"hasSub"`
	AutoStockLocationId int           `json:"autoStockLocationId"`
	AutoStockLocation   StockLocation `json:"-" gorm:"foreignKey:AutoStockLocationId;references:Id"`
	//子表信息
	SubId               int                            `json:"subId"`
	SubLog              []*StockOutboundProductsSubLog `json:"subLog" gorm:"-"`
	LocationQuantity    int                            `json:"locationQuantity"`
	LocationActQuantity int                            `json:"locationActQuantity"`
	StockLocationId     int                            `json:"stockLocationId"`
	StockLocation       StockLocation                  `json:"-" gorm:"foreignKey:StockLocationId;references:Id"`
}

func (StockEntryProducts) TableName() string {
	return "stock_entry_products"
}

func (e *StockEntryProducts) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *StockEntryProducts) GetId() interface{} {
	return e.Id
}

// 入库来源商品数量
type SourceProductsInfo struct {
	SkuCode     string `json:"skuCode"`     // SKU
	ActQuantity int    `json:"actQuantity"` // 实际出库数量
	Status      string `json:"status"`      // 来源单状态 | 出库商品 0-未出库 1-部分出库 2-出库完成 |
}

// 更新入库商品表
func (e *StockEntryProducts) Save(tx *gorm.DB) error {
	return tx.Omit(clause.Associations).Save(e).Error
}

// 是否为次品
func (e *StockEntryProducts) CheckIsDefective() bool {
	return e.IsDefective == 1
}

// 次品退货入库商品处理
func (e *StockEntryProducts) DefectiveHandle(tx *gorm.DB) (*StockEntryDefectiveProduct, error) {
	outData := &StockEntryDefectiveProduct{}
	outData.DefectiveStockLocationId = e.StockLocationId
	stockLocations, _, err := GetStockLocationsForEntry(tx, e.LogicWarehouseCode, e.GoodsId)
	if err != nil {
		return nil, err
	}
	if len(*stockLocations) == 0 {
		return nil, errors.New(e.LogicWarehouseCode + "逻辑仓下缺少有效库位")
	}
	// 更改真实入库库位
	realLocation := (*stockLocations)[0]
	e.StockLocationId = realLocation.Id
	outData.PassedStockLocationId = realLocation.Id

	if err := utils.CopyDeep(outData, e); err != nil {
		return nil, err
	}
	return outData, nil
}

func (e *StockEntryProducts) GetList(tx *gorm.DB, entryCode string) ([]StockEntryProducts, error) {
	var productList = &[]StockEntryProducts{}
	err := tx.Model(e).
		Where("entry_code = ?", entryCode).Find(productList).Error

	return *productList, err
}

func (e *StockEntryProducts) GetActQuantityMapBySourceCode(tx *gorm.DB, sourceCode string) map[string]int {
	entry := &StockEntry{}
	_ = entry.GetByEntrySourceCode(tx, sourceCode)
	productList, _ := e.GetList(tx, entry.EntryCode)
	return lo.Associate(productList, func(f StockEntryProducts) (string, int) {
		return f.SkuCode, f.ActQuantity
	})
}

func (e *StockEntryProductsSubCustom) GetList2(tx *gorm.DB, entryCode []string) (*[]StockEntryProductsSubCustom, error) {
	var productList = &[]StockEntryProductsSubCustom{}
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	err := tx.Table(wcPrefix+"."+e.TableName()).Preload("StockLocation", func(db *gorm.DB) *gorm.DB {
		return db.Table(wcPrefix + ".stock_location")
	}).Joins("LEFT JOIN "+wcPrefix+".stock_entry_products_sub productSub ON productSub.entry_product_id = stock_entry_products.id").
		Select("stock_entry_products.id,stock_entry_products.logic_warehouse_code,stock_entry_products.warehouse_code,stock_entry_products.sku_code,"+
			"stock_entry_products.goods_id,stock_entry_products.vendor_id,productSub.should_quantity as LocationQuantity,productSub.act_quantity as LocationActQuantity,stock_entry_products.quantity,stock_entry_products.act_quantity,stock_entry_products.entry_code,"+
			"productSub.entry_time,stock_entry_products.entry_end_time,stock_entry_products.status,"+
			"productSub.stock_location_id stock_location_id,productSub.id sub_id,stock_entry_products.is_defective").
		Where("stock_entry_products.act_quantity > 0").
		Where("stock_entry_products.entry_code IN ?", entryCode).
		Where("stock_entry_products.deleted_at is null").
		Debug().
		Find(productList).Error

	return productList, err
}

func (e *StockEntryProductsSubCustom) GetList23(tx *gorm.DB, entryCode []string) ([]StockEntryProductsSubCustom, error) {
	var productList = &[]StockEntryProductsSubCustom{}
	//err := tx.Model(e).
	//	Where("entry_code in ?", entryCode).Find(productList).Error
	//
	//return *productList, err

	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	err := tx.Table(wcPrefix+"."+e.TableName()).Preload("StockLocation", func(db *gorm.DB) *gorm.DB {
		return db.Table(wcPrefix + ".stock_location")
	}).Preload("AutoStockLocation", func(db *gorm.DB) *gorm.DB {
		return db.Table(wcPrefix + ".stock_location")
	}).Joins("LEFT JOIN "+wcPrefix+".stock_entry_products_sub productSub ON productSub.entry_product_id = stock_entry_products.id").
		Select("stock_entry_products.id,stock_entry_products.logic_warehouse_code,stock_entry_products.warehouse_code,stock_entry_products.sku_code,"+
			"stock_entry_products.goods_id,stock_entry_products.vendor_id,stock_entry_products.quantity,stock_entry_products.act_quantity,stock_entry_products.has_sub,stock_entry_products.entry_code,"+
			"stock_entry_products.entry_time,stock_entry_products.entry_end_time,stock_entry_products.status,"+
			"stock_entry_products.stock_location_id auto_stock_location_id,productSub.location_act_quantity,productSub.location_quantity,productSub.stock_location_id,productSub.id sub_id").
		Where("stock_entry_products.entry_code IN ?", entryCode).
		Where("stock_entry_products.deleted_at is null").
		Debug().
		Find(productList).Error

	if err != nil {
		return nil, err
	}

	// 补充数据
	for index, item := range *productList {
		// 次品自动生成的 调拨单->出库单 的商品库位id处理
		if !item.IsHasSub() {
			if item.AutoStockLocationId != 0 {
				(*productList)[index].StockLocationId = item.AutoStockLocationId
				(*productList)[index].StockLocation = item.AutoStockLocation
			}
			(*productList)[index].LocationQuantity = item.Quantity
			(*productList)[index].LocationActQuantity = item.ActQuantity
		}
	}

	return *productList, err
}

func (e *StockEntryProductsSubCustom) IsHasSub() bool {
	return e.HasSub != 0
}
