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

type Classification struct {
	service.Service
}

// GetPage 获取Classification列表
func (e *Classification) GetPage(c *dto.ClassificationGetPageReq, p *actions.DataPermission, list *[]models.Classification, count *int64) error {
	var err error
	var data models.Classification

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ClassificationService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Classification对象
func (e *Classification) Get(d *dto.ClassificationGetReq, p *actions.DataPermission, model *models.Classification) error {
	var data models.Classification

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetClassification error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建Classification对象
func (e *Classification) Insert(c *dto.ClassificationInsertReq) error {
	var err error
	var data models.Classification
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ClassificationService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改Classification对象
func (e *Classification) Update(c *dto.ClassificationUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Classification{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("ClassificationService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除Classification
func (e *Classification) Remove(d *dto.ClassificationDeleteReq, p *actions.DataPermission) error {
	var data models.Classification

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveClassification error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
