package models

import (
	"go-admin/common/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SysUser struct {
	UserId                       int      `gorm:"primaryKey;autoIncrement;comment:编码"  json:"userId"`
	Username                     string   `json:"username" gorm:"size:64;comment:用户名"`
	Password                     string   `json:"-" gorm:"size:128;comment:密码"`
	NickName                     string   `json:"nickName" gorm:"size:128;comment:中文名称"`
	NickNameEn                   string   `json:"nickNameEn" gorm:"size:64;comment:英文名称"`
	Phone                        string   `json:"phone" gorm:"size:11;comment:手机号"`
	RoleId                       int      `json:"roleId" gorm:"size:20;comment:角色ID"`
	Salt                         string   `json:"-" gorm:"size:255;comment:加盐"`
	Avatar                       string   `json:"avatar" gorm:"size:255;comment:头像"`
	Sex                          string   `json:"sex" gorm:"size:255;comment:性别"`
	Email                        string   `json:"email" gorm:"size:128;comment:邮箱"`
	DeptId                       int      `json:"deptId" gorm:"size:20;comment:部门"`
	PostId                       int      `json:"postId" gorm:"size:20;comment:岗位"`
	Telephone                    string   `json:"telephone" gorm:"size:30;comment:座机号"`
	Fax                          string   `json:"fax" gorm:"size:20;comment:传真"`
	Remark                       string   `json:"remark" gorm:"size:255;comment:备注"`
	Status                       string   `json:"status" gorm:"size:4;comment:状态"`
	SendEmailStatus              string   `json:"sendEmailStatus" gorm:"size:2;comment:生成正式领用单是否邮件通知 1.是 0.否 默认0"`
	AuthorityCompanyId           string   `json:"authorityCompanyId" gorm:"size:500;comment:用户公司权限 公司id逗号分隔"`
	AuthorityWarehouseId         string   `json:"authorityWarehouseId" gorm:"size:500;comment:用户仓库权限 仓库id逗号分隔"`
	AuthorityWarehouseAllocateId string   `json:"authorityWarehouseAllocateId" gorm:"size:500;comment:用户仓库调拨权限 仓库id逗号分隔"`
	AuthorityVendorId            string   `json:"authorityVendorId" gorm:"size:500;comment:用户货主权限 公司id逗号分隔"`
	DeptIds                      []int    `json:"deptIds" gorm:"-"`
	PostIds                      []int    `json:"postIds" gorm:"-"`
	RoleIds                      []int    `json:"roleIds" gorm:"-"`
	Dept                         *SysDept `json:"dept"`
	models.ControlBy
	models.ModelTime
}

func (SysUser) TableName() string {
	return "sys_user"
}

func (e *SysUser) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysUser) GetId() interface{} {
	return e.UserId
}

// 加密
func (e *SysUser) Encrypt() (err error) {
	if e.Password == "" {
		return
	}

	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost); err != nil {
		return
	} else {
		e.Password = string(hash)
		return
	}
}

func (e *SysUser) BeforeCreate(_ *gorm.DB) error {
	return e.Encrypt()
}

func (e *SysUser) BeforeUpdate(_ *gorm.DB) error {
	var err error
	if e.Password != "" {
		err = e.Encrypt()
	}
	return err
}

func (e *SysUser) AfterFind(_ *gorm.DB) error {
	e.DeptIds = []int{e.DeptId}
	e.PostIds = []int{e.PostId}
	e.RoleIds = []int{e.RoleId}
	return nil
}

// 用户名是否存在
func (e *SysUser) UsernameExist(tx *gorm.DB, username string) (bool, error) {
	var i int64
	err := tx.Model(&e).Where("username = ?", username).Count(&i).Error
	if err != nil {
		return true, err
	}

	if i > 0 {
		return true, nil
	}
	return false, nil
}

// 手机号是否存在
func (e *SysUser) PhoneExist(tx *gorm.DB, phone string) (bool, error) {
	var i int64
	err := tx.Model(&e).Where("phone = ?", phone).Count(&i).Error
	if err != nil {
		return true, err
	}

	if i > 0 {
		return true, nil
	}
	return false, nil
}
