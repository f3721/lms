package models

import (
	"go-admin/common/models"
	"gorm.io/gorm"
)

type Brand struct {
	models.Model

	BrandZh          string `json:"brandZh" gorm:"type:varchar(100);comment:品牌中文"`
	BrandEn          string `json:"brandEn" gorm:"type:varchar(100);comment:品牌英文"`
	FirstLetter      string `json:"firstLetter" gorm:"type:chat(1);comment:首字母"`
	BrandDescription string `json:"brandDescription" gorm:"type:varchar(600);comment:品牌描述"`
	Status           int    `json:"status" gorm:"type:tinyint(1);comment:激活状态(0:不使用1:使用)"`
	models.ModelTime
	models.ControlBy
}

func (Brand) TableName() string {
	return "brand"
}

func (e *Brand) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *Brand) GetId() interface{} {
	return e.Id
}

func (e *Brand) CheckNameZh(tx *gorm.DB, name string, id int) (bool, error) {
	var result struct {
		Id int
	}
	tx = tx.Model(&Brand{}).Select("id")
	if id != 0 {
		tx.Where("id <> ?", id)
	}
	if err := tx.Where("brand_zh = ?", name).Scan(&result).Error; err != nil {
		return false, err
	}
	if result.Id != 0 {
		return false, nil
	}
	return true, nil
}

func (e *Brand) CheckName(tx *gorm.DB, nameZh string, nameEn string, id int) (bool, error) {
	var result struct {
		Id int
	}
	tx = tx.Model(&Brand{}).Select("id")
	if id != 0 {
		tx.Where("id <> ?", id)
	}
	if err := tx.Where("brand_zh = ?", nameZh).Where("brand_en = ?", nameEn).Scan(&result).Error; err != nil {
		return false, err
	}
	if result.Id != 0 {
		return false, nil
	}
	return true, nil
}

func (e *Brand) FindOneByName(tx *gorm.DB, prefixDb string, nameZh string) (*Brand, error) {
	var data Brand
	err := tx.Table(prefixDb+"."+e.TableName()).Where("brand_zh = ?", nameZh).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (e *Brand) Insert(tx *gorm.DB, prefixDb string, req *Brand) error {
	err := tx.Table(prefixDb + "." + e.TableName()).Create(req).Error
	if err != nil {
		return err
	}
	return nil
}
