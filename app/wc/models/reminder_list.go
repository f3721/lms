package models

import (
	"go-admin/common/models"
	"time"
)

type ReminderList struct {
	models.Model

	ReminderRuleId int       `json:"reminderRuleId" gorm:"type:int(11);comment:sxyz_reminder_rule 补货规则表id"`
	CompanyId      int       `json:"companyId" gorm:"type:int(11);comment:公司id"`
	WarehouseCode  string    `json:"warehouseCode" gorm:"type:int(11);comment:仓库id"`
	VendorId       int       `json:"vendorId" gorm:"type:int(11);comment:货主id"`
	SkuCount       int       `json:"skuCount" gorm:"type:int(11);comment:SKU数量"`
	CreatedAt      time.Time `json:"createdAt" gorm:"comment:创建时间"`
}

func (ReminderList) TableName() string {
	return "reminder_list"
}

func (e *ReminderList) GetId() interface{} {
	return e.Id
}
