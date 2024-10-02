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

type StockLog struct {
	service.Service
}

// 库存日志数据权限

func StockLogPermission(p *actions.DataPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Where("stock_log.warehouse_code in ?", utils.Split(p.AuthorityWarehouseId))
		return db
	}
}

// GetPage 获取StockLog列表
func (e *StockLog) GetPage(c *dto.StockLogGetPageReq, p *actions.DataPermission, outData *[]dto.StockLogGetPageResp, count *int64) error {
	err := e.GetDataList(c, p, outData, count)
	if err != nil {
		e.Log.Errorf("StockInfoService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *StockLog) GetDataList(c *dto.StockLogGetPageReq, p *actions.DataPermission, outData *[]dto.StockLogGetPageResp, count *int64, conditions ...func(*gorm.DB) *gorm.DB) error {
	var err error
	var data models.StockLog
	//var list *[]models.StockLog
	pcPrefix := global.GetTenantPcDBNameWithDB(e.Orm)
	db := e.Orm.Model(&data).Preload("LogicWarehouse").
		Joins("LEFT JOIN warehouse ON warehouse.warehouse_code = stock_log.warehouse_code AND warehouse.deleted_at IS NULL").
		Joins("LEFT JOIN vendors ON vendors.id = stock_log.vendor_id AND vendors.deleted_at IS NULL").
		Joins("LEFT JOIN " + pcPrefix + ".product product ON product.sku_code = stock_log.sku_code AND product.deleted_at IS NULL").
		Joins("LEFT JOIN " + pcPrefix + ".brand brand ON brand.id = product.brand_id AND brand.deleted_at IS NULL").
		Joins("LEFT JOIN stock_info ON stock_info.id = stock_log.stock_info_id AND stock_info.deleted_at IS NULL").
		Joins("LEFT JOIN " + pcPrefix + ".goods goods ON goods.id = stock_info.goods_id AND goods.deleted_at IS NULL")
	if len(conditions) > 0 {
		db.Scopes(conditions...)
	} else {
		db.Scopes(
			dtoStock.GenProductSearch(c.Sku, c.ProductNo, c.ProductName, "stock_log"),
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			StockLogPermission(p),
			dtoStock.GenCreatedAtTimeSearch(c.CreatedAtStart, c.CreatedAtEnd, "stock_log"),
			dtoStock.GenWarehousesSearch(c.QueryWarehouseCode, "stock_log"),
		)
	}
	err = db.
		Select("stock_log.*,warehouse.warehouse_name,vendors.name_zh vendor_name,vendors.code vendor_code,vendors.short_name vendor_short_name," +
			"product.supplier_sku_code vendor_sku_code,product.name_zh product_name,product.mfg_model,product.sales_uom,goods.product_no,brand.brand_zh brand_name").
		Order("stock_log.id DESC").
		Find(outData).Limit(-1).Offset(-1).
		Count(count).Error
	e.FormatOutput(outData)
	if err != nil {
		e.Log.Errorf("StockLogService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *StockLog) FormatOutput(outData *[]dto.StockLogGetPageResp) {
	for index, item := range *outData {
		(*outData)[index].LogicWarehouseName = item.LogicWarehouse.LogicWarehouseName
		(*outData)[index].FromTypeName = utils.GetTFromMap(item.FromType, models.StockLogFromTypeMap)
	}
}

// Export 获取StockInfo列表
func (e *StockLog) Export(c *dto.StockLogGetPageReq, p *actions.DataPermission, outData *[]dto.StockLogGetPageResp) error {
	var err error
	var count int64
	c.PageIndex = 1
	c.PageSize = -1
	if c.Ids != "" {
		err = e.GetDataList(c, p, outData, &count, func(db *gorm.DB) *gorm.DB {
			db.Where("stock_log.id IN ?", utils.Split(c.Ids))
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
