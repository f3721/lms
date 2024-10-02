package dto

import (
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type ReminderRuleGetPageReq struct {
	dto.Pagination `search:"-"`
	CompanyId      int    `form:"companyId"  search:"type:exact;column:company_id;table:reminder_rule" comment:"公司id"`
	WarehouseCode  string `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:reminder_rule" comment:"仓库id"`
	Status         string `form:"status"  search:"type:exact;column:status;table:reminder_rule" comment:"状态 1启用 0未启用 全部不传"`
	ReminderRuleOrder
}

type ReminderRuleOrder struct {
	Id                 string `form:"idOrder"  search:"type:order;column:id;table:reminder_rule"`
	CompanyId          string `form:"companyIdOrder"  search:"type:order;column:company_id;table:reminder_rule"`
	WarehouseCode      string `form:"warehouseCodeOrder"  search:"type:order;column:warehouse_id;table:reminder_rule"`
	WarningValue       string `form:"warningValueOrder"  search:"type:order;column:warning_value;table:reminder_rule"`
	ReplenishmentValue string `form:"replenishmentValueOrder"  search:"type:order;column:replenishment_value;table:reminder_rule"`
	Status             string `form:"statusOrder"  search:"type:order;column:status;table:reminder_rule"`
	CreatedAt          string `form:"createdAtOrder"  search:"type:order;column:created_at;table:reminder_rule"`
	UpdatedAt          string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:reminder_rule"`
	CreateBy           string `form:"createByOrder"  search:"type:order;column:create_by;table:reminder_rule"`
	UpdateBy           string `form:"updateByOrder"  search:"type:order;column:update_by;table:reminder_rule"`
	CreateByName       string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:reminder_rule"`
	UpdateByName       string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:reminder_rule"`
	DeletedAt          string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:reminder_rule"`
}

func (m *ReminderRuleGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ReminderRuleData struct {
	Id                 int    `json:"id" comment:"id"`                               // id
	CompanyId          int    `json:"companyId" comment:"公司id"`                      // 公司id
	WarehouseCode      string `json:"warehouseCode" comment:"仓库code" vd:"len($)>0"`  // 仓库code
	WarningValue       int    `json:"warningValue" comment:"SKU通用预警值" vd:"$>0"`      // SKU通用预警值
	ReplenishmentValue int    `json:"replenishmentValue" comment:"设置的备货量" vd:"$>=0"` // 设置的备货量
	Status             int    `json:"status" comment:"状态 1启用 0未启用"`                  // 状态 1启用 0未启用

	CompanyName   string `json:"companyName" gorm:"-"`
	WarehouseName string `json:"warehouseName" gorm:"-"`
	common.ControlBy
	common.ModelTime
}

type ReminderRuleInsertReq struct {
	Id                 int    `json:"-" comment:"id"`                                             // id
	CompanyId          int    `json:"companyId" comment:"公司id"`                                   // 公司id
	WarehouseCode      string `json:"warehouseCode" comment:"仓库code" vd:"len($)>0"`               // 仓库code
	WarningValue       int    `json:"warningValue" comment:"SKU通用预警值" vd:"$>=0;msg:'预警值必须大于等于0'"` // SKU通用预警值
	ReplenishmentValue int    `json:"replenishmentValue" comment:"设置的备货量" vd:"$>=0"`              // 设置的备货量
	Status             int    `json:"status" comment:"状态 1启用 0未启用"`                               // 状态 1启用 0未启用
	common.ControlBy   `json:"-"`
	SkuList            []*ReminderRuleSkuInsertReq `json:"skuList"`
}

func (s *ReminderRuleInsertReq) Generate(model *models.ReminderRule) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CompanyId = s.CompanyId
	model.WarehouseCode = s.WarehouseCode
	model.WarningValue = s.WarningValue
	model.ReplenishmentValue = s.ReplenishmentValue
	model.Status = s.Status
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *ReminderRuleInsertReq) GetId() interface{} {
	return s.Id
}

type ReminderRuleUpdateReq struct {
	Id                 int    `uri:"id" comment:"id"  vd:"$>0"` // id
	CompanyId          int    `json:"companyId" comment:"公司id" vd:"$>0"`
	WarehouseCode      string `json:"warehouseCode" comment:"仓库id" vd:"len($)>0"`
	WarningValue       int    `json:"warningValue" comment:"SKU通用预警值" vd:"$>=0"`
	ReplenishmentValue int    `json:"replenishmentValue" comment:"设置的备货量" vd:"$>=0"`
	Status             int    `json:"status" comment:"状态 1启用 0未启用"`
	common.ControlBy   `json:"-"`
	SkuList            []*ReminderRuleSkuUpdateReq `json:"skuList"`
}

func (s *ReminderRuleUpdateReq) Generate(model *models.ReminderRule) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CompanyId = s.CompanyId
	model.WarehouseCode = s.WarehouseCode
	model.WarningValue = s.WarningValue
	model.ReplenishmentValue = s.ReplenishmentValue
	model.Status = s.Status
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName
}

func (s *ReminderRuleUpdateReq) GetId() interface{} {
	return s.Id
}

// ReminderRuleGetReq 功能获取请求参数
type ReminderRuleGetReq struct {
	Id int `uri:"id"`
}

type ReminderRuleGetRes struct {
	ReminderRuleData
}

func (s *ReminderRuleGetReq) GetId() interface{} {
	return s.Id
}

// ReminderRuleDeleteReq 功能删除请求参数
type ReminderRuleDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ReminderRuleDeleteReq) GetId() interface{} {
	return s.Ids
}
