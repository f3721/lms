package dto

import (
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SkuClassificationGetPageReq struct {
	dto.Pagination `search:"-"`
	SkuClassificationOrder
}

type SkuClassificationOrder struct {
	Id               string `form:"idOrder"  search:"type:order;column:id;table:sku_classification"`
	CompanyId        string `form:"companyIdOrder"  search:"type:order;column:company_id;table:sku_classification"`
	SkuCode          string `form:"skuCodeOrder"  search:"type:order;column:sku_code;table:sku_classification"`
	Status           string `form:"statusOrder"  search:"type:order;column:status;table:sku_classification"`
	ClassificationId string `form:"classificationIdOrder"  search:"type:order;column:classification_id;table:sku_classification"`
}

func (m *SkuClassificationGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type SkuClassificationInsertReq struct {
	Id               int    `json:"-" comment:""` //
	CompanyId        int    `json:"companyId" comment:"公司ID"`
	SkuCode          string `json:"skuCode" comment:"产品sku"`
	Status           int    `json:"status" comment:"是否启用 0.否 1.是 默认1"`
	ClassificationId int    `json:"classificationId" comment:"客户分类"`
	common.ControlBy
}

func (s *SkuClassificationInsertReq) Generate(model *models.SkuClassification) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CompanyId = s.CompanyId
	model.SkuCode = s.SkuCode
	model.Status = s.Status
	model.ClassificationId = s.ClassificationId
}

func (s *SkuClassificationInsertReq) GetId() interface{} {
	return s.Id
}

type SkuClassificationUpdateReq struct {
	Id               int    `uri:"id" comment:""` //
	CompanyId        int    `json:"companyId" comment:"公司ID"`
	SkuCode          string `json:"skuCode" comment:"产品sku"`
	Status           int    `json:"status" comment:"是否启用 0.否 1.是 默认1"`
	ClassificationId int    `json:"classificationId" comment:"客户分类"`
	common.ControlBy
}

func (s *SkuClassificationUpdateReq) Generate(model *models.SkuClassification) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CompanyId = s.CompanyId
	model.SkuCode = s.SkuCode
	model.Status = s.Status
	model.ClassificationId = s.ClassificationId
}

func (s *SkuClassificationUpdateReq) GetId() interface{} {
	return s.Id
}

// SkuClassificationGetReq 功能获取请求参数
type SkuClassificationGetReq struct {
	Id int `uri:"id"`
}

func (s *SkuClassificationGetReq) GetId() interface{} {
	return s.Id
}

// SkuClassificationDeleteReq 功能删除请求参数
type SkuClassificationDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *SkuClassificationDeleteReq) GetId() interface{} {
	return s.Ids
}
