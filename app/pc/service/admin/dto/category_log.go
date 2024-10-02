package dto

import (
    "go-admin/app/pc/models"
    "go-admin/common/dto"
    common "go-admin/common/models"
)

type CategoryLogGetPageReq struct {
    dto.Pagination `search:"-"`
    CategoryId     int `form:"categoryId"  search:"type:exact;column:data_id;table:category_log"`
    CategoryLogOrder
}

type CategoryLogOrder struct {
    Id           string `form:"idOrder"  search:"type:order;column:id;table:category_log"`
    DataId       string `form:"dataIdOrder"  search:"type:order;column:data_id;table:category_log"`
    Type         string `form:"typeOrder"  search:"type:order;column:type;table:category_log"`
    Data         string `form:"dataOrder"  search:"type:order;column:data;table:category_log"`
    BeforeData   string `form:"beforeDataOrder"  search:"type:order;column:before_data;table:category_log"`
    AfterData    string `form:"afterDataOrder"  search:"type:order;column:after_data;table:category_log"`
    DiffData     string `form:"diffDataOrder"  search:"type:order;column:diff_data;table:category_log"`
    CreatedAt    string `form:"createdAtOrder"  search:"type:order;column:created_at;table:category_log"`
    UpdatedAt    string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:category_log"`
    DeletedAt    string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:category_log"`
    CreateBy     string `form:"createByOrder"  search:"type:order;column:create_by;table:category_log"`
    CreateByName string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:category_log"`
    UpdateBy     string `form:"updateByOrder"  search:"type:order;column:update_by;table:category_log"`
    UpdateByName string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:category_log"`
}

func (m *CategoryLogGetPageReq) GetNeedSearch() interface{} {
    return *m
}

type CategoryLogInsertReq struct {
    Id         int    `json:"-" comment:""` //
    DataId     int    `json:"dataId" comment:"关联的主键ID"`
    Type       string `json:"type" comment:"变更类型"`
    Data       string `json:"data" comment:"源数据"`
    BeforeData string `json:"beforeData" comment:"变更前数据"`
    AfterData  string `json:"afterData" comment:"变更后数据"`
    DiffData   string `json:"diffData" comment:"差异数据"`
    common.ControlBy
}

func (s *CategoryLogInsertReq) Generate(model *models.CategoryLog) {
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

func (s *CategoryLogInsertReq) GetId() interface{} {
    return s.Id
}

type CategoryLogUpdateReq struct {
    Id         int    `uri:"id" comment:""` //
    DataId     int    `json:"dataId" comment:"关联的主键ID"`
    Type       string `json:"type" comment:"变更类型"`
    Data       string `json:"data" comment:"源数据"`
    BeforeData string `json:"beforeData" comment:"变更前数据"`
    AfterData  string `json:"afterData" comment:"变更后数据"`
    DiffData   string `json:"diffData" comment:"差异数据"`
    common.ControlBy
}

func (s *CategoryLogUpdateReq) Generate(model *models.CategoryLog) {
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

func (s *CategoryLogUpdateReq) GetId() interface{} {
    return s.Id
}

// CategoryLogGetReq 功能获取请求参数
type CategoryLogGetReq struct {
    Id int `uri:"id"`
}

func (s *CategoryLogGetReq) GetId() interface{} {
    return s.Id
}

// CategoryLogDeleteReq 功能删除请求参数
type CategoryLogDeleteReq struct {
    Ids []int `json:"ids"`
}

func (s *CategoryLogDeleteReq) GetId() interface{} {
    return s.Ids
}
