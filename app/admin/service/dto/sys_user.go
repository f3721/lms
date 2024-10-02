package dto

import (
	"go-admin/app/admin/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysUserGetPageReq struct {
	dto.Pagination `search:"-"`
	UserId         int    `form:"userId" search:"type:exact;column:user_id;table:sys_user" comment:"用户ID"`
	Username       string `form:"username" search:"type:contains;column:username;table:sys_user" comment:"用户名"`
	NickName       string `form:"nickName" search:"type:contains;column:nick_name;table:sys_user" comment:"昵称"`
	NickNameEn     string `form:"nickNameEn" search:"type:contains;column:nick_name_en;table:sys_user" comment:"英文名称"`
	Phone          string `form:"phone" search:"type:contains;column:phone;table:sys_user" comment:"手机号"`
	RoleId         string `form:"roleId" search:"type:exact;column:role_id;table:sys_user" comment:"角色ID"`
	Sex            string `form:"sex" search:"type:exact;column:sex;table:sys_user" comment:"性别"`
	Email          string `form:"email" search:"type:contains;column:email;table:sys_user" comment:"邮箱"`
	PostId         string `form:"postId" search:"type:exact;column:post_id;table:sys_user" comment:"岗位"`
	Status         string `form:"status" search:"type:exact;column:status;table:sys_user" comment:"状态"`
	DeptJoin       `search:"type:left;on:dept_id:dept_id;table:sys_user;join:sys_dept"`
	SysUserOrder
}

type SysUserOrder struct {
	UserIdOrder    string `search:"type:order;column:user_id;table:sys_user" form:"userIdOrder"`
	UsernameOrder  string `search:"type:order;column:username;table:sys_user" form:"usernameOrder"`
	StatusOrder    string `search:"type:order;column:status;table:sys_user" form:"statusOrder"`
	CreatedAtOrder string `search:"type:order;column:created_at;table:sys_user" form:"createdAtOrder"`
}

type DeptJoin struct {
	DeptId string `search:"type:contains;column:dept_path;table:sys_dept" form:"deptId"`
}

func (m *SysUserGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type ResetSysUserPwdReq struct {
	UserId   int    `json:"userId" comment:"用户ID" vd:"$>0"` // 用户ID
	Password string `json:"password" comment:"密码" vd:"len($)>0"`
	common.ControlBy
}

func (s *ResetSysUserPwdReq) GetId() interface{} {
	return s.UserId
}

func (s *ResetSysUserPwdReq) Generate(model *models.SysUser) {
	if s.UserId != 0 {
		model.UserId = s.UserId
	}
	model.Password = s.Password
}

type UpdateSysUserAvatarReq struct {
	UserId int    `json:"userId" comment:"用户ID" vd:"len($)>0"` // 用户ID
	Avatar string `json:"avatar" comment:"头像" vd:"len($)>0"`
	common.ControlBy
}

func (s *UpdateSysUserAvatarReq) GetId() interface{} {
	return s.UserId
}

func (s *UpdateSysUserAvatarReq) Generate(model *models.SysUser) {
	if s.UserId != 0 {
		model.UserId = s.UserId
	}
	model.Avatar = s.Avatar
}

type UpdateSysUserStatusReq struct {
	UserId int    `json:"userId" comment:"用户ID" vd:"$>0"` // 用户ID
	Status string `json:"status" comment:"状态" vd:"len($)>0"`
	common.ControlBy
}

func (s *UpdateSysUserStatusReq) GetId() interface{} {
	return s.UserId
}

func (s *UpdateSysUserStatusReq) Generate(model *models.SysUser) {
	if s.UserId != 0 {
		model.UserId = s.UserId
	}
	model.Status = s.Status
}

type SysUserInsertReq struct {
	UserId                       int    `json:"userId" comment:"用户ID"` // 用户ID
	Username                     string `json:"username" comment:"用户名" vd:"len($)>0"`
	Password                     string `json:"password" comment:"密码"`
	NickName                     string `json:"nickName" comment:"中文名称" vd:"len($)>0"`
	NickNameEn                   string `json:"nickNameEn" comment:"英文名称" vd:"len($)>0"`
	Phone                        string `json:"phone" comment:"手机号" vd:"len($)>0"`
	RoleId                       int    `json:"roleId" comment:"角色ID"`
	Avatar                       string `json:"avatar" comment:"头像"`
	Sex                          string `json:"sex" comment:"性别"`
	Email                        string `json:"email" comment:"邮箱" vd:"len($)>0,email"`
	DeptId                       int    `json:"deptId" comment:"部门"`
	PostId                       int    `json:"postId" comment:"岗位"`
	Telephone                    string `json:"telephone" comment:"座机号"`
	Fax                          string `json:"fax" comment:"传真"`
	Remark                       string `json:"remark" comment:"备注"`
	Status                       string `json:"status" comment:"状态" vd:"len($)>0" default:"1"`
	SendEmailStatus              string `json:"sendEmailStatus" comment:"生成正式领用单是否邮件通知 1.是 0.否" vd:"@:in($,'0','1')" default:"0"`
	AuthorityCompanyId           string `json:"authorityCompanyId" comment:"用户公司权限"`
	AuthorityWarehouseId         string `json:"authorityWarehouseId" comment:"用户仓库权限"`
	AuthorityWarehouseAllocateId string `json:"authorityWarehouseAllocateId" comment:"用户仓库调拨权限"`
	AuthorityVendorId            string `json:"authorityVendorId" comment:"用户货主权限"`
	common.ControlBy
}

func (s *SysUserInsertReq) Generate(model *models.SysUser) {
	if s.UserId != 0 {
		model.UserId = s.UserId
	}
	model.Username = s.Username
	model.Password = s.Password
	model.NickName = s.NickName
	model.NickNameEn = s.NickNameEn
	model.Phone = s.Phone
	model.RoleId = s.RoleId
	model.Avatar = s.Avatar
	model.Sex = s.Sex
	model.Email = s.Email
	//model.DeptId = s.DeptId
	model.DeptId = 1
	//model.PostId = s.PostId
	model.Remark = s.Remark
	model.Status = s.Status
	model.SendEmailStatus = s.SendEmailStatus
	model.CreateBy = s.CreateBy
	model.Telephone = s.Telephone
	model.Fax = s.Fax
	model.AuthorityCompanyId = s.AuthorityCompanyId
	model.AuthorityWarehouseId = s.AuthorityWarehouseId
	model.AuthorityWarehouseAllocateId = s.AuthorityWarehouseAllocateId
	model.AuthorityVendorId = s.AuthorityVendorId
}

func (s *SysUserInsertReq) GetId() interface{} {
	return s.UserId
}

type SysUserUpdateReq struct {
	UserId                       int    `json:"userId" comment:"用户ID"` // 用户ID
	Username                     string `json:"username" comment:"用户名" vd:"len($)>0"`
	NickName                     string `json:"nickName" comment:"昵称" vd:"len($)>0"`
	NickNameEn                   string `json:"nickNameEn" comment:"英文名称"`
	Phone                        string `json:"phone" comment:"手机号" vd:"len($)>0"`
	RoleId                       int    `json:"roleId" comment:"角色ID"`
	Avatar                       string `json:"avatar" comment:"头像"`
	Sex                          string `json:"sex" comment:"性别"`
	Email                        string `json:"email" comment:"邮箱" vd:"len($)>0,email"`
	DeptId                       int    `json:"deptId" comment:"部门"`
	PostId                       int    `json:"postId" comment:"岗位"`
	Telephone                    string `json:"telephone" comment:"座机号"`
	Fax                          string `json:"fax" comment:"传真"`
	Remark                       string `json:"remark" comment:"备注"`
	Status                       string `json:"status" comment:"状态" default:"1"`
	SendEmailStatus              string `json:"sendEmailStatus" comment:"生成正式领用单是否邮件通知 1.是 0.否" default:"0"`
	AuthorityCompanyId           string `json:"authorityCompanyId" comment:"用户公司权限"`
	AuthorityWarehouseId         string `json:"authorityWarehouseId" comment:"用户仓库权限"`
	AuthorityWarehouseAllocateId string `json:"authorityWarehouseAllocateId" comment:"用户仓库调拨权限"`
	AuthorityVendorId            string `json:"authorityVendorId" comment:"用户货主权限"`
	common.ControlBy
}

func (s *SysUserUpdateReq) Generate(model *models.SysUser) {
	if s.UserId != 0 {
		model.UserId = s.UserId
	}
	model.Username = s.Username
	model.NickName = s.NickName
	model.NickNameEn = s.NickNameEn
	model.Phone = s.Phone
	model.RoleId = s.RoleId
	model.Avatar = s.Avatar
	model.Sex = s.Sex
	model.Email = s.Email
	model.DeptId = s.DeptId
	model.PostId = s.PostId
	model.Remark = s.Remark
	model.Status = s.Status
	model.SendEmailStatus = s.SendEmailStatus
	model.Telephone = s.Telephone
	model.Fax = s.Fax
	model.AuthorityCompanyId = s.AuthorityCompanyId
	model.AuthorityWarehouseId = s.AuthorityWarehouseId
	model.AuthorityWarehouseAllocateId = s.AuthorityWarehouseAllocateId
	model.AuthorityVendorId = s.AuthorityVendorId
}

func (s *SysUserUpdateReq) GetId() interface{} {
	return s.UserId
}

type SysUserById struct {
	dto.ObjectById
	common.ControlBy
}

type UploadSignImg struct {
	dto.ObjectById
	SignImage string `form:"signImage" comment:"电子签照片"`
	common.ControlBy
}

func (s *SysUserById) GetId() interface{} {
	if len(s.Ids) > 0 {
		s.Ids = append(s.Ids, s.Id)
		return s.Ids
	}
	return s.Id
}

func (s *SysUserById) GenerateM() (common.ActiveRecord, error) {
	return &models.SysUser{}, nil
}

// PassWord 密码
type PassWord struct {
	NewPassword string `json:"newPassword" vd:"len($)>0"`
	OldPassword string `json:"oldPassword" vd:"len($)>0"`
}

type SysUserUpdatePermissionReq struct {
	UserId                       int    `json:"userId" comment:"用户ID"` // 用户ID
	AuthorityCompanyId           string `json:"authorityCompanyId" comment:"用户公司权限"`
	AuthorityWarehouseId         string `json:"authorityWarehouseId" comment:"用户仓库权限"`
	AuthorityWarehouseAllocateId string `json:"authorityWarehouseAllocateId" comment:"用户仓库调拨权限"`
	AuthorityVendorId            string `json:"authorityVendorId" comment:"用户货主权限"`
	common.ControlBy
}

func (s *SysUserUpdatePermissionReq) Generate(model *models.SysUser) {
	if s.UserId != 0 {
		model.UserId = s.UserId
	}

	if s.AuthorityCompanyId != "" {
		if model.AuthorityCompanyId != "" {
			model.AuthorityCompanyId = model.AuthorityCompanyId + "," + s.AuthorityCompanyId
		} else {
			model.AuthorityCompanyId = s.AuthorityCompanyId
		}
	}

	if s.AuthorityWarehouseId != "" {
		if model.AuthorityWarehouseId != "" {
			model.AuthorityWarehouseId = model.AuthorityWarehouseId + "," + s.AuthorityWarehouseId
		} else {
			model.AuthorityWarehouseId = s.AuthorityWarehouseId
		}
	}

	if s.AuthorityWarehouseAllocateId != "" {
		if model.AuthorityWarehouseAllocateId != "" {
			model.AuthorityWarehouseAllocateId = model.AuthorityWarehouseAllocateId + "," + s.AuthorityWarehouseAllocateId
		} else {
			model.AuthorityWarehouseAllocateId = s.AuthorityWarehouseAllocateId
		}
	}

	if s.AuthorityVendorId != "" {
		if model.AuthorityVendorId != "" {
			model.AuthorityVendorId = model.AuthorityVendorId + "," + s.AuthorityVendorId
		} else {
			model.AuthorityVendorId = s.AuthorityVendorId
		}
	}
}

func (s *SysUserUpdatePermissionReq) GetId() interface{} {
	return s.UserId
}

type SysUserGetResp struct {
	models.SysUser

	CompanyWarehouseSelected  [][]string `json:"companyWarehouseSelected" gorm:"comment:公司仓库选中项"`
	WarehouseAllocateSelected [][]string `json:"warehouseAllocateSelected" gorm:"comment:仓库调拨选中项"`
}

type SysUserGetPageResp struct {
	models.SysUser

	//AuthorityCompanyName   string   `json:"authorityCompanyName" gorm:"size:500;comment:用户公司权限"`
	AuthorityWarehouseName string `json:"authorityWarehouseName" gorm:"size:500;comment:用户仓库权限"`
	//AuthorityWarehouseAllocateName   string   `json:"authorityWarehouseAllocateName" gorm:"size:500;comment:用户仓库调拨权限"`
	AuthorityVendorName string `json:"authorityVendorName" gorm:"size:500;comment:用户货主权限"`
}

// 后台代登录请求
type SysUserProxyLoginReq struct {
	UserId int `json:"userId" comment:"用户ID" vd:"$>0"` // 用户ID
}
