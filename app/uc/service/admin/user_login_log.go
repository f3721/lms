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

type UserLoginLog struct {
	service.Service
}

// GetPage 获取UserLoginLog列表
func (e *UserLoginLog) GetPage(c *dto.UserLoginLogGetPageReq, p *actions.DataPermission, list *[]models.UserLoginLog, count *int64) error {
	var err error
	var data models.UserLoginLog

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("UserLoginLogService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取UserLoginLog对象
func (e *UserLoginLog) Get(d *dto.UserLoginLogGetReq, p *actions.DataPermission, model *models.UserLoginLog) error {
	var data models.UserLoginLog

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetUserLoginLog error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建UserLoginLog对象
func (e *UserLoginLog) Insert(c *dto.UserLoginLogInsertReq) error {
    var err error
    var data models.UserLoginLog
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("UserLoginLogService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改UserLoginLog对象
func (e *UserLoginLog) Update(c *dto.UserLoginLogUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.UserLoginLog{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("UserLoginLogService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除UserLoginLog
func (e *UserLoginLog) Remove(d *dto.UserLoginLogDeleteReq, p *actions.DataPermission) error {
	var data models.UserLoginLog

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveUserLoginLog error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
