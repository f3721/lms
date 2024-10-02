package dto

import (
     
     
     

	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type RoleInfoGetPageReq struct {
	dto.Pagination     `search:"-"`
    RoleName string `form:"roleName"  search:"type:exact;column:role_name;table:role_info" comment:"角色名称"`
    RoleStatus int `form:"roleStatus"  search:"type:exact;column:role_status;table:role_info" comment:""`
    ManageCompany string `form:"manageCompany"  search:"type:exact;column:manage_company;table:role_info" comment:"判断公司是否可以管理该权限(1xx:punchout管理,x1x:EAS管理,xx1:普通管理)"`
    RoleInfoOrder
}

type RoleInfoOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:role_info"`
    RoleName string `form:"roleNameOrder"  search:"type:order;column:role_name;table:role_info"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:role_info"`
    UpdateBy string `form:"updateByOrder"  search:"type:order;column:update_by;table:role_info"`
    RoleStatus string `form:"roleStatusOrder"  search:"type:order;column:role_status;table:role_info"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:role_info"`
    UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:role_info"`
    ManageCompany string `form:"manageCompanyOrder"  search:"type:order;column:manage_company;table:role_info"`
    
}

func (m *RoleInfoGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type RoleInfoInsertReq struct {
    Id int `json:"-" comment:"编号"` // 编号
    RoleName string `json:"roleName" comment:"角色名称"`
    RoleStatus int `json:"roleStatus" comment:""`
    ManageCompany string `json:"manageCompany" comment:"判断公司是否可以管理该权限(1xx:punchout管理,x1x:EAS管理,xx1:普通管理)"`
    common.ControlBy
}

func (s *RoleInfoInsertReq) Generate(model *models.RoleInfo)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.RoleName = s.RoleName
    model.RoleStatus = s.RoleStatus
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.ManageCompany = s.ManageCompany
}

func (s *RoleInfoInsertReq) GetId() interface{} {
	return s.Id
}

type RoleInfoUpdateReq struct {
    Id int `uri:"id" comment:"编号"` // 编号
    RoleName string `json:"roleName" comment:"角色名称"`
    RoleStatus int `json:"roleStatus" comment:""`
    ManageCompany string `json:"manageCompany" comment:"判断公司是否可以管理该权限(1xx:punchout管理,x1x:EAS管理,xx1:普通管理)"`
    common.ControlBy
}

func (s *RoleInfoUpdateReq) Generate(model *models.RoleInfo)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.RoleName = s.RoleName
    model.UpdateBy = s.UpdateBy
    model.RoleStatus = s.RoleStatus
    model.ManageCompany = s.ManageCompany
}

func (s *RoleInfoUpdateReq) GetId() interface{} {
	return s.Id
}

// RoleInfoGetReq 功能获取请求参数
type RoleInfoGetReq struct {
     Id int `uri:"id"`
}
func (s *RoleInfoGetReq) GetId() interface{} {
	return s.Id
}

// RoleInfoDeleteReq 功能删除请求参数
type RoleInfoDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *RoleInfoDeleteReq) GetId() interface{} {
	return s.Ids
}
