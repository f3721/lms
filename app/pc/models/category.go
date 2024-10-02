package models

import (
	"go-admin/common/global"
	"go-admin/common/models"
	"gorm.io/gorm"
)

type Category struct {
	models.Model

	CateLevel         int               `json:"cateLevel" gorm:"type:int(11) unsigned;comment:层级"`
	Seq               int               `json:"seq" gorm:"type:smallint(5) unsigned;comment:序列"`
	NameZh            string            `json:"nameZh" gorm:"type:varchar(255);comment:中文名"`
	NameEn            string            `json:"nameEn" gorm:"type:varchar(255);comment:英文名"`
	ParentId          int               `json:"parentId" gorm:"type:int(11) unsigned;comment:父类id"`
	Description       string            `json:"description" gorm:"type:varchar(255);comment:描述"`
	Status            int               `json:"status" gorm:"type:tinyint(1);comment:产线状态"`
	KeyWords          string            `json:"keyWords" gorm:"type:varchar(255);comment:关键字"`
	Tax               string            `json:"tax" gorm:"type:varchar(4);comment:产线税率：默认空 值有（0.13,0.06.0.09）"`
	CategoryTaxCode   string            `json:"categoryTaxCode" gorm:"type:varchar(50);comment:产线税号"`
	MediaRelationship MediaRelationship `json:"mediaRelationship" gorm:"foreignKey:BuszId"`
	models.ModelTime
	models.ControlBy
}

func (Category) TableName() string {
	return "category"
}

func (e *Category) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *Category) GetId() interface{} {
	return e.Id
}

//type SkuStruct struct {
//	SkuCode string  `json:"skuCode"`
//}

func (e *Category) GetHasCategorySku(tx *gorm.DB, skus []string) (res []string) {
	if len(skus) == 0 {
		return
	}
	pcPrefix := global.GetTenantPcDBNameWithDB(tx)
	tx.Table(pcPrefix+".product_category t").
		Select("sku_code").
		Joins("INNER JOIN "+pcPrefix+".category c on t.category_id = c.id").
		Where("sku_code in ?", skus).
		Where("ifnull(category_id,'') != ''").
		Scan(&res)
	return
}

func (e *Category) FindOneByCategoryId(tx *gorm.DB, categoryId int) (*Category, error) {
	var data Category
	err := tx.Model(&Category{}).First(&data, categoryId).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (e *Category) FindOneByParentId(tx *gorm.DB, parentId int) (*Category, error) {
	var data Category
	err := tx.Model(&Category{}).Where("parent_id = ?", parentId).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}
