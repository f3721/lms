package mall

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/uc/models"
	"go-admin/app/uc/service/mall/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type ManageMenu struct {
	service.Service
}

// GetPage 获取ManageMenu列表
func (e *ManageMenu) GetPage(c *dto.ManageMenuGetPageReq, p *actions.DataPermission, list *[]models.ManageMenu, count *int64) error {
	var err error
	var data models.ManageMenu

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ManageMenuService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取ManageMenu对象
func (e *ManageMenu) Get(d *dto.ManageMenuGetReq, p *actions.DataPermission, model *models.ManageMenu) error {
	var data models.ManageMenu

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetManageMenu error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建ManageMenu对象
func (e *ManageMenu) Insert(c *dto.ManageMenuInsertReq) error {
	var err error
	var data models.ManageMenu
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ManageMenuService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改ManageMenu对象
func (e *ManageMenu) Update(c *dto.ManageMenuUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.ManageMenu{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("ManageMenuService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除ManageMenu
func (e *ManageMenu) Remove(d *dto.ManageMenuDeleteReq, p *actions.DataPermission) error {
	var data models.ManageMenu

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveManageMenu error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
