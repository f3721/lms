package dto

import (
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type ReminderRuleSkuLogGetPageReq struct {
	dto.Pagination    `search:"-"`
	ReminderRuleSkuId int64 `form:"reminderRuleSkuId"  search:"type:exact;column:reminder_rule_sku_id;table:reminder_rule_sku_log" comment:"sku补货提醒规则表id"`
	ReminderRuleSkuLogOrder
}

type ReminderRuleSkuLogOrder struct {
	Id                 string `form:"idOrder"  search:"type:order;column:id;table:reminder_rule_sku_log"`
	ReminderRuleSkuId  string `form:"reminderRuleSkuIdOrder"  search:"type:order;column:reminder_rule_sku_id;table:reminder_rule_sku_log"`
	SkuCode            string `form:"skuCodeOrder"  search:"type:order;column:sku_code;table:reminder_rule_sku_log"`
	WarningValue       string `form:"warningValueOrder"  search:"type:order;column:warning_value;table:reminder_rule_sku_log"`
	ReplenishmentValue string `form:"replenishmentValueOrder"  search:"type:order;column:replenishment_value;table:reminder_rule_sku_log"`
	Status             string `form:"statusOrder"  search:"type:order;column:status;table:reminder_rule_sku_log"`
	LogType            string `form:"logTypeOrder"  search:"type:order;column:log_type;table:reminder_rule_sku_log"`
	CreatedAt          string `form:"createdAtOrder"  search:"type:order;column:created_at;table:reminder_rule_sku_log"`
	CreateBy           string `form:"createByOrder"  search:"type:order;column:create_by;table:reminder_rule_sku_log"`
	CreateByName       string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:reminder_rule_sku_log"`
}

func (m *ReminderRuleSkuLogGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ReminderRuleSkuLogInsertReq struct {
	Id                 int    `json:"-" comment:"id"` // id
	ReminderRuleSkuId  int    `json:"reminderRuleSkuId" comment:"sku补货提醒规则表id"`
	SkuCode            string `json:"skuCode" comment:"sku"`
	WarningValue       int    `json:"warningValue" comment:"预警值"`
	ReplenishmentValue int    `json:"replenishmentValue" comment:"设置的备货量"`
	Status             string `json:"status" comment:"状态 1启用 0未启用"`
	LogType            int    `json:"-" comment:"操作类型 1创建 2修改 3删除"`
	CreateByName       string `json:"createByName" comment:""`
	common.ControlBy
}

func (s *ReminderRuleSkuLogInsertReq) Generate(model *models.ReminderRuleSkuLog) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ReminderRuleSkuId = s.ReminderRuleSkuId
	model.SkuCode = s.SkuCode
	model.WarningValue = s.WarningValue
	model.ReplenishmentValue = s.ReplenishmentValue
	model.Status = s.Status
	model.LogType = s.LogType
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *ReminderRuleSkuLogInsertReq) GetId() interface{} {
	return s.Id
}

type ReminderRuleSkuLogUpdateReq struct {
	Id                 int    `uri:"id" comment:"id"` // id
	ReminderRuleSkuId  int    `json:"reminderRuleSkuId" comment:"sku补货提醒规则表id"`
	SkuCode            string `json:"skuCode" comment:"sku"`
	WarningValue       int    `json:"warningValue" comment:"预警值"`
	ReplenishmentValue int    `json:"replenishmentValue" comment:"设置的备货量"`
	Status             string `json:"status" comment:"状态 1启用 0未启用"`
	LogType            int    `json:"logType" comment:"操作类型 1创建 2修改 3删除"`
	CreateByName       string `json:"createByName" comment:""`
	common.ControlBy
}

func (s *ReminderRuleSkuLogUpdateReq) Generate(model *models.ReminderRuleSkuLog) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ReminderRuleSkuId = s.ReminderRuleSkuId
	model.SkuCode = s.SkuCode
	model.WarningValue = s.WarningValue
	model.ReplenishmentValue = s.ReplenishmentValue
	model.Status = s.Status
	model.LogType = s.LogType
	model.CreateByName = s.CreateByName
}

func (s *ReminderRuleSkuLogUpdateReq) GetId() interface{} {
	return s.Id
}

// ReminderRuleSkuLogGetReq 功能获取请求参数
type ReminderRuleSkuLogGetReq struct {
	Id int `uri:"id"`
}

func (s *ReminderRuleSkuLogGetReq) GetId() interface{} {
	return s.Id
}

// ReminderRuleSkuLogDeleteReq 功能删除请求参数
type ReminderRuleSkuLogDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ReminderRuleSkuLogDeleteReq) GetId() interface{} {
	return s.Ids
}
