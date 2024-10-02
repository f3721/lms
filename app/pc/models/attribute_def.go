package models

import (
	"go-admin/common/models"
	"gorm.io/gorm"
)

type AttributeDef struct {
	models.Model

	NameZh   string `json:"nameZh" gorm:"type:varchar(100);comment:属性英文名"`
	NameEn   string `json:"nameEn" gorm:"type:varchar(100);comment:属性中文名"`
	AttrType int    `json:"attrType" gorm:"type:tinyint(1) unsigned;comment:属性类型(0.主属性1:市场属性,2:技术属性)"`
	models.ModelTime
	models.ControlBy
}

func (AttributeDef) TableName() string {
	return "attribute_def"
}

func (e *AttributeDef) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *AttributeDef) GetId() interface{} {
	return e.Id
}

func (e *AttributeDef) CheckNameZh(tx *gorm.DB, name string, id int) (bool, error) {
	var result struct {
		Id int
	}
	tx = tx.Model(&AttributeDef{}).Select("id")
	if id != 0 {
		tx.Where("id <> ?", id)
	}
	if err := tx.Where("name_zh = ?", name).Scan(&result).Error; err != nil {
		return false, err
	}
	if result.Id != 0 {
		return false, nil
	}
	return true, nil
}
