package models

import (
	"go-admin/common/models"
	"gorm.io/gorm"
)

type ProductCategory struct {
	models.Model

	SkuCode      string     `json:"skuCode" gorm:"type:varchar(10);comment:sku"`
	CategoryId   int        `json:"categoryId" gorm:"type:int unsigned;comment:产线id"`
	MainCateFlag int        `json:"mainCateFlag" gorm:"type:tinyint(1);comment:主产线标志0否1是"`
	Path         []Category `json:"path" gorm:"-"`
	models.ModelTime
	models.ControlBy
}

func (ProductCategory) TableName() string {
	return "product_category"
}

func (e *ProductCategory) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *ProductCategory) GetId() interface{} {
	return e.Id
}

func (e *ProductCategory) GetMainCategoryCountBySkus(tx *gorm.DB, skus []string) int {
	var count int64
	tx.Model(e).Where("sku_code in ?", skus).Where("main_cate_flag = ?", 1).Distinct("category_id").Count(&count)
	return int(count)
}
