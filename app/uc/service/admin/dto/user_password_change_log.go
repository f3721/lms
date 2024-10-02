package dto

import (
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"time"
)

type UserPasswordChangeLogGetPageReq struct {
	dto.Pagination `search:"-"`
	UserId         int `form:"userId"  search:"type:exact;column:user_id;table:user_password_change_log" comment:"用户ID"`
	UserPasswordChangeLogOrder
}

type UserPasswordChangeLogOrder struct {
	Id         string `form:"idOrder"  search:"type:order;column:id;table:user_password_change_log"`
	UserId     string `form:"userIdOrder"  search:"type:order;column:user_id;table:user_password_change_log"`
	RecordTime string `form:"recordTimeOrder"  search:"type:order;column:record_time;table:user_password_change_log"`
	RecordName string `form:"recordNameOrder"  search:"type:order;column:record_name;table:user_password_change_log"`
}

func (m *UserPasswordChangeLogGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type UserPasswordChangeLogInsertReq struct {
	Id         int       `json:"-" comment:"编号"` // 编号
	UserId     int       `json:"userId" comment:"用户ID"`
	RecordTime time.Time `json:"recordTime" comment:"修改时间时间"`
	RecordName string    `json:"recordName" comment:"修改人"`
}

func (s *UserPasswordChangeLogInsertReq) Generate(model *models.UserPasswordChangeLog) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UserId = s.UserId
	model.RecordTime = s.RecordTime
	model.RecordName = s.RecordName
}

func (s *UserPasswordChangeLogInsertReq) GetId() interface{} {
	return s.Id
}

type UserPasswordChangeLogUpdateReq struct {
	Id         int       `uri:"id" comment:"编号"` // 编号
	UserId     int       `json:"userId" comment:"用户ID"`
	RecordTime time.Time `json:"recordTime" comment:"修改时间时间"`
	RecordName string    `json:"recordName" comment:"修改人"`
}

func (s *UserPasswordChangeLogUpdateReq) Generate(model *models.UserPasswordChangeLog) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UserId = s.UserId
	model.RecordTime = s.RecordTime
	model.RecordName = s.RecordName
}

func (s *UserPasswordChangeLogUpdateReq) GetId() interface{} {
	return s.Id
}

// UserPasswordChangeLogGetReq 功能获取请求参数
type UserPasswordChangeLogGetReq struct {
	Id int `uri:"id"`
}

func (s *UserPasswordChangeLogGetReq) GetId() interface{} {
	return s.Id
}

// UserPasswordChangeLogDeleteReq 功能删除请求参数
type UserPasswordChangeLogDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *UserPasswordChangeLogDeleteReq) GetId() interface{} {
	return s.Ids
}
