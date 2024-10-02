package dto

import (
    "go-admin/app/pc/models"
    "go-admin/common/dto"
    common "go-admin/common/models"
)

type GoodsLogGetPageReq struct {
    dto.Pagination `search:"-"`
    GoodsId        int `form:"goodsId"  search:"type:exact;column:data_id;table:goods_log"`
    GoodsLogOrder
}

type GoodsLogOrder struct {
    Id           string `form:"idOrder"  search:"type:order;column:id;table:goods_log"`
    DataId       string `form:"dataIdOrder"  search:"type:order;column:data_id;table:goods_log"`
    Type         string `form:"typeOrder"  search:"type:order;column:type;table:goods_log"`
    Data         string `form:"dataOrder"  search:"type:order;column:data;table:goods_log"`
    BeforeData   string `form:"beforeDataOrder"  search:"type:order;column:before_data;table:goods_log"`
    AfterData    string `form:"afterDataOrder"  search:"type:order;column:after_data;table:goods_log"`
    DiffData     string `form:"diffDataOrder"  search:"type:order;column:diff_data;table:goods_log"`
    CreatedAt    string `form:"createdAtOrder"  search:"type:order;column:created_at;table:goods_log"`
    UpdatedAt    string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:goods_log"`
    DeletedAt    string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:goods_log"`
    CreateBy     string `form:"createByOrder"  search:"type:order;column:create_by;table:goods_log"`
    CreateByName string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:goods_log"`
    UpdateBy     string `form:"updateByOrder"  search:"type:order;column:update_by;table:goods_log"`
    UpdateByName string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:goods_log"`
}

func (m *GoodsLogGetPageReq) GetNeedSearch() interface{} {
    return *m
}

type GoodsLogInsertReq struct {
    Id         int    `json:"-" comment:""` //
    DataId     int    `json:"dataId" comment:"关联的主键ID"`
    Type       string `json:"type" comment:"变更类型"`
    Data       string `json:"data" comment:"源数据"`
    BeforeData string `json:"beforeData" comment:"变更前数据"`
    AfterData  string `json:"afterData" comment:"变更后数据"`
    DiffData   string `json:"diffData" comment:"差异数据"`
    common.ControlBy
}

func (s *GoodsLogInsertReq) Generate(model *models.GoodsLog) {
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

func (s *GoodsLogInsertReq) GetId() interface{} {
    return s.Id
}

type GoodsLogUpdateReq struct {
    Id         int    `uri:"id" comment:""` //
    DataId     int    `json:"dataId" comment:"关联的主键ID"`
    Type       string `json:"type" comment:"变更类型"`
    Data       string `json:"data" comment:"源数据"`
    BeforeData string `json:"beforeData" comment:"变更前数据"`
    AfterData  string `json:"afterData" comment:"变更后数据"`
    DiffData   string `json:"diffData" comment:"差异数据"`
    common.ControlBy
}

func (s *GoodsLogUpdateReq) Generate(model *models.GoodsLog) {
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

func (s *GoodsLogUpdateReq) GetId() interface{} {
    return s.Id
}

// GoodsLogGetReq 功能获取请求参数
type GoodsLogGetReq struct {
    Id int `uri:"id"`
}

func (s *GoodsLogGetReq) GetId() interface{} {
    return s.Id
}

// GoodsLogDeleteReq 功能删除请求参数
type GoodsLogDeleteReq struct {
    Ids []int `json:"ids"`
}

func (s *GoodsLogDeleteReq) GetId() interface{} {
    return s.Ids
}
