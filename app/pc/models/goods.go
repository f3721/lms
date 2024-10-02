package models

import (
	"go-admin/common/global"
	"go-admin/common/models"
	"gorm.io/gorm"
)

type Goods struct {
	models.Model
	SkuCode           string  `json:"skuCode" gorm:"type:varchar(10);comment:商品SKU"`
	SupplierSkuCode   string  `json:"supplierSkuCode" gorm:"type:varchar(20);comment:货主SKU"`
	WarehouseCode     string  `json:"warehouseCode" gorm:"type:varchar(20);comment:仓库ID"`
	ProductNo         string  `json:"productNo" gorm:"type:varchar(20);comment:物料编码"`
	MarketPrice       float64 `json:"marketPrice" gorm:"type:decimal(10,2);comment:价格"`
	PriceModifyReason string  `json:"priceModifyReason" gorm:"type:varchar(50);comment:价格调整备注"`
	ApproveStatus     int     `json:"approveStatus" gorm:"type:tinyint(1);comment:审核状态 0 待审核  1 审核通过  2 审核失败"`
	ApproveRemark     string  `json:"approveRemark" gorm:"type:varchar(255);comment:审核备注"`
	Status            int     `json:"status" gorm:"type:tinyint(1);comment:商品状态 0 禁用  1启用"`
	OnlineStatus      int     `json:"onlineStatus" gorm:"type:tinyint(1);comment:上下架状态  0未上架 1上架  2下架"`
	VendorId          int     `json:"vendorId" gorm:"type:int unsigned;comment:供应商ID"`
	Product           Product `json:"product" gorm:"foreignKey:SkuCode;references:SkuCode"`
	models.ModelTime
	models.ControlBy
}

func (Goods) TableName() string {
	return "goods"
}

func (e *Goods) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *Goods) GetId() interface{} {
	return e.Id
}

func (e *Goods) GetGoodsBySku(tx *gorm.DB, skuCodes []string, warehouseCode string, vendorId int) (goods *[]Goods, err error) {
	pcPrefix := global.GetTenantPcDBNameWithDB(tx)
	err = tx.Table(pcPrefix+".goods").Where(map[string]interface{}{"sku_code": skuCodes}).Where("warehouse_code = ?", warehouseCode).Where("vendor_id = ?", vendorId).Find(&goods).Error
	return goods, err
}
