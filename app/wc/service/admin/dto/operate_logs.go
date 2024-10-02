package dto

import (
	"go-admin/common/dto"
)

type OperateLogsGetPageReq struct {
	dto.Pagination `search:"-"`
	DataId         string `form:"dataId"  search:"type:exact;column:data_id;table:operate_logs" comment:"数据id eg:WH0001"`
	ModelName      string `form:"modelName"  search:"type:exact;column:model_name;table:operate_logs" comment:"模型name"`
}

func (m *OperateLogsGetPageReq) GetNeedSearch() interface{} {
	return *m
}

// OperateLogsGetReq 功能获取请求参数
type OperateLogsGetReq struct {
	Id int `uri:"id"`
}

func (s *OperateLogsGetReq) GetId() interface{} {
	return s.Id
}

// OperateLogsDeleteReq 功能删除请求参数
type OperateLogsDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *OperateLogsDeleteReq) GetId() interface{} {
	return s.Ids
}
