package models

import (

	"go-admin/common/models"

)

type UserPayAccount struct {
    models.Model
    
    CompanyId int `json:"companyId" gorm:"type:int;comment:公司ID"` 
    UserId int `json:"userId" gorm:"type:int;comment:用户ID"` 
    PayAccount string `json:"payAccount" gorm:"type:varchar(100);comment:支付账户"` 
    ClassificationId int `json:"classificationId" gorm:"type:int;comment:客户分类"` 
    models.ModelTime
    models.ControlBy
}

func (UserPayAccount) TableName() string {
    return "user_pay_account"
}

func (e *UserPayAccount) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *UserPayAccount) GetId() interface{} {
	return e.Id
}