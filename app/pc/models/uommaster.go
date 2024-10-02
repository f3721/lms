package models

import (
	"go-admin/common/models"
)

type Uommaster struct {
	models.Model

	Uom         string `json:"uom" gorm:"type:varchar(10);comment:Uom"`
	Description string `json:"description" gorm:"type:varchar(10);comment:Description"`
	UomPy       string `json:"uomPy" gorm:"type:varchar(20);comment:UomPy"`
	models.ModelTime
	models.ControlBy
}

func (Uommaster) TableName() string {
	return "uommaster"
}
