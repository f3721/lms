package mall

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/wc/models"
	"go-admin/app/wc/service/mall/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Warehouse struct {
	service.Service
}

// GetPageWithCompanyId 获取某个用户对应公司的Warehouse列表
func (e *Warehouse) GetPageWithCompanyId(companyId int, c *dto.WarehouseGetPageReq, p *actions.DataPermission, list *[]models.Warehouse, count *int64) error {
	var err error
	var data models.Warehouse

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Where("company_id", companyId).
		Where("is_virtual", 0).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("WarehouseService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// GetPage 获取Warehouse列表
func (e *Warehouse) GetPage(c *dto.WarehouseGetPageReq, p *actions.DataPermission, list *[]models.Warehouse, count *int64) error {
	var err error
	var data models.Warehouse

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("WarehouseService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Warehouse对象
func (e *Warehouse) Get(d *dto.WarehouseGetReq, p *actions.DataPermission, model *models.Warehouse) error {
	var data models.Warehouse

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetWarehouse error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建Warehouse对象
func (e *Warehouse) Insert(c *dto.WarehouseInsertReq) error {
	var err error
	var data models.Warehouse
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("WarehouseService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改Warehouse对象
func (e *Warehouse) Update(c *dto.WarehouseUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Warehouse{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("WarehouseService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除Warehouse
func (e *Warehouse) Remove(d *dto.WarehouseDeleteReq, p *actions.DataPermission) error {
	var data models.Warehouse

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveWarehouse error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
