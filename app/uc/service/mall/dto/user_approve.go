package dto

import (
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type UserApproveGetPageReq struct {
	dto.Pagination `search:"-"`
	UserApproveOrder
}

type UserApproveOrder struct {
	Id           string `form:"idOrder"  search:"type:order;column:id;table:user_approve"`
	UserId       string `form:"userIdOrder"  search:"type:order;column:user_id;table:user_approve"`
	ApproveId    string `form:"approveIdOrder"  search:"type:order;column:approve_id;table:user_approve"`
	Status       string `form:"statusOrder"  search:"type:order;column:status;table:user_approve"`
	CreatedAt    string `form:"createdAtOrder"  search:"type:order;column:created_at;table:user_approve"`
	CreateBy     string `form:"createByOrder"  search:"type:order;column:create_by;table:user_approve"`
	UpdateBy     string `form:"updateByOrder"  search:"type:order;column:update_by;table:user_approve"`
	UpdatedAt    string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:user_approve"`
	DeletedAt    string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:user_approve"`
	CreateByName string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:user_approve"`
	UpdateByName string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:user_approve"`
}

func (m *UserApproveGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type UserApproveInsertReq struct {
	Id           int    `json:"-" comment:"id"` // id
	UserId       int    `json:"userId" comment:"用户ID"`
	ApproveId    int    `json:"approveId" comment:"审批流ID"`
	Status       int    `json:"status" comment:"状态"`
	CreateByName string `json:"createByName" comment:""`
	UpdateByName string `json:"updateByName" comment:""`
	common.ControlBy
}

func (s *UserApproveInsertReq) Generate(model *models.UserApprove) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UserId = s.UserId
	model.ApproveId = s.ApproveId
	model.Status = s.Status
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
	model.UpdateByName = s.UpdateByName
}

func (s *UserApproveInsertReq) GetId() interface{} {
	return s.Id
}

type UserApproveUpdateReq struct {
	Id           int    `uri:"id" comment:"id"` // id
	UserId       int    `json:"userId" comment:"用户ID"`
	ApproveId    int    `json:"approveId" comment:"审批流ID"`
	Status       int    `json:"status" comment:"状态"`
	CreateByName string `json:"createByName" comment:""`
	UpdateByName string `json:"updateByName" comment:""`
	common.ControlBy
}

func (s *UserApproveUpdateReq) Generate(model *models.UserApprove) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UserId = s.UserId
	model.ApproveId = s.ApproveId
	model.Status = s.Status
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.CreateByName = s.CreateByName
	model.UpdateByName = s.UpdateByName
}

func (s *UserApproveUpdateReq) GetId() interface{} {
	return s.Id
}

// UserApproveGetReq 功能获取请求参数
type UserApproveGetReq struct {
	Id int `uri:"id"`
}

func (s *UserApproveGetReq) GetId() interface{} {
	return s.Id
}

// UserApproveDeleteReq 功能删除请求参数
type UserApproveDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *UserApproveDeleteReq) GetId() interface{} {
	return s.Ids
}
