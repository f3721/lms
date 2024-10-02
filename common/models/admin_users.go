package models

import (
	"time"

	"gorm.io/gorm"
)

// AdminUsers  后台用户宽表
type AdminUsers struct {
	Id           int64     `json:"id" gorm:"column:id" db:"id"`                                   //  primary key
	TenantId     string    `json:"tenantId" gorm:"column:tenant_id" db:"tenant_id"`               //  租户id
	Username     string    `json:"username" gorm:"column:username" db:"username"`                 //  用户名
	Password     string    `json:"password" gorm:"column:password" db:"password"`                 //  密码
	Phone        string    `json:"phone" gorm:"column:phone" db:"phone"`                          //  手机号
	Status       string    `json:"status" gorm:"column:status" db:"status"`                       //  状态
	CreateBy     int64     `json:"createBy" gorm:"column:create_by" db:"create_by"`               //  创建者
	UpdateBy     int64     `json:"updateBy" gorm:"column:update_by" db:"update_by"`               //  更新者
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at" db:"created_at"`            //  创建时间
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at" db:"updated_at"`            //  最后更新时间
	DeletedAt    time.Time `json:"deleted_at" gorm:"-" db:"deleted_at"`                           //  删除时间
	CreateByName string    `json:"createByName" gorm:"column:create_by_name" db:"create_by_name"` //  创建人姓名
	UpdateByName string    `json:"updateByName" gorm:"column:update_by_name" db:"update_by_name"` //  更新人姓名
	TenantName   string    `json:"tenantName" gorm:"-" db:"-"`                                    //  自定义:租户名称
}

func (AdminUsers) TableName() string {
	return "admin_users"
}

// 根据租户id+手机查询数据
func (user *AdminUsers) GetByPhone(tx *gorm.DB, tenantId string, phone string) (*AdminUsers, error) {
	err := tx.Table("common.admin_users").Where("tenantId = ?", tenantId).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
