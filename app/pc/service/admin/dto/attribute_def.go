package dto

import (
	"go-admin/app/pc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type AttributeDefGetPageReq struct {
	dto.Pagination `search:"-"`
	Id             int    `form:"id"  search:"type:exact;column:id;table:attribute_def"`
	NameZh         string `form:"nameZh"  search:"type:contains;column:name_zh;table:attribute_def"`
	NameEn         string `form:"nameEn"  search:"type:contains;column:name_en;table:attribute_def"`
	AttributeDefOrder
}

type AttributeDefOrder struct {
	Id           string `form:"idOrder"  search:"type:order;column:id;table:attribute_def"`
	NameZh       string `form:"nameZhOrder"  search:"type:order;column:name_zh;table:attribute_def"`
	NameEn       string `form:"nameEnOrder"  search:"type:order;column:name_en;table:attribute_def"`
	AttrType     string `form:"attrTypeOrder"  search:"type:order;column:attr_type;table:attribute_def"`
	CreatedAt    string `form:"createdAtOrder"  search:"type:order;column:created_at;table:attribute_def"`
	UpdatedAt    string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:attribute_def"`
	DeletedAt    string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:attribute_def"`
	CreateBy     string `form:"createByOrder"  search:"type:order;column:create_by;table:attribute_def"`
	CreateByName string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:attribute_def"`
	UpdateBy     string `form:"updateByOrder"  search:"type:order;column:update_by;table:attribute_def"`
	UpdateByName string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:attribute_def"`
}

func (m *AttributeDefGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type AttributeDefInsertReq struct {
	Id           int    `json:"-" comment:"属性ID"` // 属性ID
	NameZh       string `json:"nameZh" comment:"属性英文名"`
	NameEn       string `json:"nameEn" comment:"属性中文名"`
	AttrType     int    `json:"attrType" comment:"属性类型(0.主属性1:市场属性,2:技术属性)"`
	CreateByName string `json:"createByName" comment:"创建人姓名"`
	UpdateByName string `json:"updateByName" comment:"修改人姓名"`
	common.ControlBy
}

func (s *AttributeDefInsertReq) Generate(model *models.AttributeDef) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.NameZh = s.NameZh
	model.NameEn = s.NameEn
	model.AttrType = s.AttrType
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
	model.UpdateByName = s.UpdateByName
}

func (s *AttributeDefInsertReq) GetId() interface{} {
	return s.Id
}

type AttributeDefUpdateReq struct {
	Id           int    `uri:"id" comment:"属性ID"` // 属性ID
	NameZh       string `json:"nameZh" comment:"属性英文名"`
	NameEn       string `json:"nameEn" comment:"属性中文名"`
	AttrType     int    `json:"attrType" comment:"属性类型(0.主属性1:市场属性,2:技术属性)"`
	CreateByName string `json:"createByName" comment:"创建人姓名"`
	UpdateByName string `json:"updateByName" comment:"修改人姓名"`
	common.ControlBy
}

func (s *AttributeDefUpdateReq) Generate(model *models.AttributeDef) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.NameZh = s.NameZh
	model.NameEn = s.NameEn
	model.AttrType = s.AttrType
	model.CreateByName = s.CreateByName
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.UpdateByName = s.UpdateByName
}

func (s *AttributeDefUpdateReq) GetId() interface{} {
	return s.Id
}

// AttributeDefGetReq 功能获取请求参数
type AttributeDefGetReq struct {
	Id int `uri:"id"`
}

func (s *AttributeDefGetReq) GetId() interface{} {
	return s.Id
}

// AttributeDefDeleteReq 功能删除请求参数
type AttributeDefDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *AttributeDefDeleteReq) GetId() interface{} {
	return s.Ids
}
