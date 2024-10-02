package admin

import (
	"errors"

    "github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/pc/models"
	"go-admin/app/pc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type LogTypes struct {
	service.Service
}

// GetPage 获取LogTypes列表
func (e *LogTypes) GetPage(c *dto.LogTypesGetPageReq, p *actions.DataPermission, list *[]models.LogTypes, count *int64) error {
	var err error
	var data models.LogTypes

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("LogTypesService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取LogTypes对象
func (e *LogTypes) Get(d *dto.LogTypesGetReq, p *actions.DataPermission, model *models.LogTypes) error {
	var data models.LogTypes

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetLogTypes error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建LogTypes对象
func (e *LogTypes) Insert(c *dto.LogTypesInsertReq) error {
    var err error
    var data models.LogTypes
    c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("LogTypesService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改LogTypes对象
func (e *LogTypes) Update(c *dto.LogTypesUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.LogTypes{}
    e.Orm.Scopes(
            actions.Permission(data.TableName(), p),
        ).First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("LogTypesService Save error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权更新该数据")
    }
    return nil
}

// Remove 删除LogTypes
func (e *LogTypes) Remove(d *dto.LogTypesDeleteReq, p *actions.DataPermission) error {
	var data models.LogTypes

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
        e.Log.Errorf("Service RemoveLogTypes error:%s \r\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("无权删除该数据")
    }
	return nil
}
