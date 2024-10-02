package dto

import (
	"go-admin/app/pc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type BrandLogGetPageReq struct {
	dto.Pagination `search:"-"`
	BrandId        int `form:"brandId"  search:"type:exact;column:data_id;table:brand_log"`
	BrandLogOrder
}

type BrandLogOrder struct {
	Id           string `form:"idOrder"  search:"type:order;column:id;table:brand_log"`
	DataId       string `form:"dataIdOrder"  search:"type:order;column:data_id;table:brand_log"`
	Type         string `form:"typeOrder"  search:"type:order;column:type;table:brand_log"`
	Data         string `form:"dataOrder"  search:"type:order;column:data;table:brand_log"`
	BeforeData   string `form:"beforeDataOrder"  search:"type:order;column:before_data;table:brand_log"`
	AfterData    string `form:"afterDataOrder"  search:"type:order;column:after_data;table:brand_log"`
	CreatedAt    string `form:"createdAtOrder"  search:"type:order;column:created_at;table:brand_log"`
	UpdatedAt    string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:brand_log"`
	DeletedAt    string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:brand_log"`
	CreateBy     string `form:"createByOrder"  search:"type:order;column:create_by;table:brand_log"`
	CreateByName string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:brand_log"`
	UpdateBy     string `form:"updateByOrder"  search:"type:order;column:update_by;table:brand_log"`
	UpdateByName string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:brand_log"`
}

func (m *BrandLogGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type BrandLogInsertReq struct {
	Id         int    `json:"-" comment:""` //
	DataId     int    `json:"dataId" comment:"关联的主键ID"`
	Type       string `json:"type" comment:"变更类型"`
	Data       string `json:"data" comment:"源数据"`
	BeforeData string `json:"beforeData" comment:"变更前数据"`
	AfterData  string `json:"afterData" comment:"变更后数据"`
	common.ControlBy
}

func (s *BrandLogInsertReq) Generate(model *models.BrandLog) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.DataId = s.DataId
	model.Type = s.Type
	model.Data = s.Data
	model.BeforeData = s.BeforeData
	model.AfterData = s.AfterData
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *BrandLogInsertReq) GetId() interface{} {
	return s.Id
}

type BrandLogUpdateReq struct {
	Id         int    `uri:"id" comment:""` //
	DataId     int    `json:"dataId" comment:"关联的主键ID"`
	Type       string `json:"type" comment:"变更类型"`
	Data       string `json:"data" comment:"源数据"`
	BeforeData string `json:"beforeData" comment:"变更前数据"`
	AfterData  string `json:"afterData" comment:"变更后数据"`
	common.ControlBy
}

func (s *BrandLogUpdateReq) Generate(model *models.BrandLog) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.DataId = s.DataId
	model.Type = s.Type
	model.Data = s.Data
	model.BeforeData = s.BeforeData
	model.AfterData = s.AfterData
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *BrandLogUpdateReq) GetId() interface{} {
	return s.Id
}

// BrandLogGetReq 功能获取请求参数
type BrandLogGetReq struct {
	Id int `uri:"id"`
}

func (s *BrandLogGetReq) GetId() interface{} {
	return s.Id
}

// BrandLogDeleteReq 功能删除请求参数
type BrandLogDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *BrandLogDeleteReq) GetId() interface{} {
	return s.Ids
}
