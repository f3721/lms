package models

import (
	"errors"
	"fmt"
	"go-admin/common/global"
	"go-admin/common/models"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type StockOutboundProducts struct {
	models.Model

	OutboundCode       string        `json:"outboundCode" gorm:"type:varchar(20);comment:出库单code"`
	Status             int           `json:"status" gorm:"type:tinyint(1);comment:0-未出库,1-部分出库,2-出库完成"`
	OutboundTime       time.Time     `json:"outboundTime" gorm:"type:datetime;comment:首次出库时间"`
	OutboundEndTime    time.Time     `json:"outboundEndTime" gorm:"type:datetime;comment:出库完成时间"`
	SkuCode            string        `json:"skuCode" gorm:"type:varchar(10);comment:sku"`
	Quantity           int           `json:"quantity" gorm:"type:int unsigned;comment:应出库数量"`
	ActQuantity        int           `json:"actQuantity" gorm:"type:int unsigned;comment:实际出库数量"`
	VendorId           int           `json:"vendorId" gorm:"type:int unsigned;comment:货主id"`
	WarehouseCode      string        `json:"warehouseCode" gorm:"type:varchar(20);comment:实体仓code"`
	LogicWarehouseCode string        `json:"logicWarehouseCode" gorm:"type:varchar(20);comment:逻辑仓code"`
	GoodsId            int           `json:"goodsId" gorm:"type:int unsigned;comment:goods表id"`
	StockLocationId    int           `json:"stockLocationId" gorm:"type:int unsigned;comment:库位id"`
	HasSub             int           `json:"hasSub" gorm:"type:tinyint(1);comment:是否为手动确认入库 0 否 1 是"`
	OriQuantity        int           `json:"oriQuantity" gorm:"type:int unsigned;comment:原始数量"`
	StockLocation      StockLocation `json:"-" gorm:"foreignKey:StockLocationId;references:Id"`
	models.ModelTime
	models.ControlBy
}

func (StockOutboundProducts) TableName() string {
	return "stock_outbound_products"
}

func (e *StockOutboundProducts) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *StockOutboundProducts) GetId() interface{} {
	return e.Id
}

func (e *StockOutboundProducts) Del(tx *gorm.DB) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	return tx.Table(wcPrefix + "." + e.TableName()).Delete(e).Error
}

func (e *StockOutboundProducts) GetById(tx *gorm.DB, id int) error {
	return tx.Table(e.TableName()).Where("id = ?", id).First(e).Error
}

// 出库单产品表扣减应出库数量 | 售后取消
func (e *StockOutboundProducts) SubQuantity(tx *gorm.DB, quantity int) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)

	// 兜底校验
	leftUnOutQuantity := (e.Quantity - e.ActQuantity)
	if quantity > leftUnOutQuantity {
		return fmt.Errorf("出库单商品ID:%v,扣除数量[%v]不能大于剩余最大应出库数量[%v]", e.Id, quantity, leftUnOutQuantity)
	}

	// 扣减应出库数量
	updateMap := map[string]interface{}{
		"quantity": gorm.Expr("quantity - ?", quantity),
	}

	// 维护产品行出库状态 | 如果剩余应出库数量被扣减完了，出库完成
	if leftUnOutQuantity == quantity {
		updateMap["status"] = "2"
		updateMap["outbound_end_time"] = time.Now()
	}

	// 防并发执行
	res := tx.Table(wcPrefix+"."+e.TableName()).Where("id = ?", e.Id).Where("quantity = ?", e.Quantity).Updates(updateMap)
	err := res.Error
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return errors.New("StockOutboundProducts SubQuantity error")
	}
	return nil
}

// 更新出库商品表实际数量

func (e *StockOutboundProducts) Save(tx *gorm.DB) error {
	return tx.Save(e).Error
}

func (e *StockOutboundProducts) GetList(tx *gorm.DB, outboundCode string) ([]StockOutboundProducts, error) {
	var productList = &[]StockOutboundProducts{}
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	err := tx.Table(wcPrefix+"."+e.TableName()).
		Where("outbound_code = ?", outboundCode).Find(productList).Error

	return *productList, err
}

func (e *StockOutboundProducts) GetActQuantityMapBySourceCode(tx *gorm.DB, sourceCode string) map[string]int {
	outbound := &StockOutbound{}
	outbound.GetBySourceCode(tx, sourceCode)
	productList, _ := e.GetList(tx, outbound.OutboundCode)
	return lo.Associate(productList, func(f StockOutboundProducts) (string, int) {
		return f.SkuCode, f.ActQuantity
	})
}

// 根据库位商品信息切分单个商品
func (e *StockOutboundProducts) SplitProductAndLockLocation(tx *gorm.DB, outboundCode, outboundType string) error {
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	stockLocationGoods := &[]StockLocationGoods{}
	quantity := e.Quantity

	if err := tx.Table(wcPrefix+".stock_location_goods").Joins("JOIN "+wcPrefix+".stock_location location ON location.id = stock_location_goods.location_id").
		Select("stock_location_goods.*").
		Where("stock_location_goods.goods_id = ?", e.GoodsId).
		Where("stock_location_goods.stock > ?", 0).
		Where("location.logic_warehouse_code = ?", e.LogicWarehouseCode).
		Where("location.status = ?", StockLocationStatus1).
		Order("stock_location_goods.stock asc").
		Find(stockLocationGoods).Error; err != nil {
		return err
	}

	stockOutboundProductsSub := &StockOutboundProductsSub{
		OutboundProductId: e.Id,
	}

	// 判断是否存在库位库存=应出库数量，定位该库位进行扣减

	var eqStockLocationGoods []StockLocationGoods = lo.Filter[StockLocationGoods](*stockLocationGoods, func(item StockLocationGoods, index int) bool {
		return item.Stock == quantity
	})
	if len(eqStockLocationGoods) > 0 {
		stockOutboundProductsSub.LocationQuantity = eqStockLocationGoods[0].Stock
		stockOutboundProductsSub.StockLocationId = eqStockLocationGoods[0].LocationId
		if err := stockOutboundProductsSub.SaveAndLockLocationGoods(tx, e.GoodsId, outboundCode, outboundType); err != nil {
			return err
		}
		return nil
	}

	for _, item := range *stockLocationGoods {
		stockOutboundProductsSub := &StockOutboundProductsSub{
			OutboundProductId: e.Id,
		}
		stockOutboundProductsSub.LocationQuantity = item.Stock
		stockOutboundProductsSub.StockLocationId = item.LocationId

		quantity -= item.Stock
		if quantity < 0 {
			stockOutboundProductsSub.LocationQuantity = item.Stock + quantity
		}
		if err := stockOutboundProductsSub.SaveAndLockLocationGoods(tx, e.GoodsId, outboundCode, outboundType); err != nil {
			return err
		}
		if quantity <= 0 {
			break
		}
	}
	if quantity > 0 {
		return errors.New(e.SkuCode + ",在当前逻辑仓(" + e.LogicWarehouseCode + ")库位库存不足")
	}
	return nil
}

// [单元测试] 根据库位商品信息切分单个商品
func (e *StockOutboundProducts) SplitProductAndLockLocationTest(tx *gorm.DB, outboundCode string) ([]map[string]any, error) {
	stockLocationGoods := &[]StockLocationGoods{}
	fmt.Println("此商品应扣减总数量:", e.Quantity)
	fmt.Println("此商品SKU:", e.SkuCode)
	quantity := e.Quantity
	res := []map[string]any{}

	if err := tx.Table("stock_location_goods").Joins("JOIN stock_location location ON location.id = stock_location_goods.location_id").
		Select("stock_location_goods.*").
		Where("stock_location_goods.goods_id = ?", e.GoodsId).
		Where("stock_location_goods.stock > ?", 0).
		Where("location.logic_warehouse_code = ?", e.LogicWarehouseCode).
		Where("location.status = ?", StockLocationStatus1).
		Order("stock_location_goods.stock asc").
		Find(stockLocationGoods).Error; err != nil {
		return res, err
	}

	// 判断是否存在库位库存=应出库数量，定位该库位进行扣减
	var eqStockLocationGoods []StockLocationGoods = lo.Filter(*stockLocationGoods, func(item StockLocationGoods, index int) bool {
		return item.Stock == quantity
	})

	// 场景一: 有库位数量相同，取当前库位
	if len(eqStockLocationGoods) > 0 {
		fmt.Println(eqStockLocationGoods[0])
		checkedLocation := eqStockLocationGoods[0]
		checkedLocationInfo := map[string]any{
			"库位ID":    checkedLocation.LocationId,
			"库存数量":    checkedLocation.Stock,
			"本库位扣减数量": quantity,
		}
		res = append(res, checkedLocationInfo)

		return res, nil
	}

	// 场景二: 没有库位数量一致，从小到大依次扣减
	for _, item := range *stockLocationGoods {
		// 依次扣减
		quantity -= item.Stock
		locationReductNum := 0
		if quantity < 0 {
			locationReductNum = item.Stock + quantity
		}
		fmt.Println("此此扣减数量:", locationReductNum)

		// 记录本库位扣减值
		checkedLocation := item
		checkedLocationInfo := map[string]any{
			"库位ID":    checkedLocation.LocationId,
			"库存数量":    checkedLocation.Stock,
			"本库位扣减数量": locationReductNum,
		}
		res = append(res, checkedLocationInfo)

		// 知道扣减结束
		if quantity <= 0 {
			break
		}
	}
	if quantity > 0 {
		return res, errors.New(e.SkuCode + ",在当前逻辑仓(" + e.LogicWarehouseCode + ")库位库存不足")
	}
	return res, nil
}

// 商品和库位记录的降级融合
type OutboundProductSubCustom struct {
	models.Model
	OutboundCode        string        `json:"outboundCode"`
	SkuCode             string        `json:"skuCode"`
	Status              int           `json:"status"`
	OutboundTime        time.Time     `json:"outboundTime"`
	OutboundEndTime     time.Time     `json:"outboundEndTime"`
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

func (OutboundProductSubCustom) TableName() string {
	return "stock_outbound_products"
}

func (e *OutboundProductSubCustom) IsHasSub() bool {
	return e.HasSub != 0
}

// 获取出库单产品信息（包含·库位信息）
func (e *OutboundProductSubCustom) GetList(tx *gorm.DB, outboundCode []string) ([]OutboundProductSubCustom, error) {
	var productList = &[]OutboundProductSubCustom{}
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	err := tx.Table(wcPrefix+"."+e.TableName()).Preload("StockLocation", func(db *gorm.DB) *gorm.DB {
		return db.Table(wcPrefix + ".stock_location")
	}).Preload("AutoStockLocation", func(db *gorm.DB) *gorm.DB {
		return db.Table(wcPrefix + ".stock_location")
	}).Joins("LEFT JOIN "+wcPrefix+".stock_outbound_products_sub productSub ON productSub.outbound_product_id = stock_outbound_products.id").
		Select("stock_outbound_products.id,stock_outbound_products.logic_warehouse_code,stock_outbound_products.warehouse_code,stock_outbound_products.sku_code,"+
			"stock_outbound_products.goods_id,stock_outbound_products.vendor_id,stock_outbound_products.quantity,stock_outbound_products.act_quantity,stock_outbound_products.has_sub,stock_outbound_products.outbound_code,"+
			"stock_outbound_products.outbound_time,stock_outbound_products.outbound_end_time,stock_outbound_products.status,"+
			"stock_outbound_products.stock_location_id auto_stock_location_id,productSub.location_act_quantity,productSub.location_quantity,productSub.stock_location_id,productSub.id sub_id").
		Where("stock_outbound_products.outbound_code IN ?", outboundCode).
		Where("stock_outbound_products.deleted_at is null").
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

// 更新出库商品子表实际数量
func (e *OutboundProductSubCustom) SetActQuantityForProductSub(tx *gorm.DB) error {
	return tx.Model(&StockOutboundProductsSub{}).Where("id = ?", e.SubId).Update("location_act_quantity", e.LocationActQuantity).Error
}
