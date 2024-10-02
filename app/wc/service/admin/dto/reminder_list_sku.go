package dto

import (
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"time"
)

type ReminderListSkuGetPageReq struct {
	dto.Pagination `search:"-"`
	ReminderListId int `form:"reminderListId"  search:"type:exact;column:reminder_list_id;table:reminder_list_sku" comment:"sxyz_reminder_list 补货清单表id"`
	ReminderListSkuOrder
}

type ReminderListSkuOrder struct {
	Id                          string `form:"idOrder"  search:"type:order;column:id;table:reminder_list_sku"`
	ReminderListId              string `form:"reminderListIdOrder"  search:"type:order;column:reminder_list_id;table:reminder_list_sku"`
	SkuCode                     string `form:"skuCodeOrder"  search:"type:order;column:sku_code;table:reminder_list_sku"`
	VendorId                    string `form:"VendorIdOrder"  search:"type:order;column:supplier_info_id;table:reminder_list_sku"`
	VendorSku                   string `form:"VendorSkuOrder"  search:"type:order;column:supplier_info_sku;table:reminder_list_sku"`
	WarningValue                string `form:"warningValueOrder"  search:"type:order;column:warning_value;table:reminder_list_sku"`
	ReplenishmentValue          string `form:"replenishmentValueOrder"  search:"type:order;column:replenishment_value;table:reminder_list_sku"`
	RecommendReplenishmentValue string `form:"recommendReplenishmentValueOrder"  search:"type:order;column:recommend_replenishment_value;table:reminder_list_sku"`
	GenuineStock                string `form:"genuineStockOrder"  search:"type:order;column:genuine_stock;table:reminder_list_sku"`
	AllStock                    string `form:"allStockOrder"  search:"type:order;column:all_stock;table:reminder_list_sku"`
	OccupyStock                 string `form:"occupyStockOrder"  search:"type:order;column:occupy_tock;table:reminder_list_sku"`
	CreatedAt                   string `form:"createdAtOrder"  search:"type:order;column:created_at;table:reminder_list_sku"`
	OrderLackStock              string `form:"orderLackStockOrder"  search:"type:order;column:order_lack_stock;table:reminder_list_sku"`
}

func (m *ReminderListSkuGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ReminderListSkuInsertReq struct {
	Id                          int    `json:"-" comment:"id"` // id
	ReminderListId              int    `json:"reminderListId" comment:"sxyz_reminder_list 补货清单表id"`
	SkuCode                     string `json:"skuCode" comment:"sku"`
	VendorId                    int    `json:"vendorId" comment:"货主id"`
	VendorSku                   string `json:"vendorSku" comment:"货主sku"`
	WarningValue                int    `json:"warningValue" comment:"预警值"`
	ReplenishmentValue          int    `json:"replenishmentValue" comment:"设置的备货量"`
	RecommendReplenishmentValue int    `json:"recommendReplenishmentValue" comment:"建议补货量：建议补货量=备货量-当前正品用数量+订单实际缺货数量"`
	GenuineStock                int    `json:"genuineStock" comment:"当前正品可用数量"`
	AllStock                    int    `json:"allStock" comment:"当前在库数量(正品仓在库数量+次品仓在库数)"`
	OccupyStock                 int    `json:"occupyStock" comment:"当前占用数量"`
	OrderLackStock              int    `json:"orderLackStock" comment:"订单缺货数量"`
	CreatedAt                   time.Time
}

func (s *ReminderListSkuInsertReq) Generate(model *models.ReminderListSku) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ReminderListId = s.ReminderListId
	model.SkuCode = s.SkuCode
	model.VendorId = s.VendorId
	model.VendorSku = s.VendorSku
	model.WarningValue = s.WarningValue
	model.ReplenishmentValue = s.ReplenishmentValue
	model.RecommendReplenishmentValue = s.RecommendReplenishmentValue
	model.GenuineStock = s.GenuineStock
	model.AllStock = s.AllStock
	model.OccupyStock = s.OccupyStock
	model.OrderLackStock = s.OrderLackStock
	model.CreatedAt = s.CreatedAt
}

func (s *ReminderListSkuInsertReq) GetId() interface{} {
	return s.Id
}

type ReminderListSkuUpdateReq struct {
	Id                          int    `uri:"id" comment:"id"` // id
	ReminderListId              int    `json:"reminderListId" comment:"sxyz_reminder_list 补货清单表id"`
	SkuCode                     string `json:"skuCode" comment:"sku"`
	VendorId                    int    `json:"VendorId" comment:"货主id"`
	VendorSku                   string `json:"vendorSku" comment:"货主sku"`
	WarningValue                int    `json:"warningValue" comment:"预警值"`
	ReplenishmentValue          int    `json:"replenishmentValue" comment:"设置的备货量"`
	RecommendReplenishmentValue int    `json:"recommendReplenishmentValue" comment:"建议补货量：建议补货量=备货量-当前正品用数量+订单实际缺货数量"`
	GenuineStock                int    `json:"genuineStock" comment:"当前正品可用数量"`
	AllStock                    int    `json:"allStock" comment:"当前在库数量(正品仓在库数量+次品仓在库数)"`
	OccupyStock                 int    `json:"occupyStock" comment:"当前占用数量"`
	OrderLackStock              int    `json:"orderLackStock" comment:"订单缺货数量"`
	common.ControlBy
}

func (s *ReminderListSkuUpdateReq) Generate(model *models.ReminderListSku) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ReminderListId = s.ReminderListId
	model.SkuCode = s.SkuCode
	model.VendorId = s.VendorId
	model.VendorSku = s.VendorSku
	model.WarningValue = s.WarningValue
	model.ReplenishmentValue = s.ReplenishmentValue
	model.RecommendReplenishmentValue = s.RecommendReplenishmentValue
	model.GenuineStock = s.GenuineStock
	model.AllStock = s.AllStock
	model.OccupyStock = s.OccupyStock
	model.OrderLackStock = s.OrderLackStock
}

func (s *ReminderListSkuUpdateReq) GetId() interface{} {
	return s.Id
}

// ReminderListSkuGetReq 功能获取请求参数
type ReminderListSkuGetReq struct {
	Id int `uri:"id"`
}

func (s *ReminderListSkuGetReq) GetId() interface{} {
	return s.Id
}

// ReminderListSkuDeleteReq 功能删除请求参数
type ReminderListSkuDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ReminderListSkuDeleteReq) GetId() interface{} {
	return s.Ids
}
