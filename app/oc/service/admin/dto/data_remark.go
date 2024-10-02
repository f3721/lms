package dto

import (
	"go-admin/app/oc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type DataRemarkGetPageReq struct {
	dto.Pagination `search:"-"`
	Type           string `form:"type"  search:"type:exact;column:type;table:data_remark" comment:"类型{sale_order:销售售后，purchase_order:采购售后}"` //类型{sale_order:销售售后，purchase_order:采购售后}
	DataId         string `form:"dataId"  search:"type:exact;column:data_id;table:data_remark" comment:""`                                   //订单号
	DataRemarkOrder
}

type DataRemarkOrder struct {
	Id           string `form:"idOrder"  search:"type:order;column:id;table:data_remark"`
	Type         string `form:"typeOrder"  search:"type:order;column:type;table:data_remark"`
	DataId       string `form:"dataIdOrder"  search:"type:order;column:data_id;table:data_remark"`
	Remark       string `form:"remarkOrder"  search:"type:order;column:remark;table:data_remark"`
	Usertype     string `form:"usertypeOrder"  search:"type:order;column:usertype;table:data_remark"`
	CreatedAt    string `form:"createdAtOrder"  search:"type:order;column:created_at;table:data_remark"`
	UpdatedAt    string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:data_remark"`
	DeletedAt    string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:data_remark"`
	CreateBy     string `form:"createByOrder"  search:"type:order;column:create_by;table:data_remark"`
	CreateByName string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:data_remark"`
	UpdateBy     string `form:"updateByOrder"  search:"type:order;column:update_by;table:data_remark"`
	UpdateByName string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:data_remark"`
}

func (m *DataRemarkGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type DataRemarkInsertReq struct {
	Id       int    `json:"-" comment:""`                                        //
	Type     string `json:"-" comment:"类型{sale_order:销售售后，purchase_order:采购售后}"` // 类型{sale_order:销售售后，purchase_order:采购售后}
	DataId   string `json:"dataId" comment:""`                                   //订单号
	Remark   string `json:"remark" comment:"备注"`                                 // 备注
	Usertype int    `json:"-" comment:"操作人员类型(0:opc,1:spc,2:other)"`             // 操作人员类型(0:opc,1:spc,2:other)
	common.ControlBy
}

func (s *DataRemarkInsertReq) Generate(model *models.DataRemark) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.Type = s.Type
	model.DataId = s.DataId
	model.Remark = s.Remark
	model.Usertype = s.Usertype
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *DataRemarkInsertReq) GetId() interface{} {
	return s.Id
}

type DataRemarkUpdateReq struct {
	Id       int    `uri:"id" comment:""`                                           //
	Type     string `json:"type" comment:"类型{sale_order:销售售后，purchase_order:采购售后}"` // 类型{sale_order:销售售后，purchase_order:采购售后}
	DataId   string `json:"dataId" comment:""`                                      //
	Remark   string `json:"remark" comment:"备注"`                                    // 备注
	Usertype int    `json:"usertype" comment:"操作人员类型(0:opc,1:spc,2:other)"`         // 操作人员类型(0:opc,1:spc,2:other)
	common.ControlBy
}

func (s *DataRemarkUpdateReq) Generate(model *models.DataRemark) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.Type = s.Type
	model.DataId = s.DataId
	model.Remark = s.Remark
	model.Usertype = s.Usertype
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *DataRemarkUpdateReq) GetId() interface{} {
	return s.Id
}

// DataRemarkGetReq 功能获取请求参数
type DataRemarkGetReq struct {
	Id int `uri:"id"`
}

func (s *DataRemarkGetReq) GetId() interface{} {
	return s.Id
}

// DataRemarkDeleteReq 功能删除请求参数
type DataRemarkDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *DataRemarkDeleteReq) GetId() interface{} {
	return s.Ids
}
