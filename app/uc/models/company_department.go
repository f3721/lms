package models

import (
	"database/sql"
	"go-admin/common/models"
)

const (
	CompanyDepartmentOperationModel  = "companyDepartment"
	CompanyDepartmentOperationCreate = "Create"
	CompanyDepartmentOperationUpdate = "Update"
	CompanyDepartmentOperationDelete = "Delete"
)

type CompanyDepartment struct {
	models.Model
	Name             string          `json:"name" gorm:"type:varchar(200)"`              // 部门名称
	Level            int             `json:"level" gorm:"type:int unsigned"`             // 部门层级
	FId              int             `json:"fId" gorm:"type:int unsigned"`               // 父级部门id
	TopId            int             `json:"topId" gorm:"type:int unsigned"`             // 一级部门ID
	CompanyId        int             `json:"companyId" gorm:"type:int unsigned"`         // 公司id
	PersonalBudget   sql.NullFloat64 `json:"personalBudget" gorm:"type:decimal(10,2)"`   // 人均预算 null不设置无限预算 >=0 有预算限制
	DepartmentBudget sql.NullFloat64 `json:"departmentBudget" gorm:"type:decimal(10,2)"` // 部门预算 null不设置无限预算 >=0 有预算限制

	models.ModelTime
	models.ControlBy
}

func (CompanyDepartment) TableName() string {
	return "company_department"
}

func (e *CompanyDepartment) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *CompanyDepartment) GetId() interface{} {
	return e.Id
}
