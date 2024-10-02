package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	common "go-admin/common/models"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/uc/models"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type UserRole struct {
	service.Service
}

// GetPage 获取UserRole列表
func (e *UserRole) GetPage(c *dto.UserRoleGetPageReq, p *actions.DataPermission, list *[]models.UserRole, count *int64) error {
	var err error
	var data models.UserRole

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("UserRoleService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取UserRole对象
func (e *UserRole) Get(d *dto.UserRoleGetReq, p *actions.DataPermission, model *models.UserRole) error {
	var data models.UserRole

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetUserRole error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建UserRole对象
func (e *UserRole) Insert(c *dto.UserRoleInsertReq) error {
	var err error
	var data models.UserRole
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("UserRoleService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改UserRole对象
func (e *UserRole) Update(c *dto.UserRoleUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.UserRole{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("UserRoleService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// 用户角色修改
func (e *UserRole) UserRoleSave(db *gorm.DB, userRoles []int, userId int) error {
	if db == nil {
		db = e.Orm
	}
	tx := db.Begin()
	tx.Where("user_id = ?", userId).Delete(&models.UserRole{})
	var userRoleAdds []*models.UserRole
	for _, roleId := range userRoles {
		userRoleAdds = append(userRoleAdds, &models.UserRole{
			UserId: userId,
			RoleId: roleId,
			ControlBy: common.ControlBy{
				CreateBy:     user.GetUserId(e.Orm.Statement.Context.(*gin.Context)),
				CreateByName: user.GetUserName(e.Orm.Statement.Context.(*gin.Context)),
			},
		})
	}
	err := tx.Create(&userRoleAdds).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Remove 删除UserRole
func (e *UserRole) Remove(d *dto.UserRoleDeleteReq, p *actions.DataPermission) error {
	var data models.UserRole

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveUserRole error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
