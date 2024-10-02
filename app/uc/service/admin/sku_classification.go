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

type SkuClassification struct {
	service.Service
}

// GetPage 获取SkuClassification列表
func (e *SkuClassification) GetPage(c *dto.SkuClassificationGetPageReq, p *actions.DataPermission, list *[]models.SkuClassification, count *int64) error {
	var err error
	var data models.SkuClassification

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("SkuClassificationService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取SkuClassification对象
func (e *SkuClassification) Get(d *dto.SkuClassificationGetReq, p *actions.DataPermission, model *models.SkuClassification) error {
	var data models.SkuClassification

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetSkuClassification error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建SkuClassification对象
func (e *SkuClassification) Insert(c *dto.SkuClassificationInsertReq) error {
	var err error
	var data models.SkuClassification
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SkuClassificationService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改SkuClassification对象
func (e *SkuClassification) Update(c *dto.SkuClassificationUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.SkuClassification{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("SkuClassificationService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除SkuClassification
func (e *SkuClassification) Remove(d *dto.SkuClassificationDeleteReq, p *actions.DataPermission) error {
	var data models.SkuClassification

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveSkuClassification error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
