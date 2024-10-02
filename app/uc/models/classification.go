package models

import (

	"go-admin/common/models"

)

type Classification struct {
    models.Model
    
    CompanyId int `json:"companyId" gorm:"type:int;comment:公司ID"` 
    Name string `json:"name" gorm:"type:varchar(80);comment:名称"` 
    Remark string `json:"remark" gorm:"type:varchar(255);comment:备注"` 
    models.ModelTime
    models.ControlBy
}

func (Classification) TableName() string {
    return "classification"
}

func (e *Classification) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *Classification) GetId() interface{} {
	return e.Id
}