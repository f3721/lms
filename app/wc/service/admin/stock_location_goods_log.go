package admin

import (
	"github.com/go-admin-team/go-admin-core/sdk/service"
	dtoStock "go-admin/common/dto/stock/dto"
	"go-admin/common/global"
	"go-admin/common/utils"
	"gorm.io/gorm"

	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type StockLocationGoodsLog struct {
	service.Service
}

// StockLocationLogPermission 库存日志数据权限
func StockLocationLogPermission(p *actions.DataPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Where("stock_location.warehouse_code in ?", utils.Split(p.AuthorityWarehouseId))
		return db
	}
}

// GetPage 获取StockLocationGoodsLog列表
func (e *StockLocationGoodsLog) GetPage(c *dto.StockLocationGoodsLogGetPageReq, p *actions.DataPermission, outData *[]dto.StockLocationGoodsLogGetPageResp, count *int64) error {
	err := e.GetDataList(c, p, outData, count)
	if err != nil {
		e.Log.Errorf("StockInfoService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *StockLocationGoodsLog) GetDataList(c *dto.StockLocationGoodsLogGetPageReq, p *actions.DataPermission, outData *[]dto.StockLocationGoodsLogGetPageResp, count *int64, conditions ...func(*gorm.DB) *gorm.DB) error {
	var err error
	var data models.StockLocationGoodsLog
	//var list *[]models.StockLog
	pcPrefix := global.GetTenantPcDBNameWithDB(e.Orm)
	db := e.Orm.Model(&data).
		Joins("LEFT JOIN stock_location_goods ON stock_location_goods.id = stock_location_goods_log.stock_location_goods_id AND stock_location_goods.deleted_at IS NULL").
		Joins("LEFT JOIN stock_location ON stock_location.id = stock_location_goods.location_id AND stock_location.deleted_at IS NULL").
		Joins("LEFT JOIN warehouse ON warehouse.warehouse_code = stock_location.warehouse_code AND warehouse.deleted_at IS NULL").
		Joins("LEFT JOIN logic_warehouse ON logic_warehouse.logic_warehouse_code = stock_location.logic_warehouse_code AND logic_warehouse.warehouse_code = stock_location.warehouse_code AND logic_warehouse.deleted_at IS NULL").
		Joins("LEFT JOIN " + pcPrefix + ".goods goods ON goods.id = stock_location_goods.goods_id AND goods.deleted_at IS NULL").
		Joins("LEFT JOIN " + pcPrefix + ".product product ON product.sku_code = goods.sku_code AND stock_location.warehouse_code = goods.warehouse_code AND product.deleted_at IS NULL").
		Joins("LEFT JOIN " + pcPrefix + ".brand brand ON brand.id = product.brand_id AND brand.deleted_at IS NULL").
		Joins("LEFT JOIN vendors ON vendors.id = product.vendor_id AND vendors.deleted_at IS NULL")
	if len(conditions) > 0 {
		db.Scopes(conditions...)
	} else {
		db.Scopes(
			dtoStock.GenProductSearch(c.Sku, c.ProductNo, c.ProductName, "goods"),
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			StockLocationLogPermission(p),
			dtoStock.GenCreatedAtTimeSearch(c.CreatedAtStart, c.CreatedAtEnd, "stock_location_goods_log"),
			dtoStock.GenWarehousesSearch(c.QueryWarehouseCode, "stock_location"),
		)
	}
	err = db.
		Select("stock_location_goods_log.*, stock_location.location_code, stock_location.warehouse_code, stock_location.logic_warehouse_code," +
			"warehouse.warehouse_name,logic_warehouse.logic_warehouse_name,vendors.name_zh vendor_name,vendors.code vendor_code,vendors.short_name vendor_short_name," +
			"product.supplier_sku_code vendor_sku_code,product.name_zh product_name,product.mfg_model,product.sales_uom,goods.product_no," +
			"brand.brand_zh brand_name,product.sku_code,stock_location.id location_id, stock_location.location_code").
		Order("stock_location_goods_log.id DESC").
		Find(outData).Limit(-1).Offset(-1).
		Count(count).Error
	e.FormatOutput(outData)
	if err != nil {
		e.Log.Errorf("StockLogService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *StockLocationGoodsLog) FormatOutput(outData *[]dto.StockLocationGoodsLogGetPageResp) {
	for index, item := range *outData {
		(*outData)[index].FromTypeName = utils.GetTFromMap(item.FromType, models.StockLocationGoodsLogFromTypeMap)
	}
}

// Export 获取StockInfo列表
func (e *StockLocationGoodsLog) Export(c *dto.StockLocationGoodsLogGetPageReq, p *actions.DataPermission, outData *[]dto.StockLocationGoodsLogGetPageResp) error {
	var err error
	var count int64
	c.PageIndex = 1
	c.PageSize = -1
	if c.Ids != "" {
		err = e.GetDataList(c, p, outData, &count, func(db *gorm.DB) *gorm.DB {
			db.Where("stock_location_goods_log.id IN ?", utils.Split(c.Ids))
			return db
		})
	} else {
		err = e.GetDataList(c, p, outData, &count)
	}

	for k := range *outData {
		(*outData)[k].CreatedTime = (*outData)[k].CreatedAt.Format("2006-01-02 15:04:05")
	}
	if err != nil {
		e.Log.Errorf("StockInfoService Export error:%s \r\n", err)
		return err
	}
	return nil
}
