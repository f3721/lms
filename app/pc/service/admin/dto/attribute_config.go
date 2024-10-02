package dto

import (
	"go-admin/app/pc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type AttributeConfigGetPageReq struct {
	dto.Pagination `search:"-"`
	AttributeConfigOrder
}

type AttributeConfigOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:attribute_config"`
	Type      string `form:"typeOrder"  search:"type:order;column:type;table:attribute_config"`
	Key       string `form:"keyOrder"  search:"type:order;column:key;table:attribute_config"`
	Value     string `form:"valueOrder"  search:"type:order;column:value;table:attribute_config"`
	SortOrder string `form:"sortOrderOrder"  search:"type:order;column:sort_order;table:attribute_config"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:attribute_config"`
	UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:attribute_config"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:attribute_config"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:attribute_config"`
	UpdateBy  string `form:"updateByOrder"  search:"type:order;column:update_by;table:attribute_config"`
}

func (m *AttributeConfigGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type AttributeConfigInsertReq struct {
	Id        int    `json:"-" comment:""` //
	Type      string `json:"type" comment:"类型"`
	Key       string `json:"key" comment:"配置键"`
	Value     string `json:"value" comment:"配置值"`
	SortOrder int    `json:"sortOrder" comment:""`
}

func (s *AttributeConfigInsertReq) Generate(model *models.AttributeConfig) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.Type = s.Type
	model.Key = s.Key
	model.Value = s.Value
	model.SortOrder = s.SortOrder
}

func (s *AttributeConfigInsertReq) GetId() interface{} {
	return s.Id
}

type AttributeConfigUpdateReq struct {
	Id        int    `uri:"id" comment:""` //
	Type      string `json:"type" comment:"类型"`
	Key       string `json:"key" comment:"配置键"`
	Value     string `json:"value" comment:"配置值"`
	SortOrder int    `json:"sortOrder" comment:""`
}

func (s *AttributeConfigUpdateReq) Generate(model *models.AttributeConfig) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.Type = s.Type
	model.Key = s.Key
	model.Value = s.Value
	model.SortOrder = s.SortOrder
}

func (s *AttributeConfigUpdateReq) GetId() interface{} {
	return s.Id
}

// AttributeConfigGetReq 功能获取请求参数
type AttributeConfigGetReq struct {
	Id int `uri:"id"`
}

func (s *AttributeConfigGetReq) GetId() interface{} {
	return s.Id
}

// AttributeConfigDeleteReq 功能删除请求参数
type AttributeConfigDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *AttributeConfigDeleteReq) GetId() interface{} {
	return s.Ids
}
