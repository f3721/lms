package dto

import (
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"gorm.io/gorm"
)

type WarehouseSelectReq struct {
	dto.Pagination `search:"-"`
	WarehouseCode  string `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:warehouse" comment:"仓库编码"`
	WarehouseName  string `form:"warehouseName"  search:"type:contains;column:warehouse_name;table:warehouse" comment:"仓库名称"`
	CompanyId      int    `form:"companyId"  search:"type:exact;column:company_id;table:warehouse" comment:"仓库对应公司d"`
	//Status         string `form:"status"  search:"type:exact;column:status;table:warehouse" comment:"是否使用 0-否，1-是"`
	IsVirtual  string `form:"isVirtual"  search:"type:exact;column:is_virtual;table:warehouse" comment:"是否为虚拟仓 0-否，1-是"`
	IsTransfer string `form:"isTransfer"  search:"-" comment:"调拨单获取仓库，权限验证"`
}

func (m *WarehouseSelectReq) GetNeedSearch() interface{} {
	return *m
}

type WarehouseSelectResp struct {
	WarehouseCode string `json:"warehouseCode"`
	WarehouseName string `json:"warehouseName"`
}

func (s *WarehouseSelectResp) ReGenerate(model *models.Warehouse) {
	s.WarehouseCode = model.WarehouseCode
	s.WarehouseName = model.WarehouseName
}

type WarehouseGetPageReq struct {
	dto.Pagination `search:"-"`
	WarehouseCode  string `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:warehouse" comment:"仓库编码"`
	WarehouseName  string `form:"warehouseName"  search:"type:exact;column:warehouse_name;table:warehouse" comment:"仓库名称"`
	CompanyId      int    `form:"companyId"  search:"type:exact;column:company_id;table:warehouse" comment:"仓库对应公司d"`
	Mobile         string `form:"mobile"  search:"type:exact;column:mobile;table:warehouse" comment:""`
	Linkman        string `form:"linkman"  search:"type:exact;column:linkman;table:warehouse" comment:"联系人"`
	Email          string `form:"email"  search:"type:exact;column:email;table:warehouse" comment:"邮箱"`
	Status         string `form:"status"  search:"type:exact;column:status;table:warehouse" comment:"是否使用 0-否，1-是"`
	IsVirtual      string `form:"isVirtual"  search:"type:exact;column:is_virtual;table:warehouse" comment:"是否为虚拟仓 0-否，1-是"`
	CreatedBy      int    `form:"createdBy"  search:"type:exact;column:created_by;table:warehouse" comment:"创建人id"`
	UpdatedBy      int    `form:"updatedBy"  search:"type:exact;column:updated_by;table:warehouse" comment:"修改人id"`
	PostCode       string `form:"postCode"  search:"type:exact;column:post_code;table:warehouse" comment:"仓库所在地址邮编"`
	Province       int    `form:"province"  search:"type:exact;column:province;table:warehouse" comment:"省"`
	City           int    `form:"city"  search:"type:exact;column:city;table:warehouse" comment:"市"`
	District       int    `form:"district"  search:"type:exact;column:district;table:warehouse" comment:"区"`
}

func (m *WarehouseGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type WarehouseGetPageResp struct {
	models.Warehouse
	CompanyName   string `json:"companyName"`
	IsVirtualName string `json:"isVirtualName"`
}

type WarehouseInsertReq struct {
	Id            int    `json:"-" comment:"id"` // id
	WarehouseName string `json:"warehouseName" comment:"仓库名称" vd:"@:len($)>0; msg:'仓库名称不能为空'"`
	CompanyId     int    `json:"companyId" comment:"仓库对应公司d" vd:"$>0; msg:'companyId不能为空'"`
	Mobile        string `json:"mobile" comment:"" vd:"regexp('^1[0-9]{10}$'); msg:'手机格号码式不正确'"`
	Linkman       string `json:"linkman" comment:"联系人" vd:"@:len($)>0; msg:'联系人不能为空'"`
	Email         string `json:"email" comment:"邮箱"`
	//Status        string `json:"status" comment:"是否使用 0-否，1-是" vd:"$=='0' || $=='1'; msg:'status为0或1'"`
	IsVirtual string `json:"isVirtual" comment:"是否为虚拟仓 0-否，1-是" vd:"$=='0' || $=='1'; msg:'isVirtual为0或1'"`
	PostCode  string `json:"postCode" comment:"仓库所在地址邮编"`
	Province  int    `json:"province" comment:"省" vd:"$>0; msg:'省不能为空'"`
	City      int    `json:"city" comment:"市" vd:"$>0; msg:'市不能为空'"`
	District  int    `json:"district" comment:"区" vd:"$>0; msg:'区不能为空'"`
	Address   string `json:"address" comment:"地址" vd:"@:len($)>0; msg:'地址不能为空'"`
	Remark    string `json:"remark" comment:""`
	common.ControlBy
}

func (s *WarehouseInsertReq) Generate(tx *gorm.DB, model *models.Warehouse) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.WarehouseName = s.WarehouseName
	model.CompanyId = s.CompanyId
	model.Mobile = s.Mobile
	model.Linkman = s.Linkman
	model.Email = s.Email
	model.Status = models.WarehouseModeStatus1
	model.IsVirtual = s.IsVirtual
	model.PostCode = s.PostCode
	model.Province = s.Province
	model.City = s.City
	model.District = s.District
	model.Address = s.Address
	model.Remark = s.Remark
	model.ControlBy = s.ControlBy

	regionsMap := models.GeRegionMapByIds(tx, []int{s.Province, s.City, s.District})
	model.GenerateRegionName(s.Province, s.City, s.District, regionsMap)
}

func (s *WarehouseInsertReq) GetId() interface{} {
	return s.Id
}

type WarehouseUpdateReq struct {
	Id            int    `uri:"id" comment:"id"` // id
	WarehouseName string `json:"warehouseName" comment:"仓库名称" vd:"@:len($)>0; msg:'仓库名称不能为空'"`
	Mobile        string `json:"mobile" comment:"" vd:"regexp('^1[0-9]{10}$'); msg:'手机格号码式不正确'"`
	Linkman       string `json:"linkman" comment:"联系人" vd:"@:len($)>0; msg:'联系人不能为空'"`
	Email         string `json:"email" comment:"邮箱"`
	//Status        string `json:"status" comment:"是否使用 0-否，1-是" vd:"$=='0' || $=='1'; msg:'status为0或1'"`
	IsVirtual string `json:"isVirtual" comment:"是否为虚拟仓 0-否，1-是" vd:"$=='0' || $=='1'; msg:'isVirtual为0或1'"`
	PostCode  string `json:"postCode" comment:"仓库所在地址邮编"`
	Province  int    `json:"province" comment:"省" vd:"$>0; msg:'省不能为空'"`
	City      int    `json:"city" comment:"市" vd:"$>0; msg:'市不能为空'"`
	District  int    `json:"district" comment:"区" vd:"$>0; msg:'区不能为空'"`
	Address   string `json:"address" comment:"地址" vd:"@:len($)>0; msg:'地址不能为空'"`
	Remark    string `json:"remark" comment:""`
	common.ControlBy
}

func (s *WarehouseUpdateReq) Generate(tx *gorm.DB, model *models.Warehouse) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.WarehouseName = s.WarehouseName
	model.Mobile = s.Mobile
	model.Linkman = s.Linkman
	model.Email = s.Email
	model.IsVirtual = s.IsVirtual
	model.PostCode = s.PostCode
	model.Province = s.Province
	model.City = s.City
	model.District = s.District
	model.Address = s.Address
	model.Remark = s.Remark
	model.ControlBy.UpdateBy = s.ControlBy.UpdateBy
	model.ControlBy.UpdateByName = s.ControlBy.UpdateByName

	regionsMap := models.GeRegionMapByIds(tx, []int{s.Province, s.City, s.District})
	model.GenerateRegionName(s.Province, s.City, s.District, regionsMap)
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

type WarehouseGetResp struct {
	models.Warehouse
	CompanyName string `json:"companyName"`
}

// WarehouseDeleteReq 功能删除请求参数
type WarehouseDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *WarehouseDeleteReq) GetId() interface{} {
	return s.Ids
}

type WarehouseGetByCodeReq struct {
	WarehouseCode string `uri:"warehouseCode"`
}

func (s *WarehouseGetByCodeReq) GetCode() interface{} {
	return s.WarehouseCode
}

type InnerWarehouseGetListByNameAndCompanyIdReq struct {
	Query []InnerWarehouseGetListByNameAndCompanyIdReqInfo `json:"query"`
}

type InnerWarehouseGetListByNameAndCompanyIdReqInfo struct {
	WarehouseName string `json:"warehouseName"  comment:"仓库名称"`
	CompanyId     int    `json:"companyId" comment:"仓库对应公司Id"`
}

type InnerWarehouseGetListReq struct {
	WarehouseCode string `form:"warehouseCode"  search:"-" comment:"仓库编码"`
	WarehouseName string `form:"warehouseName"  search:"-" comment:"仓库名称"`
	CompanyId     int    `form:"companyId"  search:"type:exact;column:company_id;table:warehouse" comment:"仓库对应公司d"`
	Status        string `form:"status"  search:"type:exact;column:status;table:warehouse" comment:"是否使用 0-否，1-是"`
	IsVirtual     string `form:"isVirtual"  search:"type:exact;column:is_virtual;table:warehouse" comment:"是否为虚拟仓 0-否，1-是"`
}

func (m *InnerWarehouseGetListReq) GetNeedSearch() interface{} {
	return *m
}

type CompanyWarehouseTreeResp struct {
	Value    int                            `json:"value"`
	Label    string                         `json:"label"`
	Children []CompanyWarehouseTreeChildren `json:"children"`
}

type CompanyWarehouseTreeChildren struct {
	Value string `json:"value"`
	Label string `json:"label"`
}
