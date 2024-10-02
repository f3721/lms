package models

import (

	"go-admin/common/models"

)

type SkuClassification struct {
    models.Model
    
    CompanyId int `json:"companyId" gorm:"type:int;comment:公司ID"` 
    SkuCode string `json:"skuCode" gorm:"type:varchar(20);comment:产品sku"` 
    Status int `json:"status" gorm:"type:tinyint;comment:是否启用 0.否 1.是 默认1"` 
    ClassificationId int `json:"classificationId" gorm:"type:int;comment:客户分类"` 
    models.ModelTime
    models.ControlBy
}

func (SkuClassification) TableName() string {
    return "sku_classification"
}

func (e *SkuClassification) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SkuClassification) GetId() interface{} {
	return e.Id
}