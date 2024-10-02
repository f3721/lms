package admin

import (
	"errors"

    "github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/uc/models"
	"go-admin/app/uc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type UserPasswordChangeLog struct {
	service.Service
}

// GetPage 获取UserPasswordChangeLog列表
func (e *UserPasswordChangeLog) GetPage(c *dto.UserPasswordChangeLogGetPageReq, p *actions.DataPermission, list *[]models.UserPasswordChangeLog, count *int64) error {
	var err error
	var data models.UserPasswordChangeLog

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("UserPasswordChangeLogService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取UserPasswordChangeLog对象
func (e *UserPasswordChangeLog) Get(d *dto.UserPasswordChangeLogGetReq, p *actions.DataPermission, model *models.UserPasswordChangeLog) error {
	var data models.UserPasswordChangeLog

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetUserPasswordChangeLog error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建UserPasswordChangeLog对象
func (e *UserPasswordChangeLog) Insert(c *dto.UserPasswordChangeLogInsertReq) error {
    var err error
    var data models.UserPasswordChangeLog
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("UserPasswordChangeLogService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改UserPasswordChangeLog对象
func (e *UserPasswordChangeLog) Update(c *dto.UserPasswordChangeLogUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.UserPasswordChangeLog{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("UserPasswordChangeLogService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除UserPasswordChangeLog
func (e *UserPasswordChangeLog) Remove(d *dto.UserPasswordChangeLogDeleteReq, p *actions.DataPermission) error {
	var data models.UserPasswordChangeLog

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveUserPasswordChangeLog error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
