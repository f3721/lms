package models

import (
	"go-admin/common/models"
)

type RoleInfo struct {
	models.Model

	RoleName      string `json:"roleName" gorm:"type:varchar(20);comment:角色名称"`
	RoleKey       string `json:"roleKey" gorm:"type:varchar(128);comment:RoleKey"`
	RoleStatus    int    `json:"roleStatus" gorm:"type:tinyint(1);comment:RoleStatus"`
	DataScope     string `json:"dataScope" gorm:"type:varchar(128);comment:DataScope"`
	ManageCompany string `json:"manageCompany" gorm:"type:varchar(20);comment:判断公司是否可以管理该权限(1xx:punchout管理,x1x:EAS管理,xx1:普通管理)"`
	Menus         string `json:"menus" gorm:"type:varchar(128);comment:菜单权限"`
	models.ModelTime
	models.ControlBy
}

func (RoleInfo) TableName() string {
	return "role_info"
}

func (e *RoleInfo) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *RoleInfo) GetId() interface{} {
	return e.Id
}
