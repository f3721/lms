package dto

import (
     
     

	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type UserRoleGetPageReq struct {
	dto.Pagination     `search:"-"`
    UserId int `form:"userId"  search:"type:exact;column:user_id;table:user_role" comment:"用户ID"`
    RoleId int `form:"roleId"  search:"type:exact;column:role_id;table:user_role" comment:"角色ID"`
    UserRoleOrder
}

type UserRoleOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:user_role"`
    UserId string `form:"userIdOrder"  search:"type:order;column:user_id;table:user_role"`
    RoleId string `form:"roleIdOrder"  search:"type:order;column:role_id;table:user_role"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:user_role"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:user_role"`
    UpdateBy string `form:"updateByOrder"  search:"type:order;column:update_by;table:user_role"`
    UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:user_role"`
    DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:user_role"`
    
}

func (m *UserRoleGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type UserRoleInsertReq struct {
    Id int `json:"-" comment:"编号"` // 编号
    UserId int `json:"userId" comment:"用户ID"`
    RoleId int `json:"roleId" comment:"角色ID"`
    common.ControlBy
}

func (s *UserRoleInsertReq) Generate(model *models.UserRole)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.UserId = s.UserId
    model.RoleId = s.RoleId
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *UserRoleInsertReq) GetId() interface{} {
	return s.Id
}

type UserRoleUpdateReq struct {
    Id int `uri:"id" comment:"编号"` // 编号
    UserId int `json:"userId" comment:"用户ID"`
    RoleId int `json:"roleId" comment:"角色ID"`
    common.ControlBy
}

func (s *UserRoleUpdateReq) Generate(model *models.UserRole)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.UserId = s.UserId
    model.RoleId = s.RoleId
    model.UpdateBy = s.UpdateBy
}

func (s *UserRoleUpdateReq) GetId() interface{} {
	return s.Id
}

// UserRoleGetReq 功能获取请求参数
type UserRoleGetReq struct {
     Id int `uri:"id"`
}
func (s *UserRoleGetReq) GetId() interface{} {
	return s.Id
}

// UserRoleDeleteReq 功能删除请求参数
type UserRoleDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *UserRoleDeleteReq) GetId() interface{} {
	return s.Ids
}
