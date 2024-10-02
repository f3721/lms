package dto

import (
    "database/sql"
    "go-admin/app/uc/models"
    "go-admin/common/dto"
    common "go-admin/common/models"
    "go-admin/common/utils"
    "gorm.io/gorm"
    "time"
)

type DepartmentBudgetGetPageReq struct {
	dto.Pagination     `search:"-"`
    DepId string `form:"depId"  search:"type:exact;column:dep_id;table:department_budget"`
    ParentDepId int `form:"parentDepId" search:"-" gorm:"type:int;comment:父部门ID"`
    CompanyId int `form:"companyId"  search:"type:exact;column:company_id;table:cd"`
    StartMonth string `form:"startMonth" search:"-" gorm:"type:string;comment:开始月份"`
    EndMonth string `form:"endMonth" search:"-" gorm:"type:string;comment:结束月份"`
    DepartmentBudgetOrder
}

type DepartmentBudgetOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:department_budget"`
    Month string `form:"monthOrder"  search:"type:order;column:month;table:department_budget"`
    //DepId string `form:"depIdOrder"  search:"type:order;column:dep_id;table:department_budget"`
    //InitialAmount string `form:"initialAmountOrder"  search:"type:order;column:initial_amount;table:department_budget"`
    //Spent string `form:"spentOrder"  search:"type:order;column:spent;table:department_budget"`
    //Balance string `form:"balanceOrder"  search:"type:order;column:balance;table:department_budget"`
    //ExcessAmount string `form:"excessAmountOrder"  search:"type:order;column:excess_amount;table:department_budget"`
    //CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:department_budget"`
    //UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:department_budget"`
    
}

func (m *DepartmentBudgetGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type DepartmentBudgetInsertReq struct {
    Id int `json:"-" comment:"自增ID"` // 自增ID
    Month string `json:"month" comment:"年月"`
    DepId int `json:"depId" comment:"部门ID"`
    InitialAmount float64 `json:"initialAmount" comment:"预算金额"`
    Spent float64 `json:"spent" comment:"已使用金额"`
    Balance float64 `json:"balance" comment:"剩余金额"`
    ExcessAmount float64 `json:"excessAmount" comment:"超出金额"`
    common.ControlBy
}

func (s *DepartmentBudgetInsertReq) Generate(model *models.DepartmentBudget)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.Month = s.Month
    model.DepId = s.DepId
    model.InitialAmount = s.InitialAmount
    model.Spent = s.Spent
    model.Balance = s.Balance
    model.ExcessAmount = s.ExcessAmount
}

func (s *DepartmentBudgetInsertReq) GetId() interface{} {
	return s.Id
}

type DepartmentBudgetUpdateReq struct {
    Id int `uri:"id" comment:"自增ID"` // 自增ID
    Month string `json:"month" comment:"年月"`
    DepId int `json:"depId" comment:"部门ID"`
    InitialAmount float64 `json:"initialAmount" comment:"预算金额"`
    Spent float64 `json:"spent" comment:"已使用金额"`
    Balance float64 `json:"balance" comment:"剩余金额"`
    ExcessAmount float64 `json:"excessAmount" comment:"超出金额"`
    common.ControlBy
}

func (s *DepartmentBudgetUpdateReq) Generate(model *models.DepartmentBudget)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.Month = s.Month
    model.DepId = s.DepId
    model.InitialAmount = s.InitialAmount
    model.Spent = s.Spent
    model.Balance = s.Balance
    model.ExcessAmount = s.ExcessAmount
}

func (s *DepartmentBudgetUpdateReq) GetId() interface{} {
	return s.Id
}

// DepartmentBudgetGetReq 功能获取请求参数
type DepartmentBudgetGetReq struct {
     Id int `uri:"id"`
}
func (s *DepartmentBudgetGetReq) GetId() interface{} {
	return s.Id
}

// DepartmentBudgetDeleteReq 功能删除请求参数
type DepartmentBudgetDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *DepartmentBudgetDeleteReq) GetId() interface{} {
	return s.Ids
}

type DepartmentBudgetListResp struct {
    common.Model

    Month string `json:"month" gorm:"type:varchar(10);comment:年月"`
    DepId int `json:"depId" gorm:"type:int;comment:部门ID"`
    InitialAmount float64 `json:"initialAmount" gorm:"type:decimal(10,2);comment:预算金额"`
    Spent float64 `json:"spent" gorm:"type:decimal(10,2);comment:已使用金额"`
    Balance float64 `json:"balance" gorm:"type:decimal(10,2);comment:剩余金额"`
    ExcessAmount float64 `json:"excessAmount" gorm:"type:decimal(10,2);comment:超出金额"`
    CreatedAt time.Time `json:"createdAt" gorm:"comment:创建时间"`
    UpdatedAt time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`
    DepName string `json:"depName" gorm:"type:varchar(200);comment:部门名称"`
    ParentDepName string `json:"parentDepName" gorm:"type:varchar(200);comment:上级部门"`
    PersonalBudget sql.NullFloat64 `json:"personalBudget" gorm:"type:decimal(10,2)"`
    IsUserBudget int `json:"isUserBudget" gorm:"-" comment:"是否使用人均预算"`
    CompanyId int `json:"companyId" gorm:"type:int;comment:公司ID"`
    CompanyName string `json:"companyName" gorm:"type:varchar(50);comment:公司名称"`
}

func GetPageMakeCondition(c *DepartmentBudgetGetPageReq) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        if c.ParentDepId > 0 {
            db.Where("(department_budget.dep_id = ? or cd.f_id = ?)", c.ParentDepId, c.ParentDepId)
        }
        if c.StartMonth != "" && c.EndMonth != "" {
            db.Where("department_budget.month in ?", utils.GetMonthsBetweenDates(c.StartMonth, c.EndMonth))
        }
        return db
    }
}

type DepartmentBudgetGetResp struct {
    models.DepartmentBudget

    DepName string `json:"depName" gorm:"type:varchar(200);comment:部门名称"`
    CompanyName string `json:"companyName" gorm:"type:varchar(50);comment:公司名称"`
}