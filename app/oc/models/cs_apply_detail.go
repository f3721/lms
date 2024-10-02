package models

import (
	"go-admin/common/models"
)

type CsApplyDetail struct {
	models.Model

	CsNo          string  `json:"csNo" gorm:"type:varchar(30);comment:售后申请编号"`                  // 售后申请编号
	ProductModel  string  `json:"productModel" gorm:"type:varchar(200);comment:售后申请产品名原始型号"`    // 售后申请产品名原始型号
	SkuCode       string  `json:"skuCode" gorm:"type:varchar(30);comment:商品sku"`                // 商品sku
	ProductName   string  `json:"productName" gorm:"type:varchar(1000);comment:商品名称"`           // 商品名称
	Quantity      int     `json:"quantity" gorm:"type:int;comment:数量"`                          // 数量
	Unit          string  `json:"unit" gorm:"type:varchar(50);comment:单位"`                      // 单位
	RefundAmt     string  `json:"refundAmt" gorm:"type:decimal(10,2);comment:退款金额"`             // 退款金额
	ReparationAmt string  `json:"reparationAmt" gorm:"type:decimal(10,2);comment:赔款金额"`         // 赔款金额
	Remark        string  `json:"remark" gorm:"type:mediumtext;comment:备注"`                     // 备注
	CsType        int     `json:"csType" gorm:"type:tinyint unsigned;comment:类型：0-用户申请、1-实际售后"` // 类型：0-用户申请、1-实际售后
	SalePrice     float64 `json:"salePrice" gorm:"type:decimal(12,4);comment:销售价"`              // 销售价
	ProductPic    string  `json:"productPic" gorm:"type:varchar(200);comment:商品图片url"`          // 商品图片url
	BrandName     string  `json:"brandName" gorm:"type:varchar(50);comment:产品的品牌"`              // 产品的品牌
	DeliveryTime  string  `json:"deliveryTime" gorm:"type:varchar(20);comment:货期"`              // 货期
	WarehouseCode string  `json:"warehouseCode" gorm:"type:varchar(20);comment:到货仓库"`           // 到货仓库
	GoodsId       int     `json:"goodsId" gorm:"type:int unsigned;comment:GoodsId"`             //
	VendorName    string  `json:"vendorName" gorm:"type:varchar(255);comment:货主"`               // 货主
	VendorSkuCode string  `json:"vendorSkuCode" gorm:"type:varchar(255);comment:货主sku"`         // 货主sku
	ProductNo     string  `json:"productNo" gorm:"type:varchar(255);comment:物料编码"`              // 物料编码
	IsDefective   int     `json:"isDefective" gorm:"type:tinyint unsigned;comment:是否次品 1是 0否"`  // 是否次品 1是 0否
	Pics          string  `json:"pics" gorm:"type:mediumtext;comment:售后申请的图片"`                  // 售后申请的图片
	models.ModelTime
	//models.ControlBy
}

func (CsApplyDetail) TableName() string {
	return "cs_apply_detail"
}

func (e *CsApplyDetail) Generate() *CsApplyDetail {
	o := *e
	return &o
}

func (e *CsApplyDetail) GetId() interface{} {
	return e.Id
}
