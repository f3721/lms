package dto

import (
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CategoryPathGetPageReq struct {
	dto.Pagination `search:"-"`
	CategoryPathOrder
}

type CategoryPathOrder struct {
	CategoryId string `form:"categoryIdOrder"  search:"type:order;column:category_id;table:category_path"`
	PathId     string `form:"pathIdOrder"  search:"type:order;column:path_id;table:category_path"`
	Level      string `form:"levelOrder"  search:"type:order;column:level;table:category_path"`
}

func (m *CategoryPathGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CategoryPathInsertReq struct {
	CategoryId int `json:"-" comment:""` //
	PathId     int `json:"-" comment:""` //
	Level      int `json:"level" comment:""`
	common.ControlBy
}

func (s *CategoryPathInsertReq) GetId() interface{} {
	return s.PathId
}

type CategoryPathUpdateReq struct {
	CategoryId int `uri:"categoryId" comment:""` //
	PathId     int `uri:"pathId" comment:""`     //
	Level      int `json:"level" comment:""`
	common.ControlBy
}

func (s *CategoryPathUpdateReq) GetId() interface{} {
	return s.PathId
}

// CategoryPathGetReq 功能获取请求参数
type CategoryPathGetReq struct {
	CategoryId int `uri:"categoryId"`
	PathId     int `uri:"pathId"`
}

func (s *CategoryPathGetReq) GetId() interface{} {
	return s.PathId
}

// CategoryPathDeleteReq 功能删除请求参数
type CategoryPathDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *CategoryPathDeleteReq) GetId() interface{} {
	return s.Ids
}
