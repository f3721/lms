package dto

import (
	"go-admin/app/pc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"strconv"
)

type ProductExtAttributeGetPageReq struct {
	dto.Pagination `search:"-"`
	ProductExtAttributeOrder
}

type ProductExtAttributeOrder struct {
	Id           string `form:"idOrder"  search:"type:order;column:id;table:product_ext_attribute"`
	SkuCode      string `form:"skuCodeOrder"  search:"type:order;column:sku_code;table:product_ext_attribute"`
	AttributeId  string `form:"attributeIdOrder"  search:"type:order;column:attribute_id;table:product_ext_attribute"`
	ValueZh      string `form:"valueZhOrder"  search:"type:order;column:value_zh;table:product_ext_attribute"`
	ValueEn      string `form:"valueEnOrder"  search:"type:order;column:value_en;table:product_ext_attribute"`
	Status       string `form:"statusOrder"  search:"type:order;column:status;table:product_ext_attribute"`
	CreatedAt    string `form:"createdAtOrder"  search:"type:order;column:created_at;table:product_ext_attribute"`
	UpdatedAt    string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:product_ext_attribute"`
	DeletedAt    string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:product_ext_attribute"`
	CreateBy     string `form:"createByOrder"  search:"type:order;column:create_by;table:product_ext_attribute"`
	CreateByName string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:product_ext_attribute"`
	UpdateBy     string `form:"updateByOrder"  search:"type:order;column:update_by;table:product_ext_attribute"`
	UpdateByName string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:product_ext_attribute"`
}

func (m *ProductExtAttributeGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ProductExtAttributeInsertReq struct {
	Id          int    `json:"-" comment:"ID"` // ID
	SkuCode     string `json:"skuCode" comment:"产品SKU"`
	AttributeId int    `json:"attributeId" comment:"属性ID"`
	ValueZh     string `json:"valueZh" comment:"属性值(中文)"`
	ValueEn     string `json:"valueEn" comment:"属性值(英文)"`
	Status      int    `json:"status" comment:"维护状态"`
	common.ControlBy
}

func (s *ProductExtAttributeInsertReq) Generate(model *models.ProductExtAttribute) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.SkuCode = s.SkuCode
	model.AttributeId = s.AttributeId
	model.ValueZh = s.ValueZh
	model.ValueEn = s.ValueEn
	model.Status = s.Status
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *ProductExtAttributeInsertReq) GetId() interface{} {
	return s.Id
}

type ProductExtAttributeUpdateReq struct {
	Id          int    `uri:"id" comment:"ID"` // ID
	SkuCode     string `json:"skuCode" comment:"产品SKU"`
	AttributeId int    `json:"attributeId" comment:"属性ID"`
	ValueZh     string `json:"valueZh" comment:"属性值(中文)"`
	ValueEn     string `json:"valueEn" comment:"属性值(英文)"`
	Status      int    `json:"status" comment:"维护状态"`
	common.ControlBy
}

func (s *ProductExtAttributeUpdateReq) Generate(model *models.ProductExtAttribute) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.SkuCode = s.SkuCode
	model.AttributeId = s.AttributeId
	model.ValueZh = s.ValueZh
	model.ValueEn = s.ValueEn
	model.Status = s.Status
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *ProductExtAttributeUpdateReq) GetId() interface{} {
	return s.Id
}

// ProductExtAttributeGetReq 功能获取请求参数
type ProductExtAttributeGetReq struct {
	Id int `uri:"id"`
}

func (s *ProductExtAttributeGetReq) GetId() interface{} {
	return s.Id
}

// ProductExtAttributeDeleteReq 功能删除请求参数
type ProductExtAttributeDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ProductExtAttributeDeleteReq) GetId() interface{} {
	return s.Ids
}

type ProductAttributeImportTemp struct {
	Key            int
	CategoryId     int
	SkuCode        string
	AttributeId    string
	AttributeValue string
	common.ControlBy
}

func (s *ProductAttributeImportTemp) Generate(model *models.ProductExtAttribute) {
	attributeId, _ := strconv.Atoi(s.AttributeId)
	model.SkuCode = s.SkuCode
	model.AttributeId = attributeId
	model.ValueZh = s.AttributeValue
	model.ValueEn = ""
	model.Status = 1
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

type GetProductExtAttributeReq struct {
	SkuCode    string         `json:"skuCode"`
	AttrList   []AttrsKeyName `json:"attrList"`
	CategoryId int            `json:"categoryId" vd:"@:$>0; msg:'分类ID必填！'"`
	Status     int            `json:"status"`
}

type ProductExtAttribute struct {
	KeyName int    `json:"key"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}
