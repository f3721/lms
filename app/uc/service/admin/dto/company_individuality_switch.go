package dto

import (
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type CompanyIndividualitySwitchGetPageReq struct {
	dto.Pagination `search:"-"`
	CompanyIndividualitySwitchOrder
}

type CompanyIndividualitySwitchOrder struct {
	Id           string `form:"idOrder"  search:"type:order;column:id;table:company_individuality_switch"`
	CompanyId    string `form:"companyIdOrder"  search:"type:order;column:company_id;table:company_individuality_switch"`
	Keyword      string `form:"keywordOrder"  search:"type:order;column:keyword;table:company_individuality_switch"`
	SwitchStatus string `form:"switchStatusOrder"  search:"type:order;column:switch_status;table:company_individuality_switch"`
	Status       string `form:"statusOrder"  search:"type:order;column:status;table:company_individuality_switch"`
	Desc         string `form:"descOrder"  search:"type:order;column:desc;table:company_individuality_switch"`
	Sort         string `form:"sortOrder"  search:"type:order;column:sort;table:company_individuality_switch"`
	IsDel        string `form:"isDel"  search:"type:order;column:is_del;table:company_individuality_switch"`
	Remark       string `form:"remarkOrder"  search:"type:order;column:remark;table:company_individuality_switch"`
}

func (m *CompanyIndividualitySwitchGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type CompanyIndividualitySwitchInsertReq struct {
	Id           int    `json:"-" comment:""` //
	CompanyId    int    `json:"companyId" comment:"公司ID"`
	Keyword      string `json:"keyword" comment:"关键字"`
	SwitchStatus string `json:"switchStatus" comment:"开关状态，每个keyword的状态（是否打开或者下拉选择某个状态）"`
	Status       int    `json:"status" comment:"状态 0关闭 1启用 默认1"`
	Desc         string `json:"desc" comment:"字段描述"`
	Sort         int    `json:"sort" comment:"排序 值越高 排名越前"`
	IsDel        int    `json:"isDel" comment:"是否删除 0.否 1.是  默认0"`
	Remark       string `json:"remark" comment:"备注"`
	common.ControlBy
}

func (s *CompanyIndividualitySwitchInsertReq) Generate(model *models.CompanyIndividualitySwitch) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CompanyId = s.CompanyId
	model.Keyword = s.Keyword
	model.SwitchStatus = s.SwitchStatus
	model.Status = s.Status
	model.Desc = s.Desc
	model.Sort = s.Sort
	model.IsDel = s.IsDel
	model.Remark = s.Remark
}

func (s *CompanyIndividualitySwitchInsertReq) GetId() interface{} {
	return s.Id
}

type CompanyIndividualitySwitchUpdateReq struct {
	Id           int    `uri:"id" comment:""` //
	CompanyId    int    `json:"companyId" comment:"公司ID"`
	Keyword      string `json:"keyword" comment:"关键字"`
	SwitchStatus string `json:"switchStatus" comment:"开关状态，每个keyword的状态（是否打开或者下拉选择某个状态）"`
	Status       int    `json:"status" comment:"状态 0关闭 1启用 默认1"`
	Desc         string `json:"desc" comment:"字段描述"`
	Sort         int    `json:"sort" comment:"排序 值越高 排名越前"`
	IsDel        int    `json:"isDel" comment:"是否删除 0.否 1.是  默认0"`
	Remark       string `json:"remark" comment:"备注"`
	common.ControlBy
}

func (s *CompanyIndividualitySwitchUpdateReq) Generate(model *models.CompanyIndividualitySwitch) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CompanyId = s.CompanyId
	model.Keyword = s.Keyword
	model.SwitchStatus = s.SwitchStatus
	model.Status = s.Status
	model.Desc = s.Desc
	model.Sort = s.Sort
	model.IsDel = s.IsDel
	model.Remark = s.Remark
}

func (s *CompanyIndividualitySwitchUpdateReq) GetId() interface{} {
	return s.Id
}

// CompanyIndividualitySwitchGetReq 功能获取请求参数
type CompanyIndividualitySwitchGetReq struct {
	Id int `uri:"id"`
}

func (s *CompanyIndividualitySwitchGetReq) GetId() interface{} {
	return s.Id
}

// CompanyIndividualitySwitchDeleteReq 功能删除请求参数
type CompanyIndividualitySwitchDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *CompanyIndividualitySwitchDeleteReq) GetId() interface{} {
	return s.Ids
}
