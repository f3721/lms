package models

import "time"

// MallUsers  商城用户宽表
type MallUsers struct {
	Id            int64     `json:"id" gorm:"column:id" db:"id" copier:"-"`                          //  primary key
	TenantId      string    `json:"tenantId" gorm:"column:tenant_id" db:"tenant_id"`                 //  租户id
	UserEmail     string    `json:"userEmail" gorm:"column:user_email" db:"user_email"`              //  用户邮箱
	LoginName     string    `json:"loginName" gorm:"column:login_name" db:"login_name"`              //  登录名称
	LoginPassword string    `json:"loginPassword" gorm:"column:login_password" db:"login_password" ` //  登录密码
	UserPhone     string    `json:"userPhone" gorm:"column:user_phone" db:"user_phone"`              //  用户手机
	UserName      string    `json:"userName" gorm:"column:user_name" db:"user_name"`                 //  用户名称
	UserStatus    int64     `json:"userStatus" gorm:"column:user_status" db:"user_status"`           //  用户状态（1可用，0不可用）
	CreateBy      int64     `json:"createBy" gorm:"column:create_by" db:"create_by"`                 //  创建者
	UpdateBy      int64     `json:"updateBy" gorm:"column:update_by" db:"update_by"`                 //  更新者
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at" db:"created_at"`              //  创建时间
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at" db:"updated_at"`              //  最后更新时间
	DeletedAt     time.Time `json:"deletedAt" gorm:"-" db:"deleted_at"`                              //  删除时间
	CreateByName  string    `json:"createByName" gorm:"column:create_by_name" db:"create_by_name"`   //  创建人姓名
	UpdateByName  string    `json:"updateByName" gorm:"column:update_by_name" db:"update_by_name"`   //  更新人姓名
	TenantName    string    `json:"tenantName" gorm:"-" db:"-"`                                      //  自定义:租户名称
}

func (MallUsers) TableName() string {
	return "mall_users"
}
