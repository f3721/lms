package dto

import (
	"go-admin/app/uc/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
	"time"
)

type UserLoginLogGetPageReq struct {
	dto.Pagination `search:"-"`
	UserId         int    `form:"userId"  search:"type:exact;column:user_id;table:user_login_log" comment:"用户ID"` //用户ID
	Status         string `form:"status"  search:"type:exact;column:status;table:user_login_log" comment:"状态"`    //状态
	UserLoginLogOrder
}

type UserLoginLogOrder struct {
	Id            string `form:"idOrder"  search:"type:order;column:id;table:user_login_log"`
	UserId        string `form:"userIdOrder"  search:"type:order;column:user_id;table:user_login_log"`
	Username      string `form:"usernameOrder"  search:"type:order;column:username;table:user_login_log"`
	Status        string `form:"statusOrder"  search:"type:order;column:status;table:user_login_log"`
	Ipaddr        string `form:"ipaddrOrder"  search:"type:order;column:ipaddr;table:user_login_log"`
	LoginLocation string `form:"loginLocationOrder"  search:"type:order;column:login_location;table:user_login_log"`
	Browser       string `form:"browserOrder"  search:"type:order;column:browser;table:user_login_log"`
	Os            string `form:"osOrder"  search:"type:order;column:os;table:user_login_log"`
	Platform      string `form:"platformOrder"  search:"type:order;column:platform;table:user_login_log"`
	LoginTime     string `form:"loginTimeOrder"  search:"type:order;column:login_time;table:user_login_log"`
	Remark        string `form:"remarkOrder"  search:"type:order;column:remark;table:user_login_log"`
	Msg           string `form:"msgOrder"  search:"type:order;column:msg;table:user_login_log"`
	CreatedAt     string `form:"createdAtOrder"  search:"type:order;column:created_at;table:user_login_log"`
	UpdatedAt     string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:user_login_log"`
	CreateBy      string `form:"createByOrder"  search:"type:order;column:create_by;table:user_login_log"`
	UpdateBy      string `form:"updateByOrder"  search:"type:order;column:update_by;table:user_login_log"`
	CreateByName  string `form:"createByNameOrder"  search:"type:order;column:create_by_name;table:user_login_log"`
	UpdateByName  string `form:"updateByNameOrder"  search:"type:order;column:update_by_name;table:user_login_log"`
}

func (m *UserLoginLogGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type UserLoginLogInsertReq struct {
	Id            int       `json:"-" comment:"主键编码"`            // 主键编码
	UserId        int       `json:"userId" comment:"用户ID"`       // 用户ID
	Username      string    `json:"username" comment:"用户名"`      // 用户名
	Status        string    `json:"status" comment:"状态"`         // 状态
	Ipaddr        string    `json:"ipaddr" comment:"ip地址"`       // ip地址
	LoginLocation string    `json:"loginLocation" comment:"归属地"` // 归属地
	Browser       string    `json:"browser" comment:"浏览器"`       // 浏览器
	Os            string    `json:"os" comment:"系统"`             // 系统
	Platform      string    `json:"platform" comment:"固件"`       // 固件
	LoginTime     time.Time `json:"loginTime" comment:"登录时间"`    // 登录时间
	Remark        string    `json:"remark" comment:"备注"`         // 备注
	Msg           string    `json:"msg" comment:"信息"`            // 信息
	common.ControlBy
}

func (s *UserLoginLogInsertReq) Generate(model *models.UserLoginLog) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UserId = s.UserId
	model.Username = s.Username
	model.Status = s.Status
	model.Ipaddr = s.Ipaddr
	model.LoginLocation = s.LoginLocation
	model.Browser = s.Browser
	model.Os = s.Os
	model.Platform = s.Platform
	model.LoginTime = s.LoginTime
	model.Remark = s.Remark
	model.Msg = s.Msg
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.CreateByName = s.CreateByName
}

func (s *UserLoginLogInsertReq) GetId() interface{} {
	return s.Id
}

type UserLoginLogUpdateReq struct {
	Id            int       `uri:"id" comment:"主键编码"`            // 主键编码
	UserId        int       `json:"userId" comment:"用户ID"`       // 用户ID
	Username      string    `json:"username" comment:"用户名"`      // 用户名
	Status        string    `json:"status" comment:"状态"`         // 状态
	Ipaddr        string    `json:"ipaddr" comment:"ip地址"`       // ip地址
	LoginLocation string    `json:"loginLocation" comment:"归属地"` // 归属地
	Browser       string    `json:"browser" comment:"浏览器"`       // 浏览器
	Os            string    `json:"os" comment:"系统"`             // 系统
	Platform      string    `json:"platform" comment:"固件"`       // 固件
	LoginTime     time.Time `json:"loginTime" comment:"登录时间"`    // 登录时间
	Remark        string    `json:"remark" comment:"备注"`         // 备注
	Msg           string    `json:"msg" comment:"信息"`            // 信息
	common.ControlBy
}

func (s *UserLoginLogUpdateReq) Generate(model *models.UserLoginLog) {
	if s.Id == 0 {
		model.Model = common.Model{Id: s.Id}
	}
	model.UserId = s.UserId
	model.Username = s.Username
	model.Status = s.Status
	model.Ipaddr = s.Ipaddr
	model.LoginLocation = s.LoginLocation
	model.Browser = s.Browser
	model.Os = s.Os
	model.Platform = s.Platform
	model.LoginTime = s.LoginTime
	model.Remark = s.Remark
	model.Msg = s.Msg
	model.UpdateBy = s.UpdateBy
	model.UpdateByName = s.UpdateByName // 添加这而，需要记录是被谁更新的
}

func (s *UserLoginLogUpdateReq) GetId() interface{} {
	return s.Id
}

// UserLoginLogGetReq 功能获取请求参数
type UserLoginLogGetReq struct {
	Id int `uri:"id"`
}

func (s *UserLoginLogGetReq) GetId() interface{} {
	return s.Id
}

// UserLoginLogDeleteReq 功能删除请求参数
type UserLoginLogDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *UserLoginLogDeleteReq) GetId() interface{} {
	return s.Ids
}
