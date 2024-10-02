package dto

import (

	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
    "gorm.io/gorm"
)

type UserBudgetGetPageReq struct {
	dto.Pagination     `search:"-"`
    UserBudgetOrder

    UserId int `form:"userId"  search:"type:exact;column:user_id;table:user_budget"`
    UserName string `form:"userName" search:"-" gorm:"type:string;comment:用户姓名"`
    Month string `form:"month"  search:"type:exact;column:month;table:user_budget"`
    DepId int `form:"depId" search:"-" gorm:"type:int;comment:部门Id"`
}

type UserBudgetOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:user_budget"`
    Month string `form:"monthOrder"  search:"type:order;column:month;table:user_budget"`
    //UserId string `form:"userIdOrder"  search:"type:order;column:user_id;table:user_budget"`
    //InitialAmount string `form:"initialAmountOrder"  search:"type:order;column:initial_amount;table:user_budget"`
    //Spent string `form:"spentOrder"  search:"type:order;column:spent;table:user_budget"`
    //Balance string `form:"balanceOrder"  search:"type:order;column:balance;table:user_budget"`
    //ExcessAmount string `form:"excessAmountOrder"  search:"type:order;column:excess_amount;table:user_budget"`
    //CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:user_budget"`
    //UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:user_budget"`
    
}

func (m *UserBudgetGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type UserBudgetInsertReq struct {
    Id int `json:"-" comment:"自增ID"` // 自增ID
    Month string `json:"month" comment:"年月"`
    UserId int `json:"userId" comment:"用户ID"`
    InitialAmount float64 `json:"initialAmount" comment:"预算金额"`
    Spent float64 `json:"spent" comment:"已使用金额"`
    Balance float64 `json:"balance" comment:"剩余金额"`
    ExcessAmount float64 `json:"excessAmount" comment:"超出金额"`
    common.ControlBy
}

func (s *UserBudgetInsertReq) Generate(model *models.UserBudget)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.Month = s.Month
    model.UserId = s.UserId
    model.InitialAmount = s.InitialAmount
    model.Spent = s.Spent
    model.Balance = s.Balance
    model.ExcessAmount = s.ExcessAmount
}

func (s *UserBudgetInsertReq) GetId() interface{} {
	return s.Id
}

type UserBudgetUpdateReq struct {
    Id int `uri:"id" comment:"自增ID"` // 自增ID
    Month string `json:"month" comment:"年月"`
    UserId int `json:"userId" comment:"用户ID"`
    InitialAmount float64 `json:"initialAmount" comment:"预算金额"`
    Spent float64 `json:"spent" comment:"已使用金额"`
    Balance float64 `json:"balance" comment:"剩余金额"`
    ExcessAmount float64 `json:"excessAmount" comment:"超出金额"`
    common.ControlBy
}

func (s *UserBudgetUpdateReq) Generate(model *models.UserBudget)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.Month = s.Month
    model.UserId = s.UserId
    model.InitialAmount = s.InitialAmount
    model.Spent = s.Spent
    model.Balance = s.Balance
    model.ExcessAmount = s.ExcessAmount
}

func (s *UserBudgetUpdateReq) GetId() interface{} {
	return s.Id
}

// UserBudgetGetReq 功能获取请求参数
type UserBudgetGetReq struct {
     Id int `uri:"id"`
}
func (s *UserBudgetGetReq) GetId() interface{} {
	return s.Id
}

// UserBudgetDeleteReq 功能删除请求参数
type UserBudgetDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *UserBudgetDeleteReq) GetId() interface{} {
	return s.Ids
}

type UserBudgetGetPageResp struct {
    models.UserBudget

    UserName string `json:"userName" gorm:"type:varchar(50);comment:用户姓名"`
}

func UserBudgetGetPageMakeCondition(c *UserBudgetGetPageReq) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        if c.UserName != "" {
            db.Where("ui.user_name like ?", "%" + c.UserName + "%")
        }
        if c.DepId > 0 {
            db.Where("ui.company_department_id = ?", c.DepId)
        }
        return db
    }
}
