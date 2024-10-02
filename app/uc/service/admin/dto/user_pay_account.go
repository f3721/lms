package dto

import (
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type UserPayAccountGetPageReq struct {
	dto.Pagination `search:"-"`
	UserPayAccountOrder
}

type UserPayAccountOrder struct {
	Id               string `form:"idOrder"  search:"type:order;column:id;table:user_pay_account"`
	CompanyId        string `form:"companyIdOrder"  search:"type:order;column:company_id;table:user_pay_account"`
	UserId           string `form:"userIdOrder"  search:"type:order;column:user_id;table:user_pay_account"`
	PayAccount       string `form:"payAccountOrder"  search:"type:order;column:pay_account;table:user_pay_account"`
	ClassificationId string `form:"classificationIdOrder"  search:"type:order;column:classification_id;table:user_pay_account"`
}

func (m *UserPayAccountGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type UserPayAccountInsertReq struct {
	Id               int    `json:"-" comment:""` //
	CompanyId        int    `json:"companyId" comment:"公司ID"`
	UserId           int    `json:"userId" comment:"用户ID"`
	PayAccount       string `json:"payAccount" comment:"支付账户"`
	ClassificationId int    `json:"classificationId" comment:"客户分类"`
	common.ControlBy
}

func (s *UserPayAccountInsertReq) Generate(model *models.UserPayAccount) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CompanyId = s.CompanyId
	model.UserId = s.UserId
	model.PayAccount = s.PayAccount
	model.ClassificationId = s.ClassificationId
}

func (s *UserPayAccountInsertReq) GetId() interface{} {
	return s.Id
}

type UserPayAccountUpdateReq struct {
	Id               int    `uri:"id" comment:""` //
	CompanyId        int    `json:"companyId" comment:"公司ID"`
	UserId           int    `json:"userId" comment:"用户ID"`
	PayAccount       string `json:"payAccount" comment:"支付账户"`
	ClassificationId int    `json:"classificationId" comment:"客户分类"`
	common.ControlBy
}

func (s *UserPayAccountUpdateReq) Generate(model *models.UserPayAccount) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.CompanyId = s.CompanyId
	model.UserId = s.UserId
	model.PayAccount = s.PayAccount
	model.ClassificationId = s.ClassificationId
}

func (s *UserPayAccountUpdateReq) GetId() interface{} {
	return s.Id
}

// UserPayAccountGetReq 功能获取请求参数
type UserPayAccountGetReq struct {
	Id int `uri:"id"`
}

func (s *UserPayAccountGetReq) GetId() interface{} {
	return s.Id
}

// UserPayAccountDeleteReq 功能删除请求参数
type UserPayAccountDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *UserPayAccountDeleteReq) GetId() interface{} {
	return s.Ids
}
