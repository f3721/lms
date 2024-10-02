package models

import (
	"github.com/samber/lo"
	"go-admin/common/global"
	"go-admin/common/models"
	"gorm.io/gorm"
)

const (
	SupplierModelName    = "supplier"
	SupplierModelInsert  = "insert"
	SupplierModelUpdate  = "update"
	SupplierModelStatus0 = "0"
	SupplierModelStatus1 = "1"
)

var SupplierModelStatus = map[string]string{
	SupplierModelStatus0: "无效",
	SupplierModelStatus1: "有效",
}

type Supplier struct {
	models.Model

	Code          string `json:"code" gorm:"type:varchar(20);comment:供应商编码"`
	NameZh        string `json:"nameZh" gorm:"type:varchar(50);comment:供应商中文名"`
	NameEn        string `json:"nameEn" gorm:"type:varchar(50);comment:供应商英文名"`
	ShortName     string `json:"shortName" gorm:"type:varchar(20);comment:供应商简称"`
	PostalCode    string `json:"postalCode" gorm:"type:varchar(20);comment:邮编"`
	Linkman       string `json:"linkman" gorm:"type:varchar(20);comment:联系人"`
	Phone         string `json:"phone" gorm:"type:varchar(20);comment:手机"`
	Email         string `json:"email" gorm:"type:varchar(50);comment:邮箱"`
	Fax           string `json:"fax" gorm:"type:varchar(50);comment:传真"`
	Address       string `json:"address" gorm:"type:varchar(200);comment:详细地址"`
	CityId        int    `json:"cityId" gorm:"type:int unsigned;comment:市id"`
	ProvinceId    int    `json:"provinceId" gorm:"type:int unsigned;comment:省id"`
	CountryId     int    `json:"countryId" gorm:"type:int unsigned;comment:国家id"`
	Telephone     string `json:"telephone" gorm:"type:varchar(50);comment:电话"`
	Remark        string `json:"remark" gorm:"type:varchar(255);comment:备注"`
	Status        string `json:"status" gorm:"type:tinyint(1);comment:0无效 1有效"`
	BackupLinkman string `json:"backupLinkman" gorm:"type:varchar(50);comment:后备联系人"`
	BackupPhone   string `json:"backupPhone" gorm:"type:varchar(50);comment:后备联系电话"`
	models.ModelTime
	models.ControlBy
}

func (Supplier) TableName() string {
	return "supplier"
}

func (e *Supplier) Generate() *Supplier {
	o := *e
	return &o
}

func (e *Supplier) GetId() interface{} {
	return e.Id
}

func (e *Supplier) CheckNameZh(tx *gorm.DB, name string, id int) (bool, error) {
	var result struct {
		Id int
	}
	tx = tx.Model(&Supplier{}).Select("id")
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

func (e *Supplier) CheckCode(tx *gorm.DB, code string, id int) (bool, error) {
	var result struct {
		Id int
	}
	tx = tx.Model(&Supplier{}).Select("id")
	if id != 0 {
		tx.Where("id <> ?", id)
	}
	if err := tx.Where("code = ?", code).Scan(&result).Error; err != nil {
		return false, err
	}
	if result.Id != 0 {
		return false, nil
	}
	return true, nil
}

func (e *Supplier) GetById(tx *gorm.DB, id int) error {
	return tx.Take(e, id).Error
}

func (e *Supplier) GetByNameZh(tx *gorm.DB, name string) (*Supplier, error) {
	var supplier = &Supplier{}
	err := tx.Where("name_zh = ?", name).First(&supplier).Error
	if supplier.Id != 0 {
		return supplier, nil
	}
	return supplier, err
}

func (e *Supplier) FindOneByName(tx *gorm.DB, nameZh string) (*Supplier, error) {
	var data Supplier
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	err := tx.Table(wcPrefix+"."+e.TableName()).Where("name_zh = ?", nameZh).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func GetSupplierListByIds(tx *gorm.DB, ids []int) (*[]Supplier, error) {
	var suppliers = &[]Supplier{}
	err := tx.Find(suppliers, ids).Error
	return suppliers, err
}

func GetSupplierMapByIds(tx *gorm.DB, ids []int) map[int]string {
	supplierSlice, _ := GetSupplierListByIds(tx, ids)
	supplierIdMap := lo.Associate(*supplierSlice, func(f Supplier) (int, string) {
		return f.Id, f.NameZh
	})
	return supplierIdMap
}
