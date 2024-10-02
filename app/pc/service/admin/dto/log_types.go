package dto

import (

	"go-admin/app/pc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type LogTypesGetPageReq struct {
	dto.Pagination     `search:"-"`
    LogTypesOrder
}

type LogTypesOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:log_types"`
    ModelName string `form:"model_nameOrder"  search:"type:order;column:model_name;table:log_types"`
    Type string `form:"typeOrder"  search:"type:order;column:type;table:log_types"`
    Mapping string `form:"mappingOrder"  search:"type:order;column:mapping;table:log_types"`
    Title string `form:"titleOrder"  search:"type:order;column:title;table:log_types"`
    
}

func (m *LogTypesGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type LogTypesInsertReq struct {
    Id int `json:"-" comment:""` // 
    ModelName string `json:"model_name" comment:""`
    Type string `json:"type" comment:""`
    Mapping string `json:"mapping" comment:"字段映射"`
    Title string `json:"title" comment:"操作名称"`
    common.ControlBy
}

func (s *LogTypesInsertReq) Generate(model *models.LogTypes)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.ModelName = s.ModelName
    model.Type = s.Type
    model.Mapping = s.Mapping
    model.Title = s.Title
}

func (s *LogTypesInsertReq) GetId() interface{} {
	return s.Id
}

type LogTypesUpdateReq struct {
    Id int `uri:"id" comment:""` // 
    ModelName string `json:"model_name" comment:""`
    Type string `json:"type" comment:""`
    Mapping string `json:"mapping" comment:"字段映射"`
    Title string `json:"title" comment:"操作名称"`
    common.ControlBy
}

func (s *LogTypesUpdateReq) Generate(model *models.LogTypes)  {
    if s.Id == 0 {
        model.Model = common.Model{ Id: s.Id }
    }
    model.ModelName = s.ModelName
    model.Type = s.Type
    model.Mapping = s.Mapping
    model.Title = s.Title
}

func (s *LogTypesUpdateReq) GetId() interface{} {
	return s.Id
}

// LogTypesGetReq 功能获取请求参数
type LogTypesGetReq struct {
     Id int `uri:"id"`
}
func (s *LogTypesGetReq) GetId() interface{} {
	return s.Id
}

// LogTypesDeleteReq 功能删除请求参数
type LogTypesDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *LogTypesDeleteReq) GetId() interface{} {
	return s.Ids
}
