package dto

import (
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SupplierSelectReq struct {
	dto.Pagination `search:"-"`
	Id             int    `form:"id"  search:"type:exact;column:id;table:supplier" comment:"供应商id"`
	Code           string `form:"code"  search:"type:exact;column:code;table:supplier" comment:"供应商编码"`
	NameZh         string `form:"nameZh"  search:"type:contains;column:name_zh;table:supplier" comment:"供应商中文名"`
	ShortName      string `form:"shortName" search:"type:contains;column:short_name;table:supplier"`
	Status         string `form:"status"  search:"type:exact;column:status;table:supplier" comment:"是否有效 0-否，1-是"`
}

func (m *SupplierSelectReq) GetNeedSearch() interface{} {
	return *m
}

type SupplierSelectResp struct {
	Id        int    `json:"id"`
	NameZh    string `json:"nameZh"`
	ShortName string `json:"shortName"`
}

func (s *SupplierSelectResp) ReGenerate(model *models.Supplier) {
	s.Id = model.Id
	s.NameZh = model.NameZh
	s.ShortName = model.ShortName
}

type SupplierGetPageReq struct {
	dto.Pagination `search:"-"`
	Code           string `form:"code"  search:"type:exact;column:code;table:supplier" comment:"供应商编码"`
	NameZh         string `form:"nameZh"  search:"type:contains;column:name_zh;table:supplier" comment:"供应商中文名"`
	NameEn         string `form:"nameEn"  search:"type:contains;column:name_en;table:supplier" comment:"供应商英文名"`
	ShortName      string `form:"shortName"  search:"type:contains;column:short_name;table:supplier" comment:"供应商简称"`
	CityId         string `form:"cityId"  search:"type:exact;column:city;table:supplier" comment:"市id"`
	ProvinceId     string `form:"provinceId"  search:"type:exact;column:province;table:supplier" comment:"省id"`
	CountryId      string `form:"countryId"  search:"type:exact;column:country;table:supplier" comment:"国家id"`
	Status         string `form:"status"  search:"type:exact;column:status;table:supplier" comment:"0无效 1有效"`
}

func (m *SupplierGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type SupplierInsertReq struct {
	Id            int    `json:"-" comment:""` //
	Code          string `json:"code" comment:"供应商编码" vd:"@:len($)>0&&mblen($)<11; msg:'供应商编码长度在0-10之间'"`
	NameZh        string `json:"nameZh" comment:"供应商中文名" vd:"@:len($)>0&&mblen($)<51; msg:'供应商中文名长度在0-50之间'"`
	NameEn        string `json:"nameEn" comment:"供应商英文名"`
	ShortName     string `json:"shortName" comment:"供应商简称" vd:"@:len($)>0&&mblen($)<21; msg:'货主编码长度在0-20之间'"`
	PostalCode    string `json:"postalCode" comment:"邮编"`
	Linkman       string `json:"linkman" comment:"联系人" vd:"@:len($)>0; msg:'主联系人必填'"`
	Phone         string `json:"phone" comment:"手机" vd:"regexp('^1[0-9]{10}$'); msg:'主联系人手机格式不正确'"`
	Email         string `json:"email" comment:"邮箱"`
	Fax           string `json:"fax" comment:"传真"`
	Address       string `json:"address" comment:"详细地址"`
	CityId        int    `json:"cityId" comment:"市id" `
	ProvinceId    int    `json:"provinceId" comment:"省id"`
	CountryId     int    `json:"countryId" comment:"国家id" vd:"$>0; msg:'国家不能为空'"`
	Telephone     string `json:"telephone" comment:"电话"`
	CompanyId     int    `json:"companyId" comment:"关联公司"`
	Remark        string `json:"remark" comment:"备注"`
	Status        string `json:"status" comment:"0无效 1有效" vd:"$=='0' || $=='1'; msg:'status为0或1'"`
	BackupLinkman string `json:"backupLinkman" comment:"后备联系人"`
	BackupPhone   string `json:"backupPhone" comment:"后备联系电话"`
	common.ControlBy
}

func (s *SupplierInsertReq) Generate(model *models.Supplier) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.Code = s.Code
	model.NameZh = s.NameZh
	model.NameEn = s.NameEn
	model.ShortName = s.ShortName
	model.PostalCode = s.PostalCode
	model.Linkman = s.Linkman
	model.Phone = s.Phone
	model.Email = s.Email
	model.Fax = s.Fax
	model.Address = s.Address
	model.CityId = s.CityId
	model.ProvinceId = s.ProvinceId
	model.CountryId = s.CountryId
	model.Telephone = s.Telephone
	model.Remark = s.Remark
	model.Status = s.Status
	model.BackupLinkman = s.BackupLinkman
	model.BackupPhone = s.BackupPhone
	model.CreateBy = s.CreateBy         // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName // 添加这而，需要记录是被谁创建的
}

func (s *SupplierInsertReq) GetId() interface{} {
	return s.Id
}

type SupplierUpdateReq struct {
	Id            int    `uri:"id" comment:""` //
	NameZh        string `json:"nameZh" comment:"供应商中文名" vd:"@:len($)>0&&mblen($)<51; msg:'供应商中文名长度在0-50之间'"`
	NameEn        string `json:"nameEn" comment:"供应商英文名"`
	ShortName     string `json:"shortName" comment:"供应商简称" vd:"@:len($)>0&&mblen($)<21; msg:'供应商编码长度在0-20之间'"`
	PostalCode    string `json:"postalCode" comment:"邮编"`
	Linkman       string `json:"linkman" comment:"联系人" vd:"@:len($)>0; msg:'主联系人必填'"`
	Phone         string `json:"phone" comment:"手机" vd:"regexp('^1[0-9]{10}$'); msg:'主联系人手机格式不正确'"`
	Email         string `json:"email" comment:"邮箱"`
	Fax           string `json:"fax" comment:"传真"`
	Address       string `json:"address" comment:"详细地址"`
	CityId        int    `json:"cityId" comment:"市id"`
	ProvinceId    int    `json:"provinceId" comment:"省id"`
	CountryId     int    `json:"countryId" comment:"国家id" vd:"$>0; msg:'国家不能为空'"`
	Telephone     string `json:"telephone" comment:"电话"`
	CompanyId     int    `json:"companyId" comment:"关联公司"`
	Remark        string `json:"remark" comment:"备注"`
	Status        string `json:"status" comment:"0无效 1有效" vd:"$=='0' || $=='1'; msg:'status为0或1'"`
	BackupLinkman string `json:"backupLinkman" comment:"后备联系人"`
	BackupPhone   string `json:"backupPhone" comment:"后备联系电话"`
	common.ControlBy
}

func (s *SupplierUpdateReq) Generate(model *models.Supplier) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.NameZh = s.NameZh
	model.NameEn = s.NameEn
	model.ShortName = s.ShortName
	model.PostalCode = s.PostalCode
	model.Linkman = s.Linkman
	model.Phone = s.Phone
	model.Email = s.Email
	model.Fax = s.Fax
	model.Address = s.Address
	model.CityId = s.CityId
	model.ProvinceId = s.ProvinceId
	model.CountryId = s.CountryId
	model.Telephone = s.Telephone
	model.Remark = s.Remark
	model.Status = s.Status
	model.BackupLinkman = s.BackupLinkman
	model.BackupPhone = s.BackupPhone
	model.UpdateBy = s.UpdateBy         // 添加这而，需要记录是被谁更新的
	model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的

}

func (s *SupplierUpdateReq) GetId() interface{} {
	return s.Id
}

// SupplierGetReq 功能获取请求参数
type SupplierGetReq struct {
	Id int `uri:"id"`
}

func (s *SupplierGetReq) GetId() interface{} {
	return s.Id
}

// supplierDeleteReq 功能删除请求参数
type SupplierDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *SupplierDeleteReq) GetId() interface{} {
	return s.Ids
}

type InnerSupplierGetListReq struct {
	Ids    string `form:"ids"  search:"-" comment:"供应商编码"`
	Code   string `form:"code"  search:"-" comment:"供应商编码"`
	NameZh string `form:"nameZh"  search:"-" comment:"供应商中文名"`
	Status string `form:"status"  search:"type:exact;column:status;table:supplier" comment:"是否有效 0-否，1-是"`
}

func (m *InnerSupplierGetListReq) GetNeedSearch() interface{} {
	return *m
}
