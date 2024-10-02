package mall_handler

import (
	"errors"
	"go-admin/common/global"
	"go-admin/common/utils"

	log "github.com/go-admin-team/go-admin-core/logger"
	"gorm.io/gorm"
)

type Login struct {
	Username string `form:"UserName" json:"username" binding:"required"`
	Password string `form:"Password" json:"password" binding:"required"`
}

func (u *Login) GetUser(tx *gorm.DB) (user UserInfo, role RoleInfo, err error) {
	ucPrefix := global.GetTenantUcDBNameWithDB(tx)
	err = tx.Table(ucPrefix+".user_info").Where("(login_name = ? or user_phone =? or user_name=? or user_email=?)  and user_status = '1'", u.Username, u.Username, u.Username, u.Username).First(&user).Error
	if err != nil {
		log.Errorf("get user error, %s", err.Error())
		return
	}

	pwd := utils.Md5Uc(u.Password)
	if pwd != user.LoginPassword {
		log.Errorf("user login error, password not right, user:%s, password:%s", u.Username, u.Password)
		err = errors.New("密码不正确")
		return
	}

	// 获取用户角色
	var userRole UserRole
	err = tx.Table(ucPrefix+".user_role").Where("user_id = ? ", user.Id).First(&userRole).Error
	if err != nil {
		log.Errorf("get user_role error, %s", err.Error())
		err = errors.New("用户角色未设置")
		return
	}

	// 获取角色详情
	err = tx.Table(ucPrefix+".role_info").Where("id = ? ", userRole.RoleId).First(&role).Error
	if err != nil {
		log.Errorf("get role_info error, %s", err.Error())
		err = errors.New("用户角色未设置")
		return
	}
	return
}
