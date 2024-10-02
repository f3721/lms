package dto

import (
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type ClassificationGetPageReq struct {
	dto.Pagination `search:"-"`
	ClassificationOrder
}

type ClassificationOrder struct {
	Id        string `form:"idOrder"  search:"type:order;column:id;table:classification"`
	CompanyId string `form:"companyIdOrder"  search:"type:order;column:company_id;table:classification"`
	Name      string `form:"nameOrder"  search:"type:order;column:name;table:classification"`
	Remark    string `form:"remarkOrder"  search:"type:order;column:remark;table:classification"`
}

func (m *ClassificationGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ClassificationInsertReq struct {
	Id        int    `json:"-" comment:""` //
	CompanyId int    `json:"companyId" comment:"公司ID"`
	Name      string `json:"name" comment:"名称"`
	Remark    string `json:"remark" comment:"备注"`
	common.ControlBy
}

func (s *ClassificationInsertReq) Generate(model *models.Classification) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CompanyId = s.CompanyId
	model.Name = s.Name
	model.Remark = s.Remark
}

func (s *ClassificationInsertReq) GetId() interface{} {
	return s.Id
}

type ClassificationUpdateReq struct {
	Id        int    `uri:"id" comment:""` //
	CompanyId int    `json:"companyId" comment:"公司ID"`
	Name      string `json:"name" comment:"名称"`
	Remark    string `json:"remark" comment:"备注"`
	common.ControlBy
}

func (s *ClassificationUpdateReq) Generate(model *models.Classification) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CompanyId = s.CompanyId
	model.Name = s.Name
	model.Remark = s.Remark
}

func (s *ClassificationUpdateReq) GetId() interface{} {
	return s.Id
}

// ClassificationGetReq 功能获取请求参数
type ClassificationGetReq struct {
	Id int `uri:"id"`
}

func (s *ClassificationGetReq) GetId() interface{} {
	return s.Id
}

// ClassificationDeleteReq 功能删除请求参数
type ClassificationDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *ClassificationDeleteReq) GetId() interface{} {
	return s.Ids
}
