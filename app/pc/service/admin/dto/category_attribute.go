package dto

import (
	"go-admin/app/pc/models"
	common "go-admin/common/models"
)

type CategoryAttributeGetPageReq struct {
	CategoryId int `form:"categoryId"  search:"type:exact;column:category_id;table:category_attribute"`
	CategoryAttributeOrder
}

type CategoryAttributeOrder struct {
	Id           string `form:"idOrder"  search:"type:order;column:id;table:category_attribute"`
	CategoryId   string `form:"categoryIdOrder"  search:"type:order;column:category_id;table:category_attribute"`
	AttributeId  string `form:"attributeIdOrder"  search:"type:order;column:attribute_id;table:category_attribute"`
	Seq          string `form:"seqOrder"  search:"type:order;column:seq;table:category_attribute"`
	RequiredFlag string `form:"requiredFlagOrder"  search:"type:order;column:required_flag;table:category_attribute"`
	FilterFlag   string `form:"filterFlagOrder"  search:"type:order;column:filter_flag;table:category_attribute"`
	RangeVal     string `form:"rangeValOrder"  search:"type:order;column:range_val;table:category_attribute"`
	CreatedAt    string `form:"createdAtOrder"  search:"type:order;column:created_at;table:category_attribute"`
	UpdatedAt    string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:category_attribute"`
	DeletedAt    string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:category_attribute"`
	CreateBy     string `form:"createByOrder"  search:"type:order;column:create_by;table:category_attribute"`
	CreateByName string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:category_attribute"`
	UpdateBy     string `form:"updateByOrder"  search:"type:order;column:update_by;table:category_attribute"`
	UpdateByName string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:category_attribute"`
}

func (m *CategoryAttributeGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CategoryAttributeInsertReq struct {
	Id           int    `json:"-" comment:"id"` // id
	CategoryId   int    `json:"categoryId" comment:"产线id"`
	AttributeId  int    `json:"attributeId" comment:"属性id"`
	Seq          int    `json:"seq" comment:"序列"`
	RequiredFlag int    `json:"requiredFlag" comment:"必填标志"`
	FilterFlag   int    `json:"filterFlag" comment:"筛选标志(0不筛选1值筛选2范围筛选)"`
	RangeVal     string `json:"rangeVal" comment:"范围值"`
	common.ControlBy
}

func (s *CategoryAttributeInsertReq) Generate(model *models.CategoryAttribute) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CategoryId = s.CategoryId
	model.AttributeId = s.AttributeId
	model.Seq = s.Seq
	model.RequiredFlag = s.RequiredFlag
	model.FilterFlag = s.FilterFlag
	model.RangeVal = s.RangeVal
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *CategoryAttributeInsertReq) GetId() interface{} {
	return s.Id
}

type CategoryAttributeUpdateReq struct {
	Id           int    `uri:"id" comment:"id"` // id
	CategoryId   int    `json:"categoryId" comment:"产线id"`
	AttributeId  int    `json:"attributeId" comment:"属性id"`
	Seq          int    `json:"seq" comment:"序列"`
	RequiredFlag int    `json:"requiredFlag" comment:"必填标志"`
	FilterFlag   int    `json:"filterFlag" comment:"筛选标志(0不筛选1值筛选2范围筛选)"`
	RangeVal     string `json:"rangeVal" comment:"范围值"`
	common.ControlBy
}

func (s *CategoryAttributeUpdateReq) Generate(model *models.CategoryAttribute) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CategoryId = s.CategoryId
	model.AttributeId = s.AttributeId
	model.Seq = s.Seq
	model.RequiredFlag = s.RequiredFlag
	model.FilterFlag = s.FilterFlag
	model.RangeVal = s.RangeVal
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *CategoryAttributeUpdateReq) GetId() interface{} {
	return s.Id
}

// CategoryAttributeGetReq 功能获取请求参数
type CategoryAttributeGetReq struct {
	Id int `uri:"id"`
}

func (s *CategoryAttributeGetReq) GetId() interface{} {
	return s.Id
}

// CategoryAttributeDeleteReq 功能删除请求参数
type CategoryAttributeDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *CategoryAttributeDeleteReq) GetId() interface{} {
	return s.Ids
}

type CategoryAttributeExportReq struct {
	Level1CatId int `json:"level1CatId" comment:"一级目录" vd:"@:len($)>0; msg:'一级目录必填'"`
	Level2CatId int `json:"level2CatId" comment:"二级目录" vd:"@:len($)>0; msg:'二级目录必填'"`
	Level3CatId int `json:"level3CatId" comment:"三级目录" vd:"@:len($)>0; msg:'三级目录必填'"`
	Level4CatId int `json:"level4CatId" comment:"四级目录" vd:"@:len($)>0; msg:'四级目录必填'"`
}

type AttrsKeyName struct {
	KeyName int    `json:"key"`
	Name    string `json:"name"`
	Value   string `json:"value" gorm:"-"`
}

type CategoryAttributeWhere struct {
	Id          int
	CategoryId  int `json:"categoryId" comment:"产线id"`
	AttributeId int `json:"attributeId" comment:"属性id"`
}
