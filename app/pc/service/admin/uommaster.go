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

type Uommaster struct {
	service.Service
}

// GetPage 获取Uommaster列表
func (e *Uommaster) GetPage(c *dto.UommasterGetPageReq, p *actions.DataPermission, list *[]models.Uommaster, count *int64) error {
	var err error
	var data models.Uommaster

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("UommasterService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Uommaster对象
func (e *Uommaster) Get(d *dto.UommasterGetReq, p *actions.DataPermission, model *models.Uommaster) error {
	var data models.Uommaster

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetUommaster error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建Uommaster对象
func (e *Uommaster) Insert(c *dto.UommasterInsertReq) error {
	var err error
	var data models.Uommaster
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("UommasterService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改Uommaster对象
func (e *Uommaster) Update(c *dto.UommasterUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Uommaster{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("UommasterService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除Uommaster
func (e *Uommaster) Remove(d *dto.UommasterDeleteReq, p *actions.DataPermission) error {
	var data models.Uommaster

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveUommaster error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (e *Uommaster) GetByName(name string, uommaster *models.Uommaster) error {
	var model models.Uommaster
	err := e.Orm.Model(&model).Where("uom = ?", name).Take(&uommaster).Error
	return err
}
