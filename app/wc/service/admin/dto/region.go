package dto

import (
	"go-admin/app/wc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type RegionGetPageReq struct {
	dto.Pagination `search:"-"`
	ParentId       string `form:"parentId"  search:"type:exact;column:parent_id;table:region" comment:""`
	Name           string `form:"name"  search:"type:contains;column:name;table:region" comment:""`
	Level          int    `form:"level"  search:"type:exact;column:level;table:region" comment:""`
}

func (m *RegionGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type RegionInsertReq struct {
	Id         int    `json:"-" comment:""` //
	ParentId   string `json:"parentId" comment:""`
	Name       string `json:"name" comment:""`
	Level      string `json:"level" comment:""`
	PostalCode string `json:"postalCode" comment:""`
	Latitude   string `json:"latitude" comment:"纬度"`
	Longitude  string `json:"longitude" comment:"经度"`
	Adcode     string `json:"adcode" comment:""`
	common.ControlBy
}

func (s *RegionInsertReq) Generate(model *models.Region) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ParentId = s.ParentId
	model.Name = s.Name
	model.Level = s.Level
	model.PostalCode = s.PostalCode
	model.Latitude = s.Latitude
	model.Longitude = s.Longitude
	model.Adcode = s.Adcode
}

func (s *RegionInsertReq) GetId() interface{} {
	return s.Id
}

type RegionUpdateReq struct {
	Id         int    `uri:"id" comment:""` //
	ParentId   string `json:"parentId" comment:""`
	Name       string `json:"name" comment:""`
	Level      string `json:"level" comment:""`
	PostalCode string `json:"postalCode" comment:""`
	Latitude   string `json:"latitude" comment:"纬度"`
	Longitude  string `json:"longitude" comment:"经度"`
	Adcode     string `json:"adcode" comment:""`
	common.ControlBy
}

func (s *RegionUpdateReq) Generate(model *models.Region) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.ParentId = s.ParentId
	model.Name = s.Name
	model.Level = s.Level
	model.PostalCode = s.PostalCode
	model.Latitude = s.Latitude
	model.Longitude = s.Longitude
	model.Adcode = s.Adcode
}

func (s *RegionUpdateReq) GetId() interface{} {
	return s.Id
}

// RegionGetReq 功能获取请求参数
type RegionGetReq struct {
	Id int `uri:"id"`
}

func (s *RegionGetReq) GetId() interface{} {
	return s.Id
}

// RegionDeleteReq 功能删除请求参数
type RegionDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *RegionDeleteReq) GetId() interface{} {
	return s.Ids
}

type InnerRegionGetByIdsReq struct {
	Ids string `form:"ids"`
}
