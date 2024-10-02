package dto

import (
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	cModels "go-admin/common/models"
	common "go-admin/common/models"
)

type ManageMenuGetPageReq struct {
	dto.Pagination `search:"-"`
	ManageMenuOrder
}

type ManageMenuOrder struct {
	Id               string `form:"idOrder"  search:"type:order;column:id;table:manage_menu"`
	GroupName        string `form:"groupNameOrder"  search:"type:order;column:group_name;table:manage_menu"`
	Title            string `form:"titleOrder"  search:"type:order;column:title;table:manage_menu"`
	IsActive         string `form:"isActiveOrder"  search:"type:order;column:is_active;table:manage_menu"`
	RouteName        string `form:"routeNameOrder"  search:"type:order;column:route_name;table:manage_menu"`
	ActiveRouteNames string `form:"activeRouteNamesOrder"  search:"type:order;column:active_route_names;table:manage_menu"`
	Order            string `form:"orderOrder"  search:"type:order;column:order;table:manage_menu"`
}

func (m *ManageMenuGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ManageMenuInsertReq struct {
	Id               int          `json:"-" comment:"主键ID"` // 主键ID
	GroupName        string       `json:"groupName" comment:"父级名称"`
	Title            string       `json:"title" comment:"菜单名称"`
	IsActive         int          `json:"isActive" comment:"是否选中"`
	RouteName        string       `json:"routeName" comment:"路由地址"`
	ActiveRouteNames cModels.Strs `json:"activeRouteNames" comment:"选中高亮标识"`
	Order            int          `json:"order" comment:"排序"`
	common.ControlBy
}

func (s *ManageMenuInsertReq) Generate(model *models.ManageMenu) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.GroupName = s.GroupName
	model.Title = s.Title
	model.IsActive = s.IsActive
	model.RouteName = s.RouteName
	model.ActiveRouteNames = s.ActiveRouteNames
	model.OrderBy = s.Order
}

func (s *ManageMenuInsertReq) GetId() interface{} {
	return s.Id
}

type ManageMenuUpdateReq struct {
	Id               int          `uri:"id" comment:"主键ID"` // 主键ID
	GroupName        string       `json:"groupName" comment:"父级名称"`
	Title            string       `json:"title" comment:"菜单名称"`
	IsActive         int          `json:"isActive" comment:"是否选中"`
	RouteName        string       `json:"routeName" comment:"路由地址"`
	ActiveRouteNames cModels.Strs `json:"activeRouteNames" comment:"选中高亮标识"`
	Order            int          `json:"order" comment:"排序"`
	common.ControlBy
}

func (s *ManageMenuUpdateReq) Generate(model *models.ManageMenu) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.GroupName = s.GroupName
	model.Title = s.Title
	model.IsActive = s.IsActive
	model.RouteName = s.RouteName
	model.ActiveRouteNames = s.ActiveRouteNames
	model.OrderBy = s.Order
}

func (s *ManageMenuUpdateReq) GetId() interface{} {
	return s.Id
}

// ManageMenuGetReq 功能获取请求参数
type ManageMenuGetReq struct {
	Id int `uri:"id"`
}

func (s *ManageMenuGetReq) GetId() interface{} {
	return s.Id
}

// ManageMenuDeleteReq 功能删除请求参数
type ManageMenuDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ManageMenuDeleteReq) GetId() interface{} {
	return s.Ids
}
