package models

import "gorm.io/gorm"

type CategoryPath struct {
	CategoryId int `json:"categoryId" gorm:"type:int unsigned;comment:CategoryId"`
	PathId     int `json:"pathId" gorm:"type:int unsigned;comment:PathId"`
	Level      int `json:"level" gorm:"type:int unsigned;comment:Level"`
}

func (CategoryPath) TableName() string {
	return "category_path"
}

func (e *CategoryPath) GetSkuCodeByCateId(db *gorm.DB, categoryId int) []string {
	var skuCode []string
	db.Model(&CategoryPath{}).
		Select("sku_code").
		Joins("INNER JOIN product_category on category_path.category_id = product_category.category_id").
		Where("path_id in (?)", categoryId).
		Group("sku_code").
		Scan(&skuCode)
	return skuCode
}

func (e *CategoryPath) GetCategoryIdsByCateId(db *gorm.DB, categoryId int) []int {
	var categoryIds []int
	db.Model(&CategoryPath{}).
		Select("category_path.category_id").
		Joins("INNER JOIN product_category on category_path.category_id = product_category.category_id").
		Where("category_path.path_id in (?)", categoryId).
		Group("category_id").
		Scan(&categoryIds)
	return categoryIds
}
