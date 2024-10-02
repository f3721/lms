package dto

import (
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type LogicWarehouseSelectReq struct {
	dto.Pagination     `search:"-"`
	LogicWarehouseCode string `form:"logicWarehouseCode"  search:"type:exact;column:logic_warehouse_code;table:logic_warehouse" comment:"逻辑仓库编码"`
	LogicWarehouseName string `form:"logicWarehouseName"  search:"type:contains;column:logic_warehouse_name;table:logic_warehouse" comment:"逻辑仓库名称"`
	WarehouseCode      string `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:logic_warehouse" comment:"逻辑仓库对应实体仓code"`
	Type               string `form:"type"  search:"type:exact;column:type;table:logic_warehouse" comment:"0 正品仓 1次品仓"`
	//Status             string `form:"status"  search:"type:exact;column:status;table:logic_warehouse" comment:"是否使用 0-否，1-是"`
	IsTransfer string `form:"isTransfer"  search:"-" comment:"调拨单获取仓库，权限验证"`
}

func (m *LogicWarehouseSelectReq) GetNeedSearch() interface{} {
	return *m
}

type LogicWarehouseSelectResp struct {
	LogicWarehouseCode string `json:"logicWarehouseCode"`
	LogicWarehouseName string `json:"logicWarehouseName"`
}

func (s *LogicWarehouseSelectResp) ReGenerate(model *models.LogicWarehouse) {
	s.LogicWarehouseCode = model.LogicWarehouseCode
	s.LogicWarehouseName = model.LogicWarehouseName
}

type LogicWarehouseGetPageReq struct {
	dto.Pagination     `search:"-"`
	LogicWarehouseCode string `form:"logicWarehouseCode"  search:"type:exact;column:logic_warehouse_code;table:logic_warehouse" comment:"逻辑仓库编码"`
	LogicWarehouseName string `form:"logicWarehouseName"  search:"type:exact;column:logic_warehouse_name;table:logic_warehouse" comment:"逻辑仓库名称"`
	WarehouseCode      string `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:logic_warehouse" comment:"逻辑仓库对应实体仓code"`
	Mobile             string `form:"mobile"  search:"type:exact;column:mobile;table:logic_warehouse" comment:""`
	Linkman            string `form:"linkman"  search:"type:exact;column:linkman;table:logic_warehouse" comment:"联系人"`
	Email              string `form:"email"  search:"type:exact;column:email;table:logic_warehouse" comment:"邮箱"`
	Type               string `form:"type"  search:"type:exact;column:type;table:logic_warehouse" comment:"0 正品仓 1次品仓"`
	Status             string `form:"status"  search:"type:exact;column:status;table:logic_warehouse" comment:"是否使用 0-否，1-是"`
}

func (m *LogicWarehouseGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type LogicWarehouseGetPageResp struct {
	models.LogicWarehouse
	WarehouseName string `json:"warehouseName"`
	TypeName      string `json:"typeName"`
}

type LogicWarehouseInsertReq struct {
	Id                 int    `json:"-" comment:"id"` // id
	LogicWarehouseName string `json:"logicWarehouseName" comment:"逻辑仓库名称" vd:"@:len($)>0; msg:'逻辑仓库名称不能为空'"`
	WarehouseCode      string `json:"warehouseCode" comment:"逻辑仓库对应实体仓code" vd:"@:len($)>0; msg:'实体仓不能为空'"`
	Mobile             string `json:"mobile" comment:""`
	Linkman            string `json:"linkman" comment:"联系人"`
	Email              string `json:"email" comment:"邮箱"`
	Type               string `json:"type" comment:"0 正品仓 1次品仓" vd:"$=='0' || $=='1'; msg:'type为0或1'"`
	Remark             string `json:"remark" comment:""`
	common.ControlBy
}

func (s *LogicWarehouseInsertReq) Generate(model *models.LogicWarehouse) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.LogicWarehouseName = s.LogicWarehouseName
	model.WarehouseCode = s.WarehouseCode
	model.Mobile = s.Mobile
	model.Linkman = s.Linkman
	model.Email = s.Email
	model.Type = s.Type
	model.Status = models.LogicWarehouseModeStatus1
	model.CreateBy = s.CreateBy         // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName // 添加这而，需要记录是被谁创建的
	model.Remark = s.Remark
}

func (s *LogicWarehouseInsertReq) GetId() interface{} {
	return s.Id
}

type LogicWarehouseUpdateReq struct {
	Id                 int    `uri:"id" comment:"id"` // id
	LogicWarehouseName string `json:"logicWarehouseName" comment:"逻辑仓库名称" vd:"@:len($)>0; msg:'逻辑仓库名称不能为空'"`
	Mobile             string `json:"mobile" comment:""`
	Linkman            string `json:"linkman" comment:"联系人"`
	Email              string `json:"email" comment:"邮箱"`
	Remark             string `json:"remark" comment:""`
	common.ControlBy
}

func (s *LogicWarehouseUpdateReq) Generate(model *models.LogicWarehouse) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.LogicWarehouseName = s.LogicWarehouseName
	model.Mobile = s.Mobile
	model.Linkman = s.Linkman
	model.Email = s.Email
	model.UpdateBy = s.UpdateBy         // 添加这而，需要记录是被谁更新的
	model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
	model.Remark = s.Remark
}

func (s *LogicWarehouseUpdateReq) GetId() interface{} {
	return s.Id
}

// LogicWarehouseGetReq 功能获取请求参数
type LogicWarehouseGetReq struct {
	Id int `uri:"id"`
}

func (s *LogicWarehouseGetReq) GetId() interface{} {
	return s.Id
}

type LogicWarehouseGetResp struct {
	models.LogicWarehouse
	WarehouseName string `json:"warehouseName"`
	TypeName      string `json:"typeName"`
}

// LogicWarehouseDeleteReq 功能删除请求参数
type LogicWarehouseDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *LogicWarehouseDeleteReq) GetId() interface{} {
	return s.Ids
}
