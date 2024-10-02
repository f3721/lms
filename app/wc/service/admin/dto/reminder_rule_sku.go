package dto

import (
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type ReminderRuleSkuGetPageReq struct {
	dto.Pagination `search:"-"`
	ReminderRuleId int    `form:"reminderRuleId"  search:"type:exact;column:reminder_rule_id;table:reminder_rule_sku" comment:"sxyz_reminder_rule 补货提醒规则表id"`
	Status         string `form:"status"  search:"type:exact;column:status;table:reminder_rule_sku" comment:"状态 1启用 0未启用 全部不传"`
	ReminderRuleSkuOrder
}

type ReminderRuleSkuOrder struct {
	Id                 string `form:"idOrder"  search:"type:order;column:id;table:reminder_rule_sku"`
	ReminderRuleId     string `form:"reminderRuleIdOrder"  search:"type:order;column:reminder_rule_id;table:reminder_rule_sku"`
	SkuCode            string `form:"skuCodeOrder"  search:"type:order;column:sku_code;table:reminder_rule_sku"`
	WarningValue       string `form:"warningValueOrder"  search:"type:order;column:warning_value;table:reminder_rule_sku"`
	ReplenishmentValue string `form:"replenishmentValueOrder"  search:"type:order;column:replenishment_value;table:reminder_rule_sku"`
	Status             string `form:"statusOrder"  search:"type:order;column:status;table:reminder_rule_sku"`
	CreatedAt          string `form:"createdAtOrder"  search:"type:order;column:created_at;table:reminder_rule_sku"`
	UpdatedAt          string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:reminder_rule_sku"`
	CreateBy           string `form:"createByOrder"  search:"type:order;column:create_by;table:reminder_rule_sku"`
	UpdateBy           string `form:"updateByOrder"  search:"type:order;column:update_by;table:reminder_rule_sku"`
	CreateByName       string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:reminder_rule_sku"`
	UpdateByName       string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:reminder_rule_sku"`
	DeletedAt          string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:reminder_rule_sku"`
}

func (m *ReminderRuleSkuGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ReminderRuleSkuInsertReq struct {
	Id                 int    `json:"-" comment:"id"` // id
	ReminderRuleId     int    `json:"-" comment:"sxyz_reminder_rule 补货提醒规则表id"`
	SkuCode            string `json:"skuCode" comment:"sku" vd:"len($)>0"`
	WarningValue       int    `json:"warningValue" comment:"预警值" vd:"$>=0"`
	ReplenishmentValue int    `json:"replenishmentValue" comment:"设置的备货量"`
	Status             string `json:"status" comment:"状态 1启用 0未启用"`
	common.ControlBy   `json:"-"`
}

type ReminderRuleSkuInsertsReq struct {
	Id                 int    `json:"-" comment:"id"` // id
	ReminderRuleId     int    `json:"-" comment:"sxyz_reminder_rule 补货提醒规则表id"`
	SkuCode            string `json:"skuCode" comment:"sku" vd:"len($)>0"`
	WarningValue       int    `json:"warningValue" comment:"预警值" vd:"$>=0"`
	ReplenishmentValue int    `json:"replenishmentValue" comment:"设置的备货量"`
	Status             string `json:"status" comment:"状态 1启用 0未启用"`
	common.ControlBy   `json:"-"`
}

func (s *ReminderRuleSkuInsertReq) Generate(model *models.ReminderRuleSku) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ReminderRuleId = s.ReminderRuleId
	model.SkuCode = s.SkuCode
	model.WarningValue = s.WarningValue
	model.ReplenishmentValue = s.ReplenishmentValue
	model.Status = s.Status
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
	model.UpdateByName = s.UpdateByName
}

func (s *ReminderRuleSkuInsertReq) GetId() interface{} {
	return s.Id
}

type ReminderRuleSkuUpdateReq struct {
	Id                 int    `uri:"id" json:"id" comment:"id"` // id
	SkuCode            string `json:"skuCode" comment:"sku" vd:"len($)>0" `
	WarningValue       int    `json:"warningValue" comment:"预警值" vd:"$>=0" `
	ReplenishmentValue int    `json:"replenishmentValue" comment:"设置的备货量" vd:"$>=0" `
	Status             string `json:"status" comment:"状态 1启用 0未启用"`
	common.ControlBy
}

func (s *ReminderRuleSkuUpdateReq) Generate(model *models.ReminderRuleSku) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	//model.ReminderRuleId = s.ReminderRuleId
	model.SkuCode = s.SkuCode
	model.WarningValue = s.WarningValue
	model.ReplenishmentValue = s.ReplenishmentValue
	model.Status = s.Status
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.CreateByName = s.CreateByName
	model.UpdateByName = s.UpdateByName
}

func (s *ReminderRuleSkuUpdateReq) GetId() interface{} {
	return s.Id
}

// ReminderRuleSkuGetReq 功能获取请求参数
type ReminderRuleSkuGetReq struct {
	Id int `uri:"id"`
}

func (s *ReminderRuleSkuGetReq) GetId() interface{} {
	return s.Id
}

// ReminderRuleSkuDeleteReq 功能删除请求参数
type ReminderRuleSkuDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ReminderRuleSkuDeleteReq) GetId() interface{} {
	return s.Ids
}
