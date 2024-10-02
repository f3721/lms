package models

import (

	"go-admin/common/models"

)

type ProductExtAttribute struct {
    models.Model
    
    SkuCode string `json:"skuCode" gorm:"type:varchar(10);comment:产品SKU"` 
    AttributeId int `json:"attributeId" gorm:"type:int unsigned;comment:属性ID"` 
    ValueZh string `json:"valueZh" gorm:"type:varchar(50);comment:属性值(中文)"` 
    ValueEn string `json:"valueEn" gorm:"type:varchar(50);comment:属性值(英文)"` 
    Status int `json:"status" gorm:"type:tinyint(1);comment:维护状态"` 
    models.ModelTime
    models.ControlBy
}

func (ProductExtAttribute) TableName() string {
    return "product_ext_attribute"
}

func (e *ProductExtAttribute) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *ProductExtAttribute) GetId() interface{} {
	return e.Id
}
