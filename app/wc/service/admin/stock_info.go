package admin

import (
	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	dtoStock "go-admin/common/dto/stock/dto"
	"go-admin/common/global"
	"go-admin/common/utils"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"
)

type StockInfo struct {
	service.Service
}

// 库存日志数据权限
func StockInfoPermission(p *actions.DataPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Where("stock_info.warehouse_code in ?", utils.Split(p.AuthorityWarehouseId))
		db.Where("stock_info.vendor_id in ?", utils.SplitToInt(p.AuthorityVendorId))
		return db
	}
}

// GetPage 获取StockInfo列表
func (e *StockInfo) GetPage(c *dto.StockInfoGetPageReq, p *actions.DataPermission, outData *[]dto.StockInfoGetPageResp, count *int64) error {
	err := e.GetDataList(c, p, outData, count)
	if err != nil {
		e.Log.Errorf("StockInfoService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *StockInfo) GetDataList(c *dto.StockInfoGetPageReq, p *actions.DataPermission, outData *[]dto.StockInfoGetPageResp, count *int64, conditions ...func(*gorm.DB) *gorm.DB) error {
	var err error
	var data models.StockInfo
	//var list *[]models.StockInfo
	pcPrefix := global.GetTenantPcDBNameWithDB(e.Orm)
	db := e.Orm.Model(&data).Preload("LogicWarehouse").
		Joins("LEFT JOIN warehouse ON warehouse.warehouse_code = stock_info.warehouse_code AND warehouse.deleted_at IS NULL").
		Joins("LEFT JOIN logic_warehouse lw on stock_info.logic_warehouse_code = lw.logic_warehouse_code and lw.deleted_at IS NULL").
		Joins("LEFT JOIN vendors ON vendors.id = stock_info.vendor_id AND vendors.deleted_at IS NULL").
		Joins("LEFT JOIN "+pcPrefix+".product product ON product.sku_code = stock_info.sku_code AND product.deleted_at IS NULL").
		Joins("LEFT JOIN "+pcPrefix+".brand brand ON brand.id = product.brand_id AND brand.deleted_at IS NULL").
		Joins("LEFT JOIN "+pcPrefix+".goods goods ON goods.id = stock_info.goods_id AND goods.deleted_at IS NULL").
		Scopes(
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			StockInfoPermission(p),
			actions.Permission(data.TableName(), p),
		)

	if len(conditions) > 0 {
		db.Scopes(conditions...)
	} else {
		db.Scopes(
			dtoStock.GenProductSearch(c.Sku, c.ProductNo, c.ProductName, "stock_info"),
			cDto.MakeCondition(c.GetNeedSearch()),
			dtoStock.GenWarehousesSearch(c.QueryWarehouseCode, "stock_info"),
		)
	}
	err = db.Select("stock_info.*,warehouse.warehouse_name,lw.logic_warehouse_name,vendors.name_zh vendor_name,vendors.code vendor_code,vendors.short_name vendor_short_name," +
		"product.supplier_sku_code vendor_sku_code,product.name_zh product_name,product.mfg_model,product.sales_uom,goods.product_no,brand.brand_zh brand_name").
		//"SUM(order_detail.quantity) order_quantity, SUM(order_detail.lock_stock) order_lock_stock").
		Order("stock_info.id DESC").
		Group("stock_info.id").
		Find(outData).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		return err
	}

	orderQuantity := make([]dto.OrderQuantity, 0)
	err = e.GetOrderQuantity(c, p, &orderQuantity)
	if err != nil {
		return err
	}

	locationStock := make(map[string]map[int][]dto.LocationStock, 0)
	stockInfoIds := make([]int, 0)
	for _, item := range *outData {
		stockInfoIds = append(stockInfoIds, item.Id)
	}
	err = e.GetLocationStock(c, p, stockInfoIds, &locationStock)
	if err != nil {
		return err
	}

	e.FormatOutput(outData, &orderQuantity, &locationStock)

	return nil
}

func (e *StockInfo) GetOrderQuantity(c *dto.StockInfoGetPageReq, p *actions.DataPermission, orderQuantity *[]dto.OrderQuantity) error {
	var data models.StockInfo
	ocPrefix := global.GetTenantOcDBNameWithDB(e.Orm)
	err := e.Orm.Model(&data).
		Joins("LEFT JOIN warehouse ON warehouse.warehouse_code = stock_info.warehouse_code AND warehouse.deleted_at IS NULL").
		Joins("LEFT JOIN "+ocPrefix+".order_detail order_detail ON order_detail.warehouse_code = stock_info.warehouse_code AND order_detail.sku_code = stock_info.sku_code").
		Joins("RIGHT JOIN "+ocPrefix+".order_info order_info ON order_info.order_id = order_detail.order_id AND order_info.order_status = 6").
		Select("stock_info.id, SUM(order_detail.quantity) order_quantity, SUM(order_detail.cancel_quantity) order_cancel_quantity, SUM(order_detail.lock_stock) order_lock_stock").
		Scopes(
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			StockInfoPermission(p),
			actions.Permission(data.TableName(), p),
		).
		Order("stock_info.id DESC").
		Group("stock_info.id").
		Find(&orderQuantity).Limit(-1).Offset(-1).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (e *StockInfo) GetLocationStock(c *dto.StockInfoGetPageReq, p *actions.DataPermission, stockInfoIds []int, locationStock *map[string]map[int][]dto.LocationStock) error {
	var data models.StockInfo
	locationStockList := make([]dto.LocationStock, 0)
	err := e.Orm.Model(&data).Debug().
		Joins("LEFT JOIN stock_location ON stock_location.warehouse_code = stock_info.warehouse_code AND stock_location.logic_warehouse_code = stock_info.logic_warehouse_code").
		Joins("LEFT JOIN stock_location_goods ON stock_location_goods.goods_id = stock_info.goods_id AND stock_location.id = stock_location_goods.location_id").
		Where("stock_info.id in ?", stockInfoIds).
		Where("stock_location.status = 1").
		Where("stock_location_goods.stock != 0 OR stock_location_goods.lock_stock != 0").
		Select("stock_info.id stock_info_id, stock_info.warehouse_code, stock_info.logic_warehouse_code, stock_location.location_code, stock_location_goods.stock, stock_location_goods.lock_stock").
		Scopes(
			StockInfoPermission(p),
			actions.Permission(data.TableName(), p),
		).
		Order("stock_info.id DESC").
		Limit(-1).Offset(-1).
		Find(&locationStockList).
		Error
	if err != nil {
		return err
	}

	locationStockInfo := make(map[int][]dto.LocationStock, 0)
	locationLockStockInfo := make(map[int][]dto.LocationStock, 0)
	locationTotalStockInfo := make(map[int][]dto.LocationStock, 0)

	for _, item := range locationStockList {
		totalStock := item.Stock + item.LockStock
		if item.Stock > 0 {
			locationStockInfo[item.StockInfoId] = append(locationStockInfo[item.StockInfoId], dto.LocationStock{
				StockInfoId:        item.StockInfoId,
				WarehouseCode:      item.WarehouseCode,
				LogicWarehouseCode: item.LogicWarehouseCode,
				LocationCode:       item.LocationCode,
				Stock:              item.Stock,
				LockStock:          item.LockStock,
				TotalStock:         totalStock,
			})
		}
		if item.LockStock > 0 {
			locationLockStockInfo[item.StockInfoId] = append(locationLockStockInfo[item.StockInfoId], dto.LocationStock{
				StockInfoId:        item.StockInfoId,
				WarehouseCode:      item.WarehouseCode,
				LogicWarehouseCode: item.LogicWarehouseCode,
				LocationCode:       item.LocationCode,
				Stock:              item.Stock,
				LockStock:          item.LockStock,
				TotalStock:         totalStock,
			})
		}
		if totalStock > 0 {
			locationTotalStockInfo[item.StockInfoId] = append(locationTotalStockInfo[item.StockInfoId], dto.LocationStock{
				StockInfoId:        item.StockInfoId,
				WarehouseCode:      item.WarehouseCode,
				LogicWarehouseCode: item.LogicWarehouseCode,
				LocationCode:       item.LocationCode,
				Stock:              item.Stock,
				LockStock:          item.LockStock,
				TotalStock:         totalStock,
			})
		}
	}

	(*locationStock)["stock"] = locationStockInfo
	(*locationStock)["lockStock"] = locationLockStockInfo
	(*locationStock)["totalStock"] = locationTotalStockInfo

	return nil
}

func (e *StockInfo) FormatOutput(outData *[]dto.StockInfoGetPageResp, OrderQuantity *[]dto.OrderQuantity, locationStock *map[string]map[int][]dto.LocationStock) {
	for index, item := range *outData {
		// (*outData)[index].LogicWarehouseName = item.LogicWarehouse.LogicWarehouseName
		(*outData)[index].TotalStock = item.Stock + item.LockStock
		(*outData)[index].LocationStocks = (*locationStock)["stock"][item.Id]
		(*outData)[index].LocationLockStocks = (*locationStock)["lockStock"][item.Id]
		(*outData)[index].LocationTotalStocks = (*locationStock)["totalStock"][item.Id]
		for _, v := range *OrderQuantity {
			if item.Id == v.Id {
				// 实缺数量 = 订单购买数量 - 订单取消数量 - 订单占位库存
				lackStock := v.OrderQuantity - v.OrderCancelQuantity - v.OrderLockStock
				if lackStock > 0 {
					(*outData)[index].LackStock = lackStock
				}
			}
		}
	}
}

// GetPage 获取StockInfo列表
func (e *StockInfo) Export(c *dto.StockInfoGetPageReq, p *actions.DataPermission, outData *[]dto.StockInfoGetPageResp) error {
	var err error
	var count int64
	c.PageIndex = 1
	c.PageSize = -1
	if c.Ids != "" {
		err = e.GetDataList(c, p, outData, &count, func(db *gorm.DB) *gorm.DB {
			db.Where("stock_info.id IN ?", utils.Split(c.Ids))
			return db
		})
	} else {
		err = e.GetDataList(c, p, outData, &count)
	}
	if err != nil {
		e.Log.Errorf("StockInfoService Export error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *StockInfo) InnerGetByGoodsIdAndLwhCode(c *dto.InnerStockInfoGetByGoodsIdAndLwhCodeReq, outData *[]models.StockInfo) error {
	var err error
	var data models.StockInfo
	db := e.Orm.Model(&data)
	for _, info := range c.Query {
		db.Or(&info)
	}
	err = db.Find(outData).Error
	if err != nil {
		e.Log.Errorf("WarehouseService InnerGetByGoodsIdAndLwhCode error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *StockInfo) InnerGetByGoodsIdAndWarehouseCode(c *dto.InnerStockInfoGetByGoodsIdAndWarehouseCodeReq, outData *[]models.StockInfo) error {
	var data models.StockInfo
	// 获取逻辑仓
	lwh := &models.LogicWarehouse{}
	if err := lwh.GetPassLogicWarehouseByWhCode(e.Orm, c.WarehouseCode); err != nil {
		return err
	}
	err := e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			func(db *gorm.DB) *gorm.DB {
				db.Where("logic_warehouse_code = ?", lwh.LogicWarehouseCode)
				return db
			},
		).Find(outData).Error
	if err != nil {
		e.Log.Errorf("WarehouseService InnerGetByGoodsIdAndWarehouseCode error:%s \r\n", err)
		return err
	}
	return nil
}
