package dto

import (
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type WarehouseGetPageReq struct {
	dto.Pagination `search:"-"`
	WarehouseOrder
}

type WarehouseOrder struct {
	Id            string `form:"idOrder"  search:"type:order;column:id;table:warehouse"`
	WarehouseCode string `form:"warehouseCodeOrder"  search:"type:order;column:warehouse_code;table:warehouse"`
	WarehouseName string `form:"warehouseNameOrder"  search:"type:order;column:warehouse_name;table:warehouse"`
	CompanyId     string `form:"companyIdOrder"  search:"type:order;column:company_id;table:warehouse"`
	Mobile        string `form:"mobileOrder"  search:"type:order;column:mobile;table:warehouse"`
	Linkman       string `form:"linkmanOrder"  search:"type:order;column:linkman;table:warehouse"`
	Email         string `form:"emailOrder"  search:"type:order;column:email;table:warehouse"`
	Status        string `form:"statusOrder"  search:"type:order;column:status;table:warehouse"`
	IsVirtual     string `form:"isVirtualOrder"  search:"type:order;column:is_virtual;table:warehouse"`
	PostCode      string `form:"postCodeOrder"  search:"type:order;column:post_code;table:warehouse"`
	Province      string `form:"provinceOrder"  search:"type:order;column:province;table:warehouse"`
	City          string `form:"cityOrder"  search:"type:order;column:city;table:warehouse"`
	District      string `form:"districtOrder"  search:"type:order;column:district;table:warehouse"`
	Address       string `form:"addressOrder"  search:"type:order;column:address;table:warehouse"`
	DistrictName  string `form:"districtNameOrder"  search:"type:order;column:district_name;table:warehouse"`
	CityName      string `form:"cityNameOrder"  search:"type:order;column:city_name;table:warehouse"`
	ProvinceName  string `form:"provinceNameOrder"  search:"type:order;column:province_name;table:warehouse"`
	Remark        string `form:"remarkOrder"  search:"type:order;column:remark;table:warehouse"`
	CreatedAt     string `form:"createdAtOrder"  search:"type:order;column:created_at;table:warehouse"`
	UpdatedAt     string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:warehouse"`
	DeletedAt     string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:warehouse"`
	CreateBy      string `form:"createByOrder"  search:"type:order;column:create_by;table:warehouse"`
	UpdateBy      string `form:"updateByOrder"  search:"type:order;column:update_by;table:warehouse"`
	CreateByName  string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:warehouse"`
	UpdateByName  string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:warehouse"`
}

func (m *WarehouseGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type WarehouseInsertReq struct {
	Id            int    `json:"-" comment:"id"` // id
	WarehouseCode string `json:"warehouseCode" comment:"仓库编码"`
	WarehouseName string `json:"warehouseName" comment:"仓库名称"`
	CompanyId     int    `json:"companyId" comment:"仓库对应公司d"`
	Mobile        string `json:"mobile" comment:""`
	Linkman       string `json:"linkman" comment:"联系人"`
	Email         string `json:"email" comment:"邮箱"`
	Status        string `json:"status" comment:"是否使用 0-否，1-是"`
	IsVirtual     string `json:"isVirtual" comment:"是否为虚拟仓 0-否，1-是"`
	PostCode      string `json:"postCode" comment:"仓库所在地址邮编"`
	Province      int    `json:"province" comment:"省"`
	City          int    `json:"city" comment:"市"`
	District      int    `json:"district" comment:"区"`
	Address       string `json:"address" comment:"地址"`
	DistrictName  string `json:"districtName" comment:"区名"`
	CityName      string `json:"cityName" comment:"市名"`
	ProvinceName  string `json:"provinceName" comment:"省名"`
	Remark        string `json:"remark" comment:""`
	CreateByName  string `json:"createByName" comment:""`
	UpdateByName  string `json:"updateByName" comment:""`
	common.ControlBy
}

func (s *WarehouseInsertReq) Generate(model *models.Warehouse) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.WarehouseCode = s.WarehouseCode
	model.WarehouseName = s.WarehouseName
	model.CompanyId = s.CompanyId
	model.Mobile = s.Mobile
	model.Linkman = s.Linkman
	model.Email = s.Email
	model.Status = s.Status
	model.IsVirtual = s.IsVirtual
	model.PostCode = s.PostCode
	model.Province = s.Province
	model.City = s.City
	model.District = s.District
	model.Address = s.Address
	model.DistrictName = s.DistrictName
	model.CityName = s.CityName
	model.ProvinceName = s.ProvinceName
	model.Remark = s.Remark
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
	model.UpdateByName = s.UpdateByName
}

func (s *WarehouseInsertReq) GetId() interface{} {
	return s.Id
}

type WarehouseUpdateReq struct {
	Id            int    `uri:"id" comment:"id"` // id
	WarehouseCode string `json:"warehouseCode" comment:"仓库编码"`
	WarehouseName string `json:"warehouseName" comment:"仓库名称"`
	CompanyId     int    `json:"companyId" comment:"仓库对应公司d"`
	Mobile        string `json:"mobile" comment:""`
	Linkman       string `json:"linkman" comment:"联系人"`
	Email         string `json:"email" comment:"邮箱"`
	Status        string `json:"status" comment:"是否使用 0-否，1-是"`
	IsVirtual     string `json:"isVirtual" comment:"是否为虚拟仓 0-否，1-是"`
	PostCode      string `json:"postCode" comment:"仓库所在地址邮编"`
	Province      int    `json:"province" comment:"省"`
	City          int    `json:"city" comment:"市"`
	District      int    `json:"district" comment:"区"`
	Address       string `json:"address" comment:"地址"`
	DistrictName  string `json:"districtName" comment:"区名"`
	CityName      string `json:"cityName" comment:"市名"`
	ProvinceName  string `json:"provinceName" comment:"省名"`
	Remark        string `json:"remark" comment:""`
	CreateByName  string `json:"createByName" comment:""`
	UpdateByName  string `json:"updateByName" comment:""`
	common.ControlBy
}

func (s *WarehouseUpdateReq) Generate(model *models.Warehouse) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.WarehouseCode = s.WarehouseCode
	model.WarehouseName = s.WarehouseName
	model.CompanyId = s.CompanyId
	model.Mobile = s.Mobile
	model.Linkman = s.Linkman
	model.Email = s.Email
	model.Status = s.Status
	model.IsVirtual = s.IsVirtual
	model.PostCode = s.PostCode
	model.Province = s.Province
	model.City = s.City
	model.District = s.District
	model.Address = s.Address
	model.DistrictName = s.DistrictName
	model.CityName = s.CityName
	model.ProvinceName = s.ProvinceName
	model.Remark = s.Remark
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.CreateByName = s.CreateByName
	model.UpdateByName = s.UpdateByName
}

func (s *WarehouseUpdateReq) GetId() interface{} {
	return s.Id
}

// WarehouseGetReq 功能获取请求参数
type WarehouseGetReq struct {
	Id int `uri:"id"`
}

func (s *WarehouseGetReq) GetId() interface{} {
	return s.Id
}

// WarehouseDeleteReq 功能删除请求参数
type WarehouseDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *WarehouseDeleteReq) GetId() interface{} {
	return s.Ids
}
