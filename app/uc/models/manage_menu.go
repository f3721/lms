package models

import (
	"go-admin/common/models"
)

type ManageMenu struct {
	models.Model

	Type             int         `json:"type" gorm:"type:tinyint(1);comment:类型1-左边菜单 2-右边下拉菜单"`
	GroupName        string      `json:"groupName" gorm:"type:varchar(20);comment:父级名称"`
	Title            string      `json:"title" gorm:"type:varchar(20);comment:菜单名称"`
	IsActive         int         `json:"isActive" gorm:"type:tinyint(1);comment:是否选中"`
	RouteName        string      `json:"routeName" gorm:"type:varchar(30);comment:路由地址"`
	ActiveRouteNames models.Strs `json:"activeRouteNames" gorm:"type:varchar(100);comment:选中高亮标识"`
	OrderBy          int         `json:"orderBy" gorm:"type:int;comment:排序"`
	models.ModelTime
	models.ControlBy
}

// 左侧菜单列表
type RouterMenu struct {
	Title    string       `json:"title" gorm:"-;comment:父级菜单"`
	Children []ManageMenu `json:"children" gorm:"-;comment:子菜单"`
}

func (ManageMenu) TableName() string {
	return "manage_menu"
}

func (e *ManageMenu) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *ManageMenu) GetId() interface{} {
	return e.Id
}
