package models

import (
     
     
     
     

	"go-admin/common/models"

)

type UserCollect struct {
    models.Model
    
    SkuCode string `json:"skuCode" gorm:"type:varchar(30);comment:产品SKU"`// 产品SKU 
    UserId int `json:"userId" gorm:"type:int unsigned;comment:用户ID"`// 用户ID 
    GoodsId int `json:"goodsId" gorm:"type:int unsigned;comment:商品表id"`// 商品表id 
    WarehouseCode string `json:"warehouseCode" gorm:"type:varchar(20);comment:仓库code"`// 仓库code 
    models.ModelTime
    models.ControlBy
}

func (UserCollect) TableName() string {
    return "user_collect"
}

func (e *UserCollect) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *UserCollect) GetId() interface{} {
	return e.Id
}
