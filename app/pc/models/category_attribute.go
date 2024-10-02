package models

import (
	"go-admin/common/models"
)

type CategoryAttribute struct {
	models.Model

	CategoryId   int    `json:"categoryId" gorm:"type:int;comment:产线id"`
	AttributeId  int    `json:"attributeId" gorm:"type:int;comment:属性id"`
	Seq          int    `json:"seq" gorm:"type:smallint unsigned;comment:序列"`
	RequiredFlag int    `json:"requiredFlag" gorm:"type:smallint;comment:必填标志"`
	FilterFlag   int    `json:"filterFlag" gorm:"type:smallint;comment:筛选标志(0不筛选1值筛选2范围筛选)"`
	RangeVal     string `json:"rangeVal" gorm:"type:varchar(1024);comment:范围值"`
	Attribute    AttributeDef
	models.ModelTime
	models.ControlBy
}

func (CategoryAttribute) TableName() string {
	return "category_attribute"
}

func (e *CategoryAttribute) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *CategoryAttribute) GetId() interface{} {
	return e.Id
}
