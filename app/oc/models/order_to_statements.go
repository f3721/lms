package models

import
	"go-admin/common/models"

type OrderToStatements struct {
	models.Model

	StatementsId int `json:"statementsId" gorm:"type:int(11);comment:对账单id"`
	OrderId string `json:"orderId" gorm:"type:varchar(30);comment:订单id"`
}

func (OrderToStatements) TableName() string {
	return "order_to_statements"
}