package dto

import common "go-admin/common/models"

// 新建审批时间请求
type ApprovalTimeInsertReq struct {
	Hour   string   `json:"hour" comment:"几时" vd:"@:in($,'5','6','7','8','9','10','11','12','13','14','15','16','17','18','19','20','21','22','23');msg:'时间必填[5-23]'"`
	Min    string   `json:"min" comment:"几分" vd:"@:in($,'0','15','30','45');msg:'分钟必填[0,15,30,45]'"`
	Weeks  []string `json:"weeks" comment:"星期几" vd:"@:len($)>=1 && len($)<=7 ;msg:'星期必选[0-6]'"`
	Repeat string   `json:"repeat" comment:"是否重复" vd:"@:in($,'0','1');msg:'范围0和1'"`
}

// 更新审批时间请求
type ApprovalTimeUpdateReq struct {
	Id     int      `uri:"id" comment:"id" vd:"@:in($,0,1,2);msg:'ID范围0-2'"`
	Hour   string   `json:"hour" comment:"几时" vd:"@:in($,'5','6','7','8','9','10','11','12','13','14','15','16','17','18','19','20','21','22','23');msg:'时间必填[5-23]'"`
	Min    string   `json:"min" comment:"几分" vd:"@:in($,'0','15','30','45');msg:'分钟必填[0,15,30,45]'"`
	Weeks  []string `json:"weeks" comment:"星期几" vd:"@:len($)>=1 && len($)<=7 ;msg:'星期必选[0-6]'"`
	Repeat string   `json:"repeat" comment:"是否重复" vd:"@:in($,'0','1');msg:'范围0和1'"`
	common.ControlBy
}

// 审批时间查询请求
type ApprovalTimeGetReq struct {
	Id int `uri:"id" comment:"id" vd:"@:in($,0,1,2);msg:'ID范围0-2'"`
}

type ApprovalTimeDeleteReq struct {
	Id int `json:"id" comment:"id" vd:"@:in($,0,1,2);msg:'ID范围0-2'"`
}
