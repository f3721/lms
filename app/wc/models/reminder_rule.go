package models

import (
	"go-admin/common/models"
)

type ReminderRule struct {
	models.Model

	CompanyId          int    `json:"companyId" gorm:"type:int(11);comment:公司id"` //公司id
	WarehouseCode      string `json:"warehouseCode" gorm:"type:int(11);comment:仓库id"`
	WarningValue       int    `json:"warningValue" gorm:"type:int(11);comment:SKU通用预警值"`
	ReplenishmentValue int    `json:"replenishmentValue" gorm:"type:int(11);comment:设置的备货量"`
	Status             int    `json:"status" gorm:"type:smallint(1) unsigned;comment:状态 1启用 0未启用"`
	models.ModelTime
	models.ControlBy
}

func (ReminderRule) TableName() string {
	return "reminder_rule"
}

func (e *ReminderRule) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *ReminderRule) GetId() interface{} {
	return e.Id
}
