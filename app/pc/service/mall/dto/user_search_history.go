package dto

import (
	"go-admin/app/pc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"gorm.io/gorm"
)

type UserSearchHistoryGetPageReq struct {
	dto.Pagination `search:"-"`
	UserId         int    `form:"-"  search:"-"`
	WarehouseCode  string `form:"-"  search:"-"`
	UserSearchHistoryOrder
}

type UserSearchHistoryOrder struct {
	Id string `form:"idOrder"  search:"type:order;column:id;table:user_search_history"`
}

func (m *UserSearchHistoryGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type UserSearchHistoryInsertReq struct {
	Id            int    `json:"-" comment:"主键"` // 主键
	UserId        int    `json:"userId" comment:"用户编号"`
	Keyword       string `json:"keyword" comment:"关键词"`
	WarehouseCode string `json:"warehouseCode" comment:"仓库CODE"`
	common.ControlBy
}

func (s *UserSearchHistoryInsertReq) Generate(model *models.UserSearchHistory) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UserId = s.UserId
	model.Keyword = s.Keyword
	model.WarehouseCode = s.WarehouseCode
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func MakeUserSearchHistoryReqCondition(c *UserSearchHistoryGetPageReq) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db.Where("user_id = ?", c.UserId)
		db.Where("warehouse_code = ?", c.WarehouseCode)
		return db
	}
}

type UserSearchHistoryGetPageResp struct {
	Keyword string `json:"keyword"`
}
