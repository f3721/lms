package models

import (
	"errors"
	"go-admin/common/global"
	"go-admin/common/models"

	"gorm.io/gorm"
)

type UserCart struct {
	models.Model

	GoodsId       int    `json:"goodsId" gorm:"type:int unsigned;comment:goods表 主键id"`
	UserId        int    `json:"userId" gorm:"type:int unsigned;comment:用户编号"`
	WarehouseCode string `json:"warehouseCode" gorm:"type:varchar(20);comment:WarehouseCode"`
	SkuCode       string `json:"skuCode" gorm:"type:varchar(10);comment:商品订货号"`
	Quantity      int    `json:"quantity" gorm:"type:int unsigned;comment:商品数量"`
	Selected      int    `json:"selected" gorm:"type:tinyint(1);comment:选中标记"`
	models.ModelTime
	models.ControlBy
}

const (
	UserCartSelected0 = 0
	UserCartSelected1 = 1
)

func (UserCart) TableName() string {
	return "user_cart"
}

func (e *UserCart) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *UserCart) GetId() interface{} {
	return e.Id
}

func (e *UserCart) GetUserCartByGoodsIdAndWarehouse(tx *gorm.DB, userId, goodsId int, warehouseCode string) error {
	pcPrefix := global.GetTenantPcDBNameWithDB(tx)
	return tx.Table(pcPrefix+"."+e.TableName()).Where("user_id = ?", userId).Where("goods_id = ?", goodsId).
		Where("warehouse_code = ?", warehouseCode).
		Take(e).Error
}

func (e *UserCart) GetUserCartWithErr(tx *gorm.DB) error {
	err := e.GetUserCartByGoodsIdAndWarehouse(tx, e.UserId, e.GoodsId, e.WarehouseCode)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("未找到购物车信息")
	}
	return err
}

func CheckByGoodsIdAndWarehouseCode(tx *gorm.DB, goodsId int, warehouseCode string) *Goods {
	goods := &Goods{}
	pcPrefix := global.GetTenantPcDBNameWithDB(tx)
	_ = tx.Table(pcPrefix+"."+goods.TableName()).Where("id = ?", goodsId).
		Where("warehouse_code = ?", warehouseCode).
		Where("status = ?", 1).
		Where("online_status = ?", 1).
		Take(goods)
	return goods
}

func CheckBySkuAndWarehouseCode(tx *gorm.DB, sku, warehouseCode string) *Goods {
	goods := &Goods{}
	pcPrefix := global.GetTenantPcDBNameWithDB(tx)
	_ = tx.Table(pcPrefix+"."+goods.TableName()).Where("sku_code = ?", sku).
		Where("warehouse_code = ?", warehouseCode).
		Where("status = ?", 1).
		Where("online_status = ?", 1).
		Take(goods)
	return goods
}

// 清空用户购物车通过goodsIds

func UserCartRemoveByGoodsIds(tx *gorm.DB, userId int, warehouseCode string, goodsId []int) error {
	userCart := &UserCart{}
	pcPrefix := global.GetTenantPcDBNameWithDB(tx)
	return tx.Table(pcPrefix+"."+userCart.TableName()).Where("user_id = ?", userId).
		Where("warehouse_code = ?", warehouseCode).
		Where("goods_id in ?", goodsId).
		Delete(userCart).Error
}

func (e *UserCart) CheckSalesMoq(tx *gorm.DB) error {
	userCartGoodsProduct := &UserCartGoodsProduct{}
	goodsResult, _ := userCartGoodsProduct.GetByGoodsIds(tx, []int{e.GoodsId}, "")
	if len(*goodsResult) == 0 {
		return errors.New("未找到该商品")
	}
	if e.Quantity < (*goodsResult)[0].SalesMoq {
		return errors.New("数量小于最小起订量")
	}
	return nil
}

func (e *UserCart) Add(tx *gorm.DB, quantity int) error {
	pcPrefix := global.GetTenantPcDBNameWithDB(tx)
	goods := CheckByGoodsIdAndWarehouseCode(tx, e.GoodsId, e.WarehouseCode)
	if goods.Id == 0 {
		return errors.New("该商品在当前仓库中，已下架或不存在")
	}
	err := e.GetUserCartByGoodsIdAndWarehouse(tx, e.UserId, e.GoodsId, e.WarehouseCode)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		e.SkuCode = goods.SkuCode
		e.Quantity = quantity
	} else {
		if err != nil {
			return err
		}
		e.Selected = UserCartSelected1
		e.Quantity += quantity
		e.UpdateBy = e.CreateBy
		e.UpdateByName = e.CreateByName
	}
	if err := e.CheckSalesMoq(tx); err != nil {
		return err
	}
	return tx.Table(pcPrefix + "." + e.TableName()).Save(e).Error
}

func (e *UserCart) AddBySku(tx *gorm.DB, quantity int) error {
	pcPrefix := global.GetTenantPcDBNameWithDB(tx)
	goods := CheckBySkuAndWarehouseCode(tx, e.SkuCode, e.WarehouseCode)
	if goods.Id == 0 {
		return errors.New(e.SkuCode + ",商品在当前仓库中，已下架或不存在")
	}
	err := e.GetUserCartByGoodsIdAndWarehouse(tx, e.UserId, goods.Id, e.WarehouseCode)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		e.GoodsId = goods.Id
		e.Quantity = quantity
	} else {
		if err != nil {
			return err
		}
		e.Quantity += quantity
		e.UpdateBy = e.CreateBy
		e.UpdateByName = e.CreateByName
	}
	e.Selected = UserCartSelected1
	return tx.Table(pcPrefix + "." + e.TableName()).Save(e).Error
}

func (e *UserCart) Edit(tx *gorm.DB, quantity int) error {
	err := e.GetUserCartByGoodsIdAndWarehouse(tx, e.UserId, e.GoodsId, e.WarehouseCode)
	if err != nil {
		return err
	}
	e.Selected = UserCartSelected1
	e.Quantity = quantity
	e.UpdateBy = e.CreateBy
	e.UpdateByName = e.CreateByName

	if err := e.CheckSalesMoq(tx); err != nil {
		return err
	}

	return tx.Save(e).Error
}

func (e *UserCart) Remove(tx *gorm.DB) error {
	err := e.GetUserCartWithErr(tx)
	if err != nil {
		return err
	}
	return tx.Delete(e).Error
}

func (e *UserCart) SelectOne(tx *gorm.DB) error {
	err := e.GetUserCartWithErr(tx)
	if err != nil {
		return err
	}
	if e.Selected == 0 {
		e.Selected = 1
	} else {
		e.Selected = 0
	}
	return tx.Save(e).Error
}

func (e *UserCart) SelectAll(tx *gorm.DB, userId int, warehouseCode string) error {
	return tx.Model(e).Where("user_id = ?", userId).
		Where("warehouse_code = ?", warehouseCode).
		Update("selected", 1).Error
}

func (e *UserCart) UnSelectAll(tx *gorm.DB, userId int, warehouseCode string) error {
	return tx.Model(e).Where("user_id = ?", userId).
		Where("warehouse_code = ?", warehouseCode).
		Update("selected", 0).Error
}

func (e *UserCart) ClearSelect(tx *gorm.DB, userId int, warehouseCode string, goodsIdsSaleStatus0 []int) error {
	pcPrefix := global.GetTenantPcDBNameWithDB(tx)
	return tx.Table(pcPrefix+"."+e.TableName()).Where("user_id = ?", userId).
		Where("warehouse_code = ?", warehouseCode).
		Where("selected = ?", 1).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if len(goodsIdsSaleStatus0) != 0 {
				return db.Where("goods_id not in ?", goodsIdsSaleStatus0)
			}
			return db
		}).
		Delete(e).Error
}

func (e *UserCart) GetUserCartCondition(tx *gorm.DB, userId int, warehouseCode string, condition func(*gorm.DB) *gorm.DB) (*[]UserCart, error) {
	pcPrefix := global.GetTenantPcDBNameWithDB(tx)
	userCart := &[]UserCart{}
	err := tx.Table(pcPrefix+"."+e.TableName()).Where("user_id = ?", userId).
		Where("warehouse_code = ?", warehouseCode).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if condition != nil {
				return condition(db)
			}
			return db
		}).
		Find(userCart).Error
	return userCart, err
}

func (e *UserCart) GetUserCartSeleted(tx *gorm.DB, userId int, warehouseCode string) (*[]UserCart, error) {
	return e.GetUserCartCondition(tx, userId, warehouseCode, func(db *gorm.DB) *gorm.DB {
		return db.Where("selected = ?", 1)
	})
}

func (e *UserCart) GetUserCartAll(tx *gorm.DB, userId int, warehouseCode string) (*[]UserCart, error) {
	return e.GetUserCartCondition(tx, userId, warehouseCode, nil)
}

type UserCartPageGoodsProduct struct {
	UserCartGoodsProduct
	Collected bool `json:"collected"`
}
type UserCartGoodsProduct struct {
	GoodsId             int           `json:"goodsId"`
	NameZh              string        `json:"nameZh"`
	MarketPrice         float64       `json:"marketPrice"`
	NakedSalePrice      float64       `json:"nakedSalePrice" gorm:"-"`
	VendorId            int           `json:"vendorId"`
	VendorName          string        `json:"vendorName" gorm:"-"`
	SkuCode             string        `json:"skuCode"`
	WarehouseCode       string        `json:"warehouseCode"`
	Status              int           `json:"status"`
	OnlineStatus        int           `json:"onlineStatus"`
	Stock               int           `json:"stock"`
	ProductNo           string        `json:"productNo"`
	MfgModel            string        `json:"mfgModel"`
	RefundFlag          int           `json:"refundFlag"`
	BrandZh             string        `json:"brandZh"`
	Tax                 string        `json:"tax"`
	Image               MediaInstance `json:"image" gorm:"-"`
	Quantity            int           `json:"quantity" gorm:"-"`
	SalesMoq            int           `json:"salesMoq"`
	TotalNakedSalePrice float64       `json:"totalNakedSalePrice" gorm:"-"`
	TotalMarketPrice    float64       `json:"totalMarketPrice" gorm:"-"`
	TotalTaxPrice       float64       `json:"totalTaxPrice" gorm:"-"`
	Selected            int           `json:"selected" gorm:"-"`
	ShowCartProduct     int           `json:"showCartProduct" gorm:"-"`
	SaleStatus          int           `json:"saleStatus" gorm:"-"`
	SaleStatusRemark    string        `json:"saleStatusRemark" gorm:"-"`
}

func (e *UserCartGoodsProduct) GetByGoodsIds(tx *gorm.DB, goodsIds []int, logicWarehouseType string) (*[]UserCartGoodsProduct, error) {
	var data Goods
	var goodsResult = &[]UserCartGoodsProduct{}
	pcPrefix := global.GetTenantPcDBNameWithDB(tx)
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	ucPrefix := global.GetTenantUcDBNameWithDB(tx)
	if err := tx.Table(pcPrefix+"."+data.TableName()).Debug().
		Joins("LEFT JOIN "+pcPrefix+".product ON goods.sku_code = product.sku_code").
		Joins("LEFT JOIN "+pcPrefix+".brand ON product.brand_id = brand.id").
		Joins("LEFT JOIN "+wcPrefix+".stock_info ON goods.sku_code = stock_info.sku_code and goods.warehouse_code = stock_info.warehouse_code").
		Joins("LEFT JOIN "+wcPrefix+".warehouse ON warehouse.warehouse_code = goods.warehouse_code").
		Joins("LEFT JOIN "+wcPrefix+".logic_warehouse ON stock_info.warehouse_code = logic_warehouse.warehouse_code and stock_info.logic_warehouse_code = logic_warehouse.logic_warehouse_code").
		Joins("LEFT JOIN "+ucPrefix+".company_info ON company_info.id = warehouse.company_id").
		Select([]string{
			"goods.id as GoodsId",
			"goods.sku_code",
			"goods.market_price",
			"goods.product_no",
			"goods.online_status",
			"goods.status",
			"goods.vendor_id",
			"goods.warehouse_code",
			"product.name_zh",
			"product.tax",
			"product.mfg_model",
			"product.refund_flag",
			"product.sales_moq",
			"brand.brand_zh",
			"stock_info.stock",
			"warehouse.company_id",
			"company_info.check_stock_status",
		}).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if logicWarehouseType != "" {
				db.Where("logic_warehouse.type = ?", logicWarehouseType)
			}
			return db
		}).
		Where("goods.id in ?", goodsIds).
		Find(&goodsResult).Error; err != nil {
		return goodsResult, err
	}
	return goodsResult, nil
}
