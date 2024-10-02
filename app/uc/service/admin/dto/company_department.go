package dto

import (
	"database/sql"
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CompanyDepartmentGetPageReq struct {
	dto.Pagination   `search:"-"`
	Id               int    `form:"id"  search:"type:exact;column:id;table:company_department" comment:"id"`
	Ids              string `form:"ids"  search:"-" comment:"ids"`
	Name             string `form:"name"  search:"type:contains;column:name;table:company_department" comment:"部门名称"`
	Level            int    `form:"level"  search:"type:exact;column:level;table:company_department" comment:"部门层级"`
	FId              int    `form:"fId"  search:"type:exact;column:f_id;table:company_department" comment:"父级部门id"`
	TopId            int    `form:"topId"  search:"type:exact;column:top_id;table:company_department" comment:"一级部门ID"`
	CompanyId        int    `form:"companyId"  search:"type:exact;column:company_id;table:company_department" comment:"公司id"`
	PersonalBudget   string `form:"personalBudget"  search:"type:exact;column:personal_budget;table:company_department" comment:"人均预算 null不设置无限预算 >=0 有预算限制 "`
	DepartmentBudget string `form:"departmentBudget"  search:"type:exact;column:department_budget;table:company_department" comment:"部门预算 null不设置无限预算 >=0 有预算限制 "`
	QueryFid         int    `form:"queryFid" search:"-" comment:"查询自己和下级部门"`
	CompanyDepartmentOrder
}

type CompanyDepartmentGetPageData struct {
	Id               int      `json:"id"`
	Name             string   `json:"name"`             // 部门名称
	Level            int      `json:"level"`            // 部门层级
	FId              int      `json:"fId"`              // 父级部门id
	TopId            int      `json:"topId"`            // 一级部门ID
	CompanyId        int      `json:"companyId"`        // 公司id
	PersonalBudget   *float64 `json:"personalBudget"`   // 人均预算 null不设置无限预算 >=0 有预算限制
	DepartmentBudget *float64 `json:"departmentBudget"` // 部门预算 null不设置无限预算 >=0 有预算限制

	common.ModelTime
	common.ControlBy
}

type CompanyDepartmentGetPageListData struct {
	Id               int      `json:"id"`
	Name             string   `json:"name"`             // 部门名称
	Level            int      `json:"level"`            // 部门层级
	FId              int      `json:"fId"`              // 父级部门id
	TopId            int      `json:"topId"`            // 一级部门ID
	CompanyId        int      `json:"companyId"`        // 公司id
	PersonalBudget   *float64 `json:"personalBudget"`   // 人均预算 null不设置无限预算 >=0 有预算限制
	DepartmentBudget *float64 `json:"departmentBudget"` // 部门预算 null不设置无限预算 >=0 有预算限制
	CompanyName      string   `json:"companyName"`      // 公司名
	FDepartmentName  string   `json:"fDepartmentName"`  // 上级部门名称

	common.ModelTime
	common.ControlBy
}

type CompanyDepartmentOrder struct {
	Id               string `form:"idOrder"  search:"type:order;column:id;table:company_department"`
	Name             string `form:"nameOrder"  search:"type:order;column:name;table:company_department"`
	Level            string `form:"levelOrder"  search:"type:order;column:level;table:company_department"`
	FId              string `form:"fIdOrder"  search:"type:order;column:f_id;table:company_department"`
	TopId            string `form:"topIdOrder"  search:"type:order;column:top_id;table:company_department"`
	CompanyId        string `form:"companyIdOrder"  search:"type:order;column:company_id;table:company_department"`
	PersonalBudget   string `form:"personalBudgetOrder"  search:"type:order;column:personal_budget;table:company_department"`
	DepartmentBudget string `form:"departmentBudgetOrder"  search:"type:order;column:department_budget;table:company_department"`
	CreatedAt        string `form:"createdAtOrder"  search:"type:order;column:created_at;table:company_department"`
	UpdatedAt        string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:company_department"`
	CreateBy         string `form:"createByOrder"  search:"type:order;column:create_by;table:company_department"`
	UpdateBy         string `form:"updateByOrder"  search:"type:order;column:update_by;table:company_department"`
	DeletedAt        string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:company_department"`
}

func (m *CompanyDepartmentGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CompanyDepartmentGetPageRes struct {
	List []*CompanyDepartmentGetPageCompanyData `json:"list"`
}

type CompanyDepartmentGetPageCompanyData struct {
	models.CompanyInfo
	DepartmentList []*CompanyDepartmentGetPageDepartmentData `json:"departmentList"`
}

type CompanyDepartmentGetPageDepartmentData struct {
	Id               int    `json:"id" comment:""` //
	Name             string `json:"name" comment:"部门名称"`
	Level            int    `json:"level" comment:"部门层级"`
	FId              int    `json:"fId" comment:"父级部门id"`
	FName            int    `json:"FName" comment:"父级部门名称"`
	TopId            int    `json:"topId" comment:"一级部门ID"`
	CompanyId        int    `json:"companyId" comment:"公司id"`
	PersonalBudget   string `json:"personalBudget" comment:"人均预算 null不设置无限预算 >=0 有预算限制 "`
	DepartmentBudget string `json:"departmentBudget" comment:"部门预算 null不设置无限预算 >=0 有预算限制 "`
	common.ControlBy
}

type CompanyDepartmentInsertReq struct {
	Id               int      `json:"-" comment:""` //
	Name             string   `json:"name" comment:"部门名称" vd:"len($)>0 && len($)<=50" `
	FId              int      `json:"fId" comment:"父级部门id"`
	CompanyId        int      `json:"companyId" comment:"公司id" vd:"$>0" `
	PersonalBudget   *float64 `json:"personalBudget" comment:"人均预算 null不设置无限预算 >=0 有预算限制 "`
	DepartmentBudget *float64 `json:"departmentBudget" comment:"部门预算 null不设置无限预算 >=0 有预算限制 "`

	TopId int `json:"-" comment:"一级部门ID"`

	common.ControlBy `json:"-"`
}

func (s *CompanyDepartmentInsertReq) Generate(model *models.CompanyDepartment) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.Name = s.Name
	model.FId = s.FId
	model.TopId = s.FId
	model.CompanyId = s.CompanyId
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的

	if s.PersonalBudget != nil {
		//f, _ := strconv.ParseFloat(s.PersonalBudget, 64)
		model.PersonalBudget = sql.NullFloat64{
			Float64: *s.PersonalBudget,
			Valid:   true,
		}
	} else {
		model.PersonalBudget = sql.NullFloat64{
			Float64: 0,
			Valid:   false,
		}
	}
	if s.DepartmentBudget != nil {
		//f, _ := strconv.ParseFloat(s.DepartmentBudget, 64)
		model.DepartmentBudget = sql.NullFloat64{
			Float64: *s.DepartmentBudget,
			Valid:   true,
		}
	} else {
		model.DepartmentBudget = sql.NullFloat64{
			Float64: 0,
			Valid:   false,
		}
	}
}

func (s *CompanyDepartmentInsertReq) GetId() interface{} {
	return s.Id
}

func (s *CompanyDepartmentInsertReq) check() interface{} {
	return s.Id
}

type CompanyDepartmentUpdateReq struct {
	Id               int      `uri:"id" comment:""` //
	FId              int      `json:"fId" comment:"父级部门id"`
	PersonalBudget   *float64 `json:"personalBudget" comment:"人均预算 null不设置无限预算 >=0 有预算限制 "`
	DepartmentBudget *float64 `json:"departmentBudget" comment:"部门预算 null不设置无限预算 >=0 有预算限制 "`
	common.ControlBy `json:"-"`
}

func (s *CompanyDepartmentUpdateReq) Generate(model *models.CompanyDepartment) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.FId = s.FId
	model.TopId = s.FId
	if s.PersonalBudget != nil {
		//f, _ := strconv.ParseFloat(s.PersonalBudget, 64)
		model.PersonalBudget = sql.NullFloat64{
			Float64: *s.PersonalBudget,
			Valid:   true,
		}
	} else {
		model.PersonalBudget = sql.NullFloat64{
			Float64: 0,
			Valid:   false,
		}
	}
	if s.DepartmentBudget != nil {
		//f, _ := strconv.ParseFloat(s.DepartmentBudget, 64)
		model.DepartmentBudget = sql.NullFloat64{
			Float64: *s.DepartmentBudget,
			Valid:   true,
		}
	} else {
		model.DepartmentBudget = sql.NullFloat64{
			Float64: 0,
			Valid:   false,
		}
	}
}

func (s *CompanyDepartmentUpdateReq) GetId() interface{} {
	return s.Id
}

type CompanyDepartmentUpdateBudgetReq struct {
	Id               int      `uri:"id" comment:""` //
	PersonalBudget   *float64 `json:"personalBudget" comment:"人均预算 null不设置无限预算 >=0 有预算限制 "`
	DepartmentBudget *float64 `json:"departmentBudget" comment:"部门预算 null不设置无限预算 >=0 有预算限制 "`
	common.ControlBy
}

func (s *CompanyDepartmentUpdateBudgetReq) Generate(model *models.CompanyDepartment) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}

	if s.PersonalBudget != nil {
		//f, _ := strconv.ParseFloat(s.PersonalBudget, 64)
		model.PersonalBudget = sql.NullFloat64{
			Float64: *s.PersonalBudget,
			Valid:   true,
		}
	}
	if s.DepartmentBudget != nil {
		//f, _ := strconv.ParseFloat(s.PersonalBudget, 64)
		model.DepartmentBudget = sql.NullFloat64{
			Float64: *s.PersonalBudget,
			Valid:   true,
		}
	}
}

func (s *CompanyDepartmentUpdateBudgetReq) GetId() interface{} {
	return s.Id
}

// CompanyDepartmentGetReq 功能获取请求参数
type CompanyDepartmentGetReq struct {
	Id int `uri:"id"`
}

func (s *CompanyDepartmentGetReq) GetId() interface{} {
	return s.Id
}

type CompanyDepartmentGetRes struct {
	Id               int      `json:"id"`
	Name             string   `json:"name"`             // 部门名称
	Level            int      `json:"level"`            // 部门层级
	FId              *int     `json:"fId"`              // 父级部门id
	CompanyId        int      `json:"companyId"`        // 公司id
	PersonalBudget   *float64 `json:"personalBudget"`   // 人均预算 null不设置无限预算 >=0 有预算限制
	DepartmentBudget *float64 `json:"departmentBudget"` // 部门预算 null不设置无限预算 >=0 有预算限制

	common.ModelTime
	common.ControlBy
}

// CompanyDepartmentDeleteReq 功能删除请求参数
type CompanyDepartmentDeleteReq struct {
	Id int `json:"id"`
}

func (s *CompanyDepartmentDeleteReq) GetId() int {
	return s.Id
}

// CompanyDepartmentIsCanDeleteRes
type CompanyDepartmentIsCanDeleteRes struct {
	IsCanDelete bool `json:"isCanDelete"`
}

type ImportReq struct {
	Objects string `json:"objects"`
}

type CompanyDepartmentImportData struct {
	CompanyId        string `json:"companyId" vd:"len($)>0" ` // 公司id
	CompanyName      string `json:"companyName"`              // 公司名
	Name             string `json:"name" vd:"len($)>0"`       // 部门名称
	DepartmentBudget string `json:"departmentBudget"`         // 部门预算 null不设置无限预算 >=0 有预算限制
	PersonalBudget   string `json:"personalBudget"`           // 人均预算 null不设置无限预算 >=0 有预算限制
	FDepartmentName  string `json:"fDepartmentName"`          // 上级部门名称

	common.ModelTime
	common.ControlBy
}

type ExportRes struct {
	Tpl       string      `json:"tpl"`
	FileName  string      `json:"file_name"`
	Data      interface{} `json:"data"`
	PageTotal int         `json:"page_total"`
}

type CompanyDepartmentExportData struct {
	CompanyId        int    `json:"companyId"`        // 公司id
	CompanyName      string `json:"companyName"`      // 公司名
	Name             string `json:"name"`             // 部门名称
	FDepartmentName  string `json:"fDepartmentName"`  // 上级部门名称
	DepartmentBudget string `json:"departmentBudget"` // 部门预算 null不设置无限预算 >=0 有预算限制
	PersonalBudget   string `json:"personalBudget"`   // 人均预算 null不设置无限预算 >=0 有预算限制
}
