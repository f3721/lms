package models

import (
	"go-admin/common/models"
	"time"
)

type ReminderListSku struct {
	models.Model

	ReminderListId              int       `json:"reminderListId" gorm:"type:int(11);comment:sxyz_reminder_list 补货清单表id"`
	SkuCode                     string    `json:"skuCode" gorm:"type:varchar(10);comment:sku"`
	VendorId                    int       `json:"vendorId" gorm:"type:int(11);comment:货主id"`
	VendorSku                   string    `json:"vendorSku" gorm:"type:varchar(10);comment:货主sku"`
	WarningValue                int       `json:"warningValue" gorm:"type:int(11);comment:预警值"`
	ReplenishmentValue          int       `json:"replenishmentValue" gorm:"type:int(11);comment:设置的备货量"`
	RecommendReplenishmentValue int       `json:"recommendReplenishmentValue" gorm:"type:int(11);comment:建议补货量：建议补货量=备货量-当前正品用数量+订单实际缺货数量"`
	GenuineStock                int       `json:"genuineStock" gorm:"type:int(11);comment:当前正品可用数量"`
	AllStock                    int       `json:"allStock" gorm:"type:int(11);comment:当前在库数量(正品仓在库数量+次品仓在库数)"`
	OccupyStock                 int       `json:"occupyStock" gorm:"type:int(11);comment:当前占用数量"`
	OrderLackStock              int       `json:"orderLackStock" gorm:"type:int(11);comment:订单缺货数量"`
	CreatedAt                   time.Time `json:"createdAt" gorm:"comment:创建时间"`
}

func (ReminderListSku) TableName() string {
	return "reminder_list_sku"
}

func (e *ReminderListSku) GetId() interface{} {
	return e.Id
}
