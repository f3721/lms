package models

import (
	"go-admin/common/models"
	"time"
)

type ReminderRuleSkuLog struct {
	models.Model

	ReminderRuleSkuId  int       `json:"reminderRuleSkuId" gorm:"type:int(11);comment:sku补货提醒规则表id"`
	SkuCode            string    `json:"skuCode" gorm:"type:varchar(10);comment:sku"`
	WarningValue       int       `json:"warningValue" gorm:"type:int(11);comment:预警值"`
	ReplenishmentValue int       `json:"replenishmentValue" gorm:"type:int(11);comment:设置的备货量"`
	Status             string    `json:"status" gorm:"type:smallint(1) unsigned;comment:状态 1启用 0未启用"`
	LogType            int       `json:"logType" gorm:"type:varchar(10);comment:操作类型 1创建 2修改 3删除"`
	CreateByName       string    `json:"createByName" gorm:"type:varchar(255);comment:CreateByName"`
	CreatedAt          time.Time `json:"createdAt" gorm:"comment:创建时间"`
	CreateBy           int       `json:"createBy" gorm:"index;comment:创建者"`
}

func (ReminderRuleSkuLog) TableName() string {
	return "reminder_rule_sku_log"
}

//func (e *ReminderRuleSkuLog) Generate() models.ActiveRecord {
//	o := *e
//	return &o
//}

func (e *ReminderRuleSkuLog) GetId() interface{} {
	return e.Id
}
