package dto

import (
	"go-admin/app/oc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CsApplyLogGetPageReq struct {
	dto.Pagination `search:"-"`
	DataId         string `form:"dataId"  search:"type:exact;column:cs_no;table:cs_apply_log" comment:"售后申请编号"`
	CsApplyLogOrder
}

type CsApplyLogOrder struct {
	Id         string `form:"idOrder"  search:"type:order;column:id;table:cs_apply_log"`
	CsNo       string `form:"csNoOrder"  search:"type:order;column:cs_no;table:cs_apply_log"`
	UserId     string `form:"userIdOrder"  search:"type:order;column:user_id;table:cs_apply_log"`
	HandlerLog string `form:"handlerLogOrder"  search:"type:order;column:handler_log;table:cs_apply_log"`
	UserName   string `form:"userNameOrder"  search:"type:order;column:user_name;table:cs_apply_log"`
}

func (m *CsApplyLogGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CsApplyLogInsertReq struct {
	CsNo       string `json:"csNo" comment:"售后申请编号(对应于cs_apply.cs_no)"`
	UserId     string `json:"userId" comment:"备注人id"`
	HandlerLog string `json:"handlerLog" comment:"处理记录"`
	UserName   string `json:"userName" comment:""`
	common.ControlBy
}

func (s *CsApplyLogInsertReq) Generate(model *models.CsApplyLog) {
	model.Model = common.Model{Id: 0}
	model.CsNo = s.CsNo
	model.HandlerLog = s.HandlerLog
}

type CsApplyLogUpdateReq struct {
	Id         int    `uri:"id" comment:""` //
	CsNo       string `json:"csNo" comment:"售后申请编号(对应于cs_apply.cs_no)"`
	UserId     string `json:"userId" comment:"备注人id"`
	HandlerLog string `json:"handlerLog" comment:"处理记录"`
	UserName   string `json:"userName" comment:""`
	common.ControlBy
}

func (s *CsApplyLogUpdateReq) Generate(model *models.CsApplyLog) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CsNo = s.CsNo
	model.UserId = s.UserId
	model.HandlerLog = s.HandlerLog
	model.UserName = s.UserName
}

func (s *CsApplyLogUpdateReq) GetId() interface{} {
	return s.Id
}

// CsApplyLogGetReq 功能获取请求参数
type CsApplyLogGetReq struct {
	Id int `uri:"id"`
}

func (s *CsApplyLogGetReq) GetId() interface{} {
	return s.Id
}

// CsApplyLogDeleteReq 功能删除请求参数
type CsApplyLogDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *CsApplyLogDeleteReq) GetId() interface{} {
	return s.Ids
}
