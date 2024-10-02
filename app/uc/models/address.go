package models

import (
	"go-admin/common/global"
	"go-admin/common/models"

	"gorm.io/gorm"
)

const AddressModelName = "userInfo"
const AddressOperationCreate = "createAddress"
const AddressOperationUpdate = "updateAddress"

var AddressType = map[int]string{
	1: "收货地址",
	2: "发票地址",
	3: "产品使用地址",
	4: "通用",
}

type Address struct {
	models.Model

	UserId                  int    `json:"userId" gorm:"type:int unsigned;comment:用户ID"`             //用户ID
	ReceiverName            string `json:"receiverName" gorm:"type:varchar(64);comment:收货姓名"`        //收货姓名
	DetailAddress           string `json:"detailAddress" gorm:"type:varchar(256);comment:详细地址"`      //详细地址
	PostalCode              string `json:"postalCode" gorm:"type:varchar(10);comment:邮政编码"`          //邮政编码
	CellPhone               string `json:"cellPhone" gorm:"type:varchar(20);comment:收件人手机号码"`        //收件人手机号码
	Telephone               string `json:"telephone" gorm:"type:varchar(64);comment:座机号码"`           //座机号码
	IsDefault               int    `json:"isDefault" gorm:"type:tinyint unsigned;comment:是否默认(1默认)"` //是否默认(1默认)
	CountryId               int    `json:"countryId" gorm:"type:int unsigned;comment:国家ID"`          //国家ID
	CountryName             string `json:"countryName" gorm:"type:varchar(64);comment:国家名称"`         //国家名称
	ProvinceName            string `json:"provinceName" gorm:"type:varchar(64);comment:省份名称"`        //省份名称
	ProvinceId              int    `json:"provinceId" gorm:"type:int unsigned;comment:省份ID"`         //省份ID
	CityId                  int    `json:"cityId" gorm:"type:int unsigned;comment:城市ID"`             //城市ID
	CityName                string `json:"cityName" gorm:"type:varchar(64);comment:城市名称"`            //城市名称
	AreaId                  int    `json:"areaId" gorm:"type:int unsigned;comment:区域ID"`             //区域ID
	AreaName                string `json:"areaName" gorm:"type:varchar(64);comment:区域名称"`            //区域名称
	TownId                  int    `json:"townId" gorm:"type:int unsigned;comment:镇ID"`              //镇(街道)ID
	TownName                string `json:"townName" gorm:"type:varchar(50);comment:镇名"`              //镇(街道)ID
	CompanyId               int    `json:"companyId" gorm:"type:int unsigned;comment:公司ID"`
	CompanyName             string `json:"companyName" gorm:"type:varchar(50);comment:公司名称"`
	AddressType             int    `json:"addressType" gorm:"type:tinyint unsigned;comment:地址类型（1收货地址，2发票地址,3产品使用地址，4通用）"` //地址类型（1收货地址)
	IsDeadlineReceivingAddr int    `json:"isDeadlineReceivingAddr" gorm:"type:tinyint unsigned;comment:是否账期收货地址"`
	Email                   string `json:"email" gorm:"type:varchar(64);comment:邮箱（用于发送出库提醒）"` //邮箱（用于发送出库提醒）
	Department              string `json:"department" gorm:"type:varchar(50);comment:部门"`
	PcCode                  string `json:"pcCode" gorm:"type:varchar(256);comment:仓库编号"` //仓库编号
	models.ModelTime
	models.ControlBy
}

func (Address) TableName() string {
	return "address"
}

func (e *Address) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *Address) GetId() interface{} {
	return e.Id
}

func (e *Address) Get(tx *gorm.DB, id int) (err error) {
	ucPrefix := global.GetTenantUcDBNameWithDB(tx)
	err = tx.Table(ucPrefix+"."+e.TableName()).First(e, id).Error
	return
}

func (e *Address) IsAddressTypeExists(db *gorm.DB, addressType int64, userId int) bool {
	var count int64
	db.Model(Address{}).
		Where("address_type = ?", addressType).
		Where("user_id = ?", userId).
		Limit(1).Count(&count)
	return count > 0
}
