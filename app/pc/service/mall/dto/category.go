package dto

import (
	"go-admin/app/pc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CategoryGetPageReq struct {
	dto.Pagination `search:"-"`
	Id             int    `form:"id"  search:"type:exact;column:id;table:category"`
	NameZh         string `form:"nameZh"  search:"type:contains;column:name_zh;table:category"`
	Status         string `form:"status"  search:"type:exact;column:status;table:category"`
	ParentId       string `form:"parentId"  search:"type:exact;column:parent_id;table:category"`
	CateLevel      int    `form:"cateLevel"  search:"type:exact;column:cate_level;table:category"`
	CategoryOrder
}

type CategoryOrder struct {
	Id              string `form:"idOrder"  search:"type:order;column:id;table:category"`
	CateLevel       string `form:"cateLevelOrder"  search:"type:order;column:cate_level;table:category"`
	Seq             string `form:"seqOrder"  search:"type:order;column:seq;table:category"`
	NameZh          string `form:"nameZhOrder"  search:"type:order;column:name_zh;table:category"`
	NameEn          string `form:"nameEnOrder"  search:"type:order;column:name_en;table:category"`
	ParentId        string `form:"parentIdOrder"  search:"type:order;column:parent_id;table:category"`
	Description     string `form:"descriptionOrder"  search:"type:order;column:description;table:category"`
	Status          string `form:"statusOrder"  search:"type:order;column:status;table:category"`
	KeyWords        string `form:"keyWordsOrder"  search:"type:order;column:key_words;table:category"`
	Tax             string `form:"taxOrder"  search:"type:order;column:tax;table:category"`
	CategoryTaxCode string `form:"categoryTaxCodeOrder"  search:"type:order;column:category_tax_code;table:category"`
	CreatedAt       string `form:"createdAtOrder"  search:"type:order;column:created_at;table:category"`
	UpdatedAt       string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:category"`
	DeletedAt       string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:category"`
	CreateBy        string `form:"createByOrder"  search:"type:order;column:create_by;table:category"`
	CreateByName    string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:category"`
	UpdateBy        string `form:"updateByOrder"  search:"type:order;column:update_by;table:category"`
	UpdateByName    string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:category"`
}

func (m *CategoryGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CategoryUpdateReq struct {
	Id              int    `uri:"id" comment:"id"` // id
	CateLevel       int    `json:"cateLevel" comment:"层级"`
	Seq             int    `json:"seq" comment:"序列"`
	NameZh          string `json:"nameZh" comment:"中文名"`
	NameEn          string `json:"nameEn" comment:"英文名"`
	ParentId        int    `json:"parentId" comment:"父类id"`
	Description     string `json:"description" comment:"描述"`
	Status          int    `json:"status" comment:"产线状态"`
	KeyWords        string `json:"keyWords" comment:"关键字"`
	Tax             string `json:"tax" comment:"产线税率：默认空 值有（0.13,0.06.0.09）"`
	CategoryTaxCode string `json:"categoryTaxCode" comment:"产线税号"`
	common.ControlBy
}

func (s *CategoryUpdateReq) Generate(model *models.Category) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CateLevel = s.CateLevel
	model.Seq = s.Seq
	model.NameZh = s.NameZh
	model.NameEn = s.NameEn
	model.ParentId = s.ParentId
	model.Description = s.Description
	model.Status = s.Status
	model.KeyWords = s.KeyWords
	model.Tax = s.Tax
	model.CategoryTaxCode = s.CategoryTaxCode
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *CategoryUpdateReq) GetId() interface{} {
	return s.Id
}

// CategoryGetReq 功能获取请求参数
type CategoryGetReq struct {
	Id int `uri:"id"`
}

func (s *CategoryGetReq) GetId() interface{} {
	return s.Id
}

type CategoryList struct {
	Id1     int
	Id2     int
	Id3     int
	Id4     int
	NameZh1 string
	NameZh2 string
	NameZh3 string
	NameZh4 string
}

type CategoryNav struct {
	CategoryId int    `json:"categoryId"`
	NameZh     string `json:"nameZh"`
	Selected   int    `json:"selected"`
}

type CategoryGetPageResp struct {
	models.Category
	Addchild bool `gorm:"-" json:"addchild"`
	Haschild bool `json:"haschild"`
}
