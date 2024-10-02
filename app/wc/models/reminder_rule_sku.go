package models

import (
	"go-admin/common/models"
)

type ReminderRuleSku struct {
	models.Model

	ReminderRuleId     int    `json:"reminderRuleId" gorm:"type:int(11);comment:sxyz_reminder_rule 补货提醒规则表id"`
	SkuCode            string `json:"skuCode" gorm:"type:varchar(10);comment:sku"`
	WarningValue       int    `json:"warningValue" gorm:"type:int(11);comment:预警值"`
	ReplenishmentValue int    `json:"replenishmentValue" gorm:"type:int(11);comment:设置的备货量"`
	Status             string `json:"status" gorm:"type:smallint(1) unsigned;comment:状态 1启用 0未启用"`
	models.ModelTime
	models.ControlBy
}

func (ReminderRuleSku) TableName() string {
	return "reminder_rule_sku"
}

func (e *ReminderRuleSku) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *ReminderRuleSku) GetId() interface{} {
	return e.Id
}
