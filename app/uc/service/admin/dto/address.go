package dto

import (
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type AddressGetPageReq struct {
	dto.Pagination `search:"-"`
	UserId         int `form:"userId"  search:"type:exact;column:user_id;table:address" comment:"用户ID"`
	AddressType    int `form:"addressType" search:"type:exact;column:address_type;table:address"`
	AddressOrder
}

type AddressOrder struct {
	Id                      string `form:"idOrder"  search:"type:order;column:id;table:address"`
	UserId                  string `form:"userIdOrder"  search:"type:order;column:user_id;table:address"`
	ReceiverName            string `form:"receiverNameOrder"  search:"type:order;column:receiver_name;table:address"`
	DetailAddress           string `form:"detailAddressOrder"  search:"type:order;column:detail_address;table:address"`
	PostalCode              string `form:"postalCodeOrder"  search:"type:order;column:postal_code;table:address"`
	CellPhone               string `form:"cellPhoneOrder"  search:"type:order;column:cell_phone;table:address"`
	Telephone               string `form:"telephoneOrder"  search:"type:order;column:telephone;table:address"`
	IsDefault               string `form:"isDefaultOrder"  search:"type:order;column:is_default;table:address"`
	CountryId               string `form:"countryIdOrder"  search:"type:order;column:country_id;table:address"`
	CountryName             string `form:"countryNameOrder"  search:"type:order;column:country_name;table:address"`
	ProvinceName            string `form:"provinceNameOrder"  search:"type:order;column:province_name;table:address"`
	ProvinceId              string `form:"provinceIdOrder"  search:"type:order;column:province_id;table:address"`
	CityId                  string `form:"cityIdOrder"  search:"type:order;column:city_id;table:address"`
	CityName                string `form:"cityNameOrder"  search:"type:order;column:city_name;table:address"`
	AreaId                  string `form:"areaIdOrder"  search:"type:order;column:area_id;table:address"`
	AreaName                string `form:"areaNameOrder"  search:"type:order;column:area_name;table:address"`
	TownId                  string `form:"townIdOrder"  search:"type:order;column:town_id;table:address"`
	TownName                string `form:"townNameOrder"  search:"type:order;column:town_name;table:address"`
	CompanyId               string `form:"companyIdOrder"  search:"type:order;column:company_id;table:address"`
	CompanyName             string `form:"companyNameOrder"  search:"type:order;column:company_name;table:address"`
	AddressType             string `form:"addressTypeOrder"  search:"type:order;column:address_type;table:address"`
	IsDeadlineReceivingAddr string `form:"isDeadlineReceivingAddrOrder"  search:"type:order;column:is_deadline_receiving_addr;table:address"`
	Email                   string `form:"emailOrder"  search:"type:order;column:email;table:address"`
	Department              string `form:"departmentOrder"  search:"type:order;column:department;table:address"`
	DeletedAt               string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:address"`
	CreatedAt               string `form:"createdAtOrder"  search:"type:order;column:created_at;table:address"`
	CreateBy                string `form:"createByOrder"  search:"type:order;column:create_by;table:address"`
	UpdateBy                string `form:"updateByOrder"  search:"type:order;column:update_by;table:address"`
	UpdatedAt               string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:address"`
	PcCode                  string `form:"pcCodeOrder"  search:"type:order;column:pc_code;table:address"`
}

func (m *AddressGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type AddressGetPageRes struct {
	models.Address
	AddressTypeText string `json:"addressTypeText"`
}

type AddressInsertReq struct {
	Id            int    `json:"-" comment:"编号"` // 编号
	UserId        int    `json:"userId" comment:"用户ID" vd:"$>0" `
	ReceiverName  string `json:"receiverName" comment:"收货姓名"`
	DetailAddress string `json:"detailAddress" comment:"详细地址"`
	PostalCode    string `json:"postalCode" comment:"邮政编码"`
	CellPhone     string `json:"cellPhone" comment:"收件人手机号码"`
	Telephone     string `json:"telephone" comment:"座机号码"`
	IsDefault     int    `json:"isDefault" comment:"是否默认(1默认)"`
	CountryId     int    `json:"countryId" comment:"国家ID" vd:"$>0" `
	CountryName   string `json:"countryName" comment:"国家名称"`
	ProvinceName  string `json:"provinceName" comment:"省份名称"`
	ProvinceId    int    `json:"provinceId" comment:"省份ID"`
	CityId        int    `json:"cityId" comment:"城市ID"`
	CityName      string `json:"cityName" comment:"城市名称"`
	AreaId        int    `json:"areaId" comment:"区域ID"`
	AreaName      string `json:"areaName" comment:"区域名称"`
	TownId        int    `json:"townId" comment:"镇ID"`
	TownName      string `json:"townName" comment:"镇名"`
	//CompanyId int `json:"companyId" comment:"公司ID"`
	//CompanyName string `json:"companyName" comment:"公司名称"`
	AddressType int `json:"addressType" comment:"地址类型（1收货地址，2发票地址,3产品使用地址，4通用）"`
	//IsDeadlineReceivingAddr int    `json:"isDeadlineReceivingAddr" comment:"是否账期收货地址"`
	Email string `json:"email" comment:"邮箱（用于发送出库提醒）"`
	//Department string `json:"department" comment:"部门"`
	//PcCode string `json:"pcCode" comment:"仓库编号"`
	common.ControlBy `json:"-"`
}

func (s *AddressInsertReq) Generate(model *models.Address) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UserId = s.UserId
	model.ReceiverName = s.ReceiverName
	model.DetailAddress = s.DetailAddress
	model.PostalCode = s.PostalCode
	model.CellPhone = s.CellPhone
	model.Telephone = s.Telephone
	model.IsDefault = s.IsDefault
	model.CountryId = s.CountryId
	model.CountryName = s.CountryName
	model.ProvinceName = s.ProvinceName
	model.ProvinceId = s.ProvinceId
	model.CityId = s.CityId
	model.CityName = s.CityName
	model.AreaId = s.AreaId
	model.AreaName = s.AreaName
	model.TownId = s.TownId
	model.TownName = s.TownName
	model.AddressType = s.AddressType
	model.Email = s.Email
	model.CreateBy = s.CreateBy         // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName // 添加这而，需要记录是被谁创建的
}

func (s *AddressInsertReq) GetId() interface{} {
	return s.Id
}

type AddressUpdateReq struct {
	Id            int    `uri:"id" comment:"编号"` // 编号
	UserId        int    `json:"userId" comment:"用户ID"`
	ReceiverName  string `json:"receiverName" comment:"收货姓名"`
	DetailAddress string `json:"detailAddress" comment:"详细地址"`
	PostalCode    string `json:"postalCode" comment:"邮政编码"`
	CellPhone     string `json:"cellPhone" comment:"收件人手机号码"`
	Telephone     string `json:"telephone" comment:"座机号码"`
	IsDefault     int    `json:"isDefault" comment:"是否默认(1默认)"`
	CountryId     int    `json:"countryId" comment:"国家ID"`
	CountryName   string `json:"countryName" comment:"国家名称"`
	ProvinceName  string `json:"provinceName" comment:"省份名称"`
	ProvinceId    int    `json:"provinceId" comment:"省份ID"`
	CityId        int    `json:"cityId" comment:"城市ID"`
	CityName      string `json:"cityName" comment:"城市名称"`
	AreaId        int    `json:"areaId" comment:"区域ID"`
	AreaName      string `json:"areaName" comment:"区域名称"`
	TownId        int    `json:"townId" comment:"镇ID"`
	TownName      string `json:"townName" comment:"镇名"`
	//CompanyId     int    `json:"companyId" comment:"公司ID"`
	//CompanyName   string `json:"companyName" comment:"公司名称"`
	AddressType int `json:"addressType" comment:"地址类型（1收货地址，2发票地址,3产品使用地址，4通用）"`
	//IsDeadlineReceivingAddr int    `json:"isDeadlineReceivingAddr" comment:"是否账期收货地址"`
	Email            string `json:"email" comment:"邮箱（用于发送出库提醒）"`
	Department       string `json:"department" comment:"部门"`
	PcCode           string `json:"pcCode" comment:"仓库编号"`
	common.ControlBy `json:"-"`
}

func (s *AddressUpdateReq) Generate(model *models.Address) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UserId = s.UserId
	model.ReceiverName = s.ReceiverName
	model.DetailAddress = s.DetailAddress
	model.PostalCode = s.PostalCode
	model.CellPhone = s.CellPhone
	model.Telephone = s.Telephone
	model.IsDefault = s.IsDefault
	model.CountryId = s.CountryId
	model.CountryName = s.CountryName
	model.ProvinceName = s.ProvinceName
	model.ProvinceId = s.ProvinceId
	model.CityId = s.CityId
	model.CityName = s.CityName
	model.AreaId = s.AreaId
	model.AreaName = s.AreaName
	model.TownId = s.TownId
	model.TownName = s.TownName
	model.AddressType = s.AddressType
	model.Email = s.Email
	model.Department = s.Department
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName
	model.PcCode = s.PcCode
}

func (s *AddressUpdateReq) GetId() interface{} {
	return s.Id
}

// AddressGetReq 功能获取请求参数
type AddressGetReq struct {
	Id int `uri:"id"`
}

func (s *AddressGetReq) GetId() interface{} {
	return s.Id
}

// AddressDeleteReq 功能删除请求参数
type AddressDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *AddressDeleteReq) GetId() interface{} {
	return s.Ids
}
