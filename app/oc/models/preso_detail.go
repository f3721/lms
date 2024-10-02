package models

import (

	"go-admin/common/models"
    "time"
)

type PresoDetail struct {
    models.Model

    PresoNo string `json:"presoNo" gorm:"type:varchar(30);comment:审批单编号"`
    UserId int `json:"userId" gorm:"type:int unsigned;comment:用户ID"` 
    UserName string `json:"userName" gorm:"type:varchar(200);comment:用户名称"` 
    GoodsId int `json:"goodsId" gorm:"type:int unsigned;comment:goods表ID"` 
    SkuCode string `json:"skuCode" gorm:"type:varchar(20);comment:sku"` 
    VendorId int `json:"vendorId" gorm:"type:int unsigned;comment:货主ID"` 
    ProductNo string `json:"productNo" gorm:"type:varchar(30);comment:物料编码"` 
    WarehouseCode string `json:"warehouseCode" gorm:"type:varchar(10);comment:发货仓"` 
    MarketPrice float64 `json:"marketPrice" gorm:"type:decimal(10,2);comment:系统价格"`
    SalePrice float64 `json:"salePrice" gorm:"type:decimal(10,2);comment:商品销售价格"`
    Quantity int `json:"quantity" gorm:"type:int unsigned;comment:商品数量"` 
    ApproveQuantity int `json:"approveQuantity" gorm:"type:int unsigned;comment:商品审核通过的数量"`
    ProductName string `json:"productName" gorm:"type:varchar(512);comment:商品名称"`
    ProductPic string `json:"productPic" gorm:"type:varchar(200);comment:商品名称"`
    ApproveStatus int `json:"approveStatus" gorm:"type:tinyint;comment:审批状态 -1 审批不通过  1 审批通过 0 初始提交审批 10 审批中 -2 超时 -3 撤回"`
    ApproveRemark string `json:"approveRemark" gorm:"type:varchar(400);comment:审批备注"`
    RejectByName string `json:"rejectByName" gorm:"index;comment:驳回人"`
    Step int `json:"step" gorm:"type:tinyint unsigned;comment:审批到第几级"` 
    UserProductRemark string `json:"userProductRemark" gorm:"type:varchar(255);comment:客户商品物料采购备注"`
    CreatedAt time.Time      `json:"createdAt" gorm:"comment:创建时间"`
    UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:最后更新时间"`
    models.ControlBy
}

func (PresoDetail) TableName() string {
    return "preso_detail"
}

func (e *PresoDetail) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *PresoDetail) GetId() interface{} {
	return e.Id
}
