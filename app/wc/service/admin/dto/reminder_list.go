package dto

import (
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"time"
)

type ReminderListGetPageReq struct {
	dto.Pagination `search:"-"`
	ReminderRuleId int    `form:"reminderRuleId"  search:"type:exact;column:reminder_rule_id;table:reminder_list" comment:"sxyz_reminder_rule 补货规则表id"`
	CompanyId      int    `form:"companyId"  search:"type:exact;column:company_id;table:reminder_list" comment:"公司id"`
	WarehouseCode  string `form:"warehouseCode"  search:"type:exact;column:warehouse_code;table:reminder_list" comment:"仓库id"`
	VendorId       int    `form:"vendorId"  search:"type:exact;column:vendor_id;table:reminder_list" comment:"货主id"`
	CreatedAtStart string `form:"createdAtStart" search:"type:gte;column:created_at;table:reminder_list" comment:"时间筛选"`
	CreatedAtEnd   string `form:"createdAtEnd" search:"type:lte;column:created_at;table:reminder_list" comment:"时间筛选"`
	ReminderListOrder
}

type ReminderListOrder struct {
	Id             string `form:"idOrder"  search:"type:order;column:id;table:reminder_list"`
	ReminderRuleId string `form:"reminderRuleIdOrder"  search:"type:order;column:reminder_rule_id;table:reminder_list"`
	CompanyId      string `form:"companyIdOrder"  search:"type:order;column:company_id;table:reminder_list"`
	WarehouseCode  string `form:"warehouseCodeOrder"  search:"type:order;column:warehouse_id;table:reminder_list"`
	VendorId       string `form:"VendorIdOrder"  search:"type:order;column:supplier_info_id;table:reminder_list"`
	SkuCount       string `form:"skuCountOrder"  search:"type:order;column:sku_count;table:reminder_list"`
	CreatedAt      string `form:"createdAtOrder"  search:"type:order;column:created_at;table:reminder_list"`
}

func (m *ReminderListGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ReminderListData struct {
	Id             int       `json:"id" comment:"id"` // id
	ReminderRuleId int       `json:"reminderRuleId" comment:"sxyz_reminder_rule 补货规则表id"`
	CompanyId      int       `json:"companyId" comment:"公司id"`
	WarehouseCode  string    `json:"warehouseCode" comment:"仓库id"`
	VendorId       int       `json:"vendorId" comment:"货主id"`
	SkuCount       int       `json:"skuCount" comment:"SKU数量"`
	CreatedAt      time.Time `json:"createdAt"`
	CompanyName    string    `json:"companyName" gorm:"-"`
	WarehouseName  string    `json:"warehouseName" gorm:"-"`
	VendorName     string    `json:"vendorName" gorm:"-"`
}

type ReminderListInsertReq struct {
	Id             int    `json:"-" comment:"id"` // id
	ReminderRuleId int    `json:"reminderRuleId" comment:"sxyz_reminder_rule 补货规则表id"`
	CompanyId      int    `json:"companyId" comment:"公司id"`
	WarehouseCode  string `json:"warehouseCode" comment:"仓库id"`
	VendorId       int    `json:"vendorId" comment:"货主id"`
	SkuCount       int    `json:"skuCount" comment:"SKU数量"`
}

func (s *ReminderListInsertReq) Generate(model *models.ReminderList) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ReminderRuleId = s.ReminderRuleId
	model.CompanyId = s.CompanyId
	model.WarehouseCode = s.WarehouseCode
	model.VendorId = s.VendorId
	model.SkuCount = s.SkuCount
}

func (s *ReminderListInsertReq) GetId() interface{} {
	return s.Id
}

type ReminderListUpdateReq struct {
	Id             int    `uri:"id" comment:"id"` // id
	ReminderRuleId int    `json:"reminderRuleId" comment:"sxyz_reminder_rule 补货规则表id"`
	CompanyId      int    `json:"companyId" comment:"公司id"`
	WarehouseCode  string `json:"warehouseCode" comment:"仓库id"`
	VendorId       int    `json:"vendorId" comment:"货主id"`
	SkuCount       int    `json:"skuCount" comment:"SKU数量"`
	common.ControlBy
}

func (s *ReminderListUpdateReq) Generate(model *models.ReminderList) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ReminderRuleId = s.ReminderRuleId
	model.CompanyId = s.CompanyId
	model.WarehouseCode = s.WarehouseCode
	model.VendorId = s.VendorId
	model.SkuCount = s.SkuCount
}

func (s *ReminderListUpdateReq) GetId() interface{} {
	return s.Id
}

// ReminderListGetReq 功能获取请求参数
type ReminderListGetReq struct {
	Id int `uri:"id"`
}

func (s *ReminderListGetReq) GetId() interface{} {
	return s.Id
}

// ReminderListDeleteReq 功能删除请求参数
type ReminderListDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ReminderListDeleteReq) GetId() interface{} {
	return s.Ids
}

// 补货清单添加结构体
type ReminderListGenuine struct {
	models.StockInfo        //其实是StockInfo 还没有
	Id               int    `json:"id"`
	Stock            int    `json:"stock"`
	VendorId         int    `json:"vendorId"`
	VendorSku        string `json:"vendorSku"`
	SkuCode          string `json:"skuCode"`
	LockStock        int    `json:"lockStock"`
}

// 导出SKU结构体
type ExportSkuData struct {
	CompanyName                 string `json:"companyName"`
	WarehouseName               string `json:"warehouseName"`
	Date                        string `json:"date"`
	SkuCode                     string `json:"skuCode"`
	VendorName                  string `json:"vendorName"`
	VendorSku                   string `json:"vendorSku"`
	WarningValue                int    `json:"warningValue"`
	ReplenishmentValue          int    `json:"replenishmentValue"`
	RecommendReplenishmentValue int    `json:"recommendReplenishmentValue"`
	GenuineStock                int    `json:"genuineStock"`
	AllStock                    int    `json:"allStock"`
	OccupyStock                 int    `json:"occupyStock"`
}

type ExportRes struct {
	Tpl       string      `json:"tpl"`
	FileName  string      `json:"file_name"`
	Data      interface{} `json:"data"`
	PageTotal int         `json:"page_total"`
}
