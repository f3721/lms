package models

import (
	"github.com/samber/lo"
	"go-admin/common/global"
	"go-admin/common/models"
	"gorm.io/gorm"
)

const (
	VendorsModelName    = "vendors"
	VendorsModelInsert  = "insert"
	VendorsModelUpdate  = "update"
	VendorsModelStatus0 = "0"
	VendorsModelStatus1 = "1"
)

var VendorsModelStatus = map[string]string{
	VendorsModelStatus0: "无效",
	VendorsModelStatus1: "有效",
}

type Vendors struct {
	models.Model

	Code          string `json:"code" gorm:"type:varchar(20);comment:货主编码"`
	NameZh        string `json:"nameZh" gorm:"type:varchar(50);comment:货主中文名"`
	NameEn        string `json:"nameEn" gorm:"type:varchar(50);comment:供应商英文名"`
	ShortName     string `json:"shortName" gorm:"type:varchar(20);comment:供应商简称"`
	PostalCode    string `json:"postalCode" gorm:"type:varchar(20);comment:邮编"`
	Linkman       string `json:"linkman" gorm:"type:varchar(20);comment:联系人"`
	Phone         string `json:"phone" gorm:"type:varchar(20);comment:手机"`
	Email         string `json:"email" gorm:"type:varchar(50);comment:邮箱"`
	Fax           string `json:"fax" gorm:"type:varchar(50);comment:传真"`
	Address       string `json:"address" gorm:"type:varchar(200);comment:详细地址"`
	City          int    `json:"city" gorm:"type:int unsigned;comment:市id"`
	Province      int    `json:"province" gorm:"type:int unsigned;comment:省id"`
	Country       int    `json:"country" gorm:"type:int unsigned;comment:国家id"`
	Telephone     string `json:"telephone" gorm:"type:varchar(50);comment:电话"`
	Remark        string `json:"remark" gorm:"type:varchar(255);comment:备注"`
	Status        string `json:"status" gorm:"type:tinyint(1);comment:0无效 1有效"`
	BackupLinkman string `json:"backupLinkman" gorm:"type:varchar(50);comment:后备联系人"`
	BackupPhone   string `json:"backupPhone" gorm:"type:varchar(50);comment:后备联系电话"`
	models.ModelTime
	models.ControlBy
}

func (Vendors) TableName() string {
	return "vendors"
}

func (e *Vendors) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *Vendors) GetId() interface{} {
	return e.Id
}

func (e *Vendors) CheckNameZh(tx *gorm.DB, name string, id int) (bool, error) {
	var result struct {
		Id int
	}
	tx = tx.Model(&Vendors{}).Select("id")
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

func (e *Vendors) CheckCode(tx *gorm.DB, code string, id int) (bool, error) {
	var result struct {
		Id int
	}
	tx = tx.Model(&Vendors{}).Select("id")
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

func (e *Vendors) GetById(tx *gorm.DB, id int) error {
	return tx.Take(e, id).Error
}

func (e *Vendors) GetByNameZh(tx *gorm.DB, name string) (*Vendors, error) {
	var vendor = &Vendors{}
	err := tx.Where("name_zh = ?", name).First(&vendor).Error
	if vendor.Id != 0 {
		return vendor, nil
	}
	return vendor, err
}

func GetVendorListByIds(tx *gorm.DB, ids []int) (*[]Vendors, error) {
	var vendors = &[]Vendors{}
	err := tx.Find(vendors, ids).Error
	return vendors, err
}

func GetVendorsMapByIds(tx *gorm.DB, ids []int) map[int]string {
	vendorSlice, _ := GetVendorListByIds(tx, ids)
	vendorIdMap := lo.Associate(*vendorSlice, func(f Vendors) (int, string) {
		return f.Id, f.NameZh
	})
	return vendorIdMap
}

func (e *Vendors) FindOneByName(tx *gorm.DB, nameZh string) (*Vendors, error) {
	var data Vendors
	wcPrefix := global.GetTenantWcDBNameWithDB(tx)
	err := tx.Table(wcPrefix+"."+e.TableName()).Where("name_zh = ?", nameZh).First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}
