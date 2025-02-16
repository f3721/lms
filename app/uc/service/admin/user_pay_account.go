package admin

import (
	"errors"
	"go-admin/app/uc/service/admin/dto"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/uc/models"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type UserPayAccount struct {
	service.Service
}

// GetPage 获取UserPayAccount列表
func (e *UserPayAccount) GetPage(c *dto.UserPayAccountGetPageReq, p *actions.DataPermission, list *[]models.UserPayAccount, count *int64) error {
	var err error
	var data models.UserPayAccount

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("UserPayAccountService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取UserPayAccount对象
func (e *UserPayAccount) Get(d *dto.UserPayAccountGetReq, p *actions.DataPermission, model *models.UserPayAccount) error {
	var data models.UserPayAccount

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetUserPayAccount error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建UserPayAccount对象
func (e *UserPayAccount) Insert(c *dto.UserPayAccountInsertReq) error {
	var err error
	var data models.UserPayAccount
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("UserPayAccountService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改UserPayAccount对象
func (e *UserPayAccount) Update(c *dto.UserPayAccountUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.UserPayAccount{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("UserPayAccountService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除UserPayAccount
func (e *UserPayAccount) Remove(d *dto.UserPayAccountDeleteReq, p *actions.DataPermission) error {
	var data models.UserPayAccount

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveUserPayAccount error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
