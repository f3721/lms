package dto

import (
    "go-admin/app/pc/models"
    "go-admin/common/dto"
    common "go-admin/common/models"
)

type ProductLogGetPageReq struct {
    dto.Pagination `search:"-"`
    ProductId      int `form:"productId"  search:"type:exact;column:data_id;table:product_log"`
    ProductLogOrder
}

type ProductLogOrder struct {
    Id           string `form:"idOrder"  search:"type:order;column:id;table:product_log"`
    DataId       string `form:"dataIdOrder"  search:"type:order;column:data_id;table:product_log"`
    Type         string `form:"typeOrder"  search:"type:order;column:type;table:product_log"`
    Data         string `form:"dataOrder"  search:"type:order;column:data;table:product_log"`
    BeforeData   string `form:"beforeDataOrder"  search:"type:order;column:before_data;table:product_log"`
    AfterData    string `form:"afterDataOrder"  search:"type:order;column:after_data;table:product_log"`
    DiffData     string `form:"diffDataOrder"  search:"type:order;column:diff_data;table:product_log"`
    CreatedAt    string `form:"createdAtOrder"  search:"type:order;column:created_at;table:product_log"`
    UpdatedAt    string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:product_log"`
    DeletedAt    string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:product_log"`
    CreateBy     string `form:"createByOrder"  search:"type:order;column:create_by;table:product_log"`
    CreateByName string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:product_log"`
    UpdateBy     string `form:"updateByOrder"  search:"type:order;column:update_by;table:product_log"`
    UpdateByName string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:product_log"`
}

func (m *ProductLogGetPageReq) GetNeedSearch() interface{} {
    return *m
}

type ProductLogInsertReq struct {
    Id         int    `json:"-" comment:""` //
    DataId     int    `json:"dataId" comment:"关联的主键ID"`
    Type       string `json:"type" comment:"变更类型"`
    Data       string `json:"data" comment:"源数据"`
    BeforeData string `json:"beforeData" comment:"变更前数据"`
    AfterData  string `json:"afterData" comment:"变更后数据"`
    DiffData   string `json:"diffData" comment:"差异数据"`
    common.ControlBy
}

func (s *ProductLogInsertReq) Generate(model *models.ProductLog) {
    if s.Id == 0 {
        model.Model = common.Model{Id: s.Id}
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

func (s *ProductLogInsertReq) GetId() interface{} {
    return s.Id
}

type ProductLogUpdateReq struct {
    Id         int    `uri:"id" comment:""` //
    DataId     int    `json:"dataId" comment:"关联的主键ID"`
    Type       string `json:"type" comment:"变更类型"`
    Data       string `json:"data" comment:"源数据"`
    BeforeData string `json:"beforeData" comment:"变更前数据"`
    AfterData  string `json:"afterData" comment:"变更后数据"`
    DiffData   string `json:"diffData" comment:"差异数据"`
    common.ControlBy
}

func (s *ProductLogUpdateReq) Generate(model *models.ProductLog) {
    if s.Id == 0 {
        model.Model = common.Model{Id: s.Id}
    }
    model.DataId = s.DataId
    model.Type = s.Type
    model.Data = s.Data
    model.BeforeData = s.BeforeData
    model.AfterData = s.AfterData
    model.DiffData = s.DiffData
    model.UpdateBy = s.UpdateBy
    model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *ProductLogUpdateReq) GetId() interface{} {
    return s.Id
}

// ProductLogGetReq 功能获取请求参数
type ProductLogGetReq struct {
    Id int `uri:"id"`
}

func (s *ProductLogGetReq) GetId() interface{} {
    return s.Id
}

// ProductLogDeleteReq 功能删除请求参数
type ProductLogDeleteReq struct {
    Ids []int `json:"ids"`
}

func (s *ProductLogDeleteReq) GetId() interface{} {
    return s.Ids
}
