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

type MediaInstance struct {
	service.Service
}

// GetPage 获取MediaInstance列表
func (e *MediaInstance) GetPage(c *dto.MediaInstanceGetPageReq, p *actions.DataPermission, list *[]models.MediaInstance, count *int64) error {
	var err error
	var data models.MediaInstance

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("MediaInstanceService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取MediaInstance对象
func (e *MediaInstance) Get(d *dto.MediaInstanceGetReq, p *actions.DataPermission, model *models.MediaInstance) error {
	var data models.MediaInstance

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetMediaInstance error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建MediaInstance对象
func (e *MediaInstance) Insert(c *dto.MediaInstanceInsertReq, model *models.MediaInstance) error {
	c.Generate(model)
	err := e.Orm.Create(model).Error
	if err != nil {
		return err
	} else {
		return nil
	}
}

// Update 修改MediaInstance对象
func (e *MediaInstance) Update(model *models.MediaInstance) error {
	db := e.Orm.Save(model)
	if err := db.Error; err != nil {
		e.Log.Errorf("CategoryService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除MediaInstance
func (e *MediaInstance) Remove(d *dto.MediaInstanceDeleteReq) error {
	var data models.MediaInstance

	db := e.Orm.Model(&data).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveMediaInstance error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
