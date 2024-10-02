package dto

import (
	"go-admin/app/pc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type UommasterGetPageReq struct {
	dto.Pagination `search:"-"`
	Uom            string `form:"uom"  search:"type:contains;column:uom;table:uommaster"`
}

type UommasterOrder struct {
	UomId       string `form:"uomIdOrder"  search:"type:order;column:uom_id;table:uommaster"`
	Uom         string `form:"uomOrder"  search:"type:order;column:uom;table:uommaster"`
	Description string `form:"descriptionOrder"  search:"type:order;column:description;table:uommaster"`
	UomPy       string `form:"uomPyOrder"  search:"type:order;column:uom_py;table:uommaster"`
}

func (m *UommasterGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type UommasterInsertReq struct {
	UomId       int    `json:"-" comment:""` //
	Uom         string `json:"uom" comment:""`
	Description string `json:"description" comment:""`
	UomPy       string `json:"uomPy" comment:""`
	common.ControlBy
}

func (s *UommasterInsertReq) Generate(model *models.Uommaster) {
	//if s.UomId == 0 {
	//    model.Model = common.Model{ UomId: s.UomId }
	//}
	model.Uom = s.Uom
	model.Description = s.Description
	model.UomPy = s.UomPy
}

func (s *UommasterInsertReq) GetId() interface{} {
	return s.UomId
}

type UommasterUpdateReq struct {
	UomId       int    `uri:"uomId" comment:""` //
	Uom         string `json:"uom" comment:""`
	Description string `json:"description" comment:""`
	UomPy       string `json:"uomPy" comment:""`
	common.ControlBy
}

func (s *UommasterUpdateReq) Generate(model *models.Uommaster) {
	//if s.UomId == 0 {
	//    model.Model = common.Model{ UomId: s.UomId }
	//}
	model.Uom = s.Uom
	model.Description = s.Description
	model.UomPy = s.UomPy
}

func (s *UommasterUpdateReq) GetId() interface{} {
	return s.UomId
}

// UommasterGetReq 功能获取请求参数
type UommasterGetReq struct {
	UomId int `uri:"uomId"`
}

func (s *UommasterGetReq) GetId() interface{} {
	return s.UomId
}

// UommasterDeleteReq 功能删除请求参数
type UommasterDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *UommasterDeleteReq) GetId() interface{} {
	return s.Ids
}
