package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	"go-admin/common/global"
	"gorm.io/gorm"
)

type StockLocationGoods struct {
	service.Service
}

// GetPage 获取StockLocationGoods列表
func (e *StockLocationGoods) GetPage(c *dto.StockLocationGoodsGetPageReq, p *actions.DataPermission, list *[]dto.StockLocationGoodsListResp, count *int64) error {
	var err error
	var data models.StockLocationGoods

	pcPrefix := global.GetTenantPcDBNameWithDB(e.Orm)
	err = e.Orm.Model(&data).
		Select("stock_location_goods.*, (stock_location_goods.stock + stock_location_goods.lock_stock) total_stock, sl.location_code, w.warehouse_name, g.sku_code, g.supplier_sku_code, "+
			"g.product_no, v.name_zh vendor_name, p.name_zh product_name, p.mfg_model, b.brand_zh brand_name").
		Joins("left join stock_location sl on stock_location_goods.location_id = sl.id").
		Joins("left join warehouse w on sl.warehouse_code = w.warehouse_code").
		Joins("left join "+pcPrefix+".goods g on stock_location_goods.goods_id = g.id").
		Joins("left join "+pcPrefix+".product p on g.sku_code = p.sku_code").
		Joins("left join "+pcPrefix+".brand b on p.brand_id = b.id").
		Joins("left join vendors v on g.vendor_id = v.id").
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				if len(c.Ids) > 0 {
					db = db.Where("stock_location_goods.id IN ?", c.Ids)
				}

				db = db.Where("w.is_virtual = 0")
				db = db.Where("(stock_location_goods.stock > 0 or stock_location_goods.lock_stock > 0)")
				return db
			},
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			actions.SysUserPermission("w", p, 2),
		).
		Scan(&list).Limit(-1).Offset(-1).
		Count(count).Error

	if err != nil {
		e.Log.Errorf("StockLocationGoodsService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取StockLocationGoods对象
func (e *StockLocationGoods) Get(d *dto.StockLocationGoodsGetReq, p *actions.DataPermission, model *dto.StockLocationGoodsListResp) error {
	var data models.StockLocationGoods

	pcPrefix := global.GetTenantPcDBNameWithDB(e.Orm)
	err := e.Orm.Model(&data).
		Select("stock_location_goods.*, (stock_location_goods.stock + stock_location_goods.lock_stock) total_stock, sl.location_code, w.warehouse_name, g.sku_code, g.supplier_sku_code, "+
			"g.product_no, v.name_zh vendor_name, p.name_zh product_name, p.mfg_model, b.brand_zh brand_name").
		Joins("left join stock_location sl on stock_location_goods.location_id = sl.id").
		Joins("left join warehouse w on sl.warehouse_code = w.warehouse_code").
		Joins("left join "+pcPrefix+".goods g on stock_location_goods.goods_id = g.id").
		Joins("left join "+pcPrefix+".product p on g.sku_code = p.sku_code").
		Joins("left join "+pcPrefix+".brand b on p.brand_id = b.id").
		Joins("left join vendors v on g.vendor_id = v.id").
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetStockLocationGoods error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// GetSameLogicStockLocationList 获取同一个逻辑仓下的仓位列表
func (e *StockLocationGoods) GetSameLogicStockLocationList(c *dto.SameLogicStockLocationReq, p *actions.DataPermission, list *[]cDto.Option, count *int64) error {
	var err error
	var data models.StockLocation

	err = e.Orm.Model(&data).
		Select("stock_location.id as value, stock_location.location_code as label").
		Joins("left join stock_location sl on stock_location.logic_warehouse_code = sl.logic_warehouse_code").
		Joins("left join logic_warehouse lw on stock_location.logic_warehouse_code = lw.warehouse_code").
		Where("stock_location.status = 1").
		Where("stock_location.location_code <> ?", c.LocationCode).
		Where("sl.location_code = ?", c.LocationCode).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Scan(&list).Error

	if err != nil {
		e.Log.Errorf("获取同一个逻辑仓下的仓位列表 error:%s \r\n", err)
		return err
	}
	return nil
}

// TransferStock 转移数量
func (e *StockLocationGoods) TransferStock(c *dto.TransferStockReq, p *actions.DataPermission, gc *gin.Context) error {
	var err error
	var oriData models.StockLocationGoods

	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	err = tx.Where("id = ?", c.Id).First(&oriData).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return errors.New("查询不到原库存商品")
	}
	if oriData.Stock < c.TransferStock {
		tx.Rollback()
		return errors.New("原库位数量不足，请刷新页面重写尝试")
	}
	// 扣除原库位库存
	beforeStock := oriData.Stock
	oriData.Stock = oriData.Stock - c.TransferStock
	oriData.SetUpdateBy(user.GetUserId(gc))
	err = tx.Save(&oriData).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	log := models.StockLocationGoodsLog{
		FromType:             "3",
		StockLocationGoodsId: oriData.Id,
		BeforeStock:          beforeStock + oriData.LockStock,
		AfterStock:           oriData.Stock + oriData.LockStock,
		ChangeStock:          -c.TransferStock,
		CreateBy:             user.GetUserId(gc),
		CreateByName:         user.GetUserName(gc),
	}
	if _, err := log.GenerateStockLocationGoodsLogCode(tx); err != nil {
		return err
	}
	err = tx.Create(&log).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 添加新库存
	var targetData models.StockLocationGoods
	err = tx.Where("goods_id = ?", c.GoodsId).Where("location_id = ?", c.TargetLocationId).First(&targetData).Error
	//var log2 models.StockLocationGoodsLog

	targetBeforeStock := 0
	if errors.Is(err, gorm.ErrRecordNotFound) {

		targetData.GoodsId = c.GoodsId
		targetData.LocationId = c.TargetLocationId
		targetData.Stock = c.TransferStock
		targetData.SetCreateBy(user.GetUserId(gc))
	} else {
		targetBeforeStock = targetData.Stock + targetData.LockStock
		targetData.Stock = targetData.Stock + c.TransferStock
		targetData.SetUpdateBy(user.GetUserId(gc))
	}
	err = tx.Save(&targetData).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	log2 := models.StockLocationGoodsLog{
		DocketCode:           log.DocketCode,
		FromType:             "3",
		StockLocationGoodsId: targetData.Id,
		BeforeStock:          targetBeforeStock,
		AfterStock:           targetBeforeStock + c.TransferStock,
		ChangeStock:          c.TransferStock,
		CreateBy:             user.GetUserId(gc),
		CreateByName:         user.GetUserName(gc),
	}
	err = e.Orm.Create(&log2).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
