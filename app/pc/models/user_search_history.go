package models

import (
	"go-admin/common/models"
)

type UserSearchHistory struct {
	models.Model

	UserId        int    `json:"userId" gorm:"type:int unsigned;comment:用户编号"`
	WarehouseCode string `json:"warehouseCode" gorm:"type:varchar(20);comment:WarehouseCode"`
	Keyword       string `json:"keyword" gorm:"type:varchar(100);comment:关键词"`
	models.ModelTime
	models.ControlBy
}

func (UserSearchHistory) TableName() string {
	return "user_search_history"
}

func (e *UserSearchHistory) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *UserSearchHistory) GetId() interface{} {
	return e.Id
}
