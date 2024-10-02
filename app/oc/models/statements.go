package models

import (
	"time"

	"go-admin/common/models"
)

type Statements struct {
	models.Model

	StatementsNo string    `json:"statementsNo" gorm:"type:varchar(10);comment:对账单编号"`
	CompanyId    int       `json:"companyId" gorm:"type:int(11);comment:公司"`
	TotalAmount  float64   `json:"totalAmount" gorm:"type:decimal(10,2);comment:账单总金额"`
	StartTime    time.Time `json:"startTime" gorm:"type:datetime;comment:账单开始日期"`
	EndTime      time.Time `json:"endTime" gorm:"type:datetime;comment:账单结束日期"`
	OrderCount   int       `json:"orderCount" gorm:"type:int(11);comment:订单总数"`
	ProductCount int       `json:"productCount" gorm:"type:int(11);comment:商品总数"`
	CreatedAt    time.Time `json:"createdAt" gorm:"comment:创建时间"`
}

func (Statements) TableName() string {
	return "statements"
}

//func (e *Statements) Generate() models.ActiveRecord {
//	o := *e
//	return &o
//}

func (e *Statements) GetId() interface{} {
	return e.Id
}
