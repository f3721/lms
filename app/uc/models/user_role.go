package models

import (
	"go-admin/common/models"
)

type UserRole struct {
	models.Model

	UserId int `json:"userId" gorm:"type:int unsigned;comment:用户ID"`
	RoleId int `json:"roleId" gorm:"type:int unsigned;comment:角色ID"`
	models.ModelTime
	models.ControlBy
}

func (UserRole) TableName() string {
	return "user_role"
}

func (e *UserRole) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *UserRole) GetId() interface{} {
	return e.Id
}
