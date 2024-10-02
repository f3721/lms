package models

import (
	"go-admin/common/models"
)

type AttributeConfig struct {
	models.Model

	Type      string `json:"type" gorm:"type:varchar(32);comment:类型"`
	Key       string `json:"key" gorm:"type:varchar(32);comment:配置键"`
	Value     string `json:"value" gorm:"type:varchar(100);comment:配置值"`
	SortOrder int    `json:"sortOrder" gorm:"type:int(11);comment:SortOrder"`
}

func (AttributeConfig) TableName() string {
	return "attribute_config"
}

func (e *AttributeConfig) GetId() interface{} {
	return e.Id
}
