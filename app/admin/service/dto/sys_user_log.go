package dto

import (
    "go-admin/common/dto"
)

type SysUserLogGetPageReq struct {
	dto.Pagination     `search:"-"`
    DataId         int    `form:"dataId" search:"type:exact;column:data_id;table:sys_user_log" comment:"用户ID"`
    SysUserLogOrder
}

type SysUserLogOrder struct {
    Id string `form:"idOrder"  search:"type:order;column:id;table:sys_user_log"`
    DataId string `form:"dataIdOrder"  search:"type:order;column:data_id;table:sys_user_log"`
    Type string `form:"typeOrder"  search:"type:order;column:type;table:sys_user_log"`
    Data string `form:"dataOrder"  search:"type:order;column:data;table:sys_user_log"`
    BeforeData string `form:"beforeDataOrder"  search:"type:order;column:before_data;table:sys_user_log"`
    AfterData string `form:"afterDataOrder"  search:"type:order;column:after_data;table:sys_user_log"`
    DiffData string `form:"diffDataOrder"  search:"type:order;column:diff_data;table:sys_user_log"`
    CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:sys_user_log"`
    CreateBy string `form:"createByOrder"  search:"type:order;column:create_by;table:sys_user_log"`
    CreateByName string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:sys_user_log"`
    
}

func (m *SysUserLogGetPageReq) GetNeedSearch() interface{} {
	return *m
}
// SysUserLogGetReq 功能获取请求参数
type SysUserLogGetReq struct {
     Id int `uri:"id"`
}
func (s *SysUserLogGetReq) GetId() interface{} {
	return s.Id
}
