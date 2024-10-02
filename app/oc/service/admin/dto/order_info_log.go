package dto

import (

	"go-admin/app/oc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type OrderInfoLogGetPageReq struct {
	dto.Pagination     `search:"-"`
    DataId  int         `form:"dataId" search:"type:exact;column:data_id;table:order_info_log" comment:"订单ID"`
    OrderInfoLogOrder
}

type OrderInfoLogOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:order_info_log"`
    DataId string `form:"dataIdOrder"  search:"type:order;column:data_id;table:order_info_log"`
    //Type string `form:"typeOrder"  search:"type:order;column:type;table:order_info_log"`
    //Data string `form:"dataOrder"  search:"type:order;column:data;table:order_info_log"`
    //BeforeData string `form:"beforeDataOrder"  search:"type:order;column:before_data;table:order_info_log"`
    //AfterData string `form:"afterDataOrder"  search:"type:order;column:after_data;table:order_info_log"`
    //DiffData string `form:"diffDataOrder"  search:"type:order;column:diff_data;table:order_info_log"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:order_info_log"`
    //CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:order_info_log"`
    //CreateByName string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:order_info_log"`
    
}

func (m *OrderInfoLogGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type OrderInfoLogInsertReq struct {
    Id int `json:"-" comment:""` // 
    DataId int `json:"dataId" comment:"关联的主键ID"`
    Type string `json:"type" comment:"变更类型"`
    Data string `json:"data" comment:"源数据"`
    BeforeData string `json:"beforeData" comment:"变更前数据"`
    AfterData string `json:"afterData" comment:"变更后数据"`
    DiffData string `json:"diffData" comment:"差异数据"`
    common.ControlBy
}

func (s *OrderInfoLogInsertReq) Generate(model *models.OrderInfoLog)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.DataId = s.DataId
    model.Type = s.Type
    model.Data = s.Data
    model.BeforeData = s.BeforeData
    model.AfterData = s.AfterData
    model.DiffData = s.DiffData
    model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
    model.CreateByName = s.CreateByName
}

func (s *OrderInfoLogInsertReq) GetId() interface{} {
	return s.Id
}

type OrderInfoLogUpdateReq struct {
    Id int `uri:"id" comment:""` // 
    DataId int `json:"dataId" comment:"关联的主键ID"`
    Type string `json:"type" comment:"变更类型"`
    Data string `json:"data" comment:"源数据"`
    BeforeData string `json:"beforeData" comment:"变更前数据"`
    AfterData string `json:"afterData" comment:"变更后数据"`
    DiffData string `json:"diffData" comment:"差异数据"`
    common.ControlBy
}

func (s *OrderInfoLogUpdateReq) Generate(model *models.OrderInfoLog)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.DataId = s.DataId
    model.Type = s.Type
    model.Data = s.Data
    model.BeforeData = s.BeforeData
    model.AfterData = s.AfterData
    model.DiffData = s.DiffData
}

func (s *OrderInfoLogUpdateReq) GetId() interface{} {
	return s.Id
}

// OrderInfoLogGetReq 功能获取请求参数
type OrderInfoLogGetReq struct {
     Id int `uri:"id"`
}
func (s *OrderInfoLogGetReq) GetId() interface{} {
	return s.Id
}

// OrderInfoLogDeleteReq 功能删除请求参数
type OrderInfoLogDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *OrderInfoLogDeleteReq) GetId() interface{} {
	return s.Ids
}
