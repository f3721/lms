package admin

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/app/pc/models"
	"go-admin/app/pc/service/admin/dto"
	"go-admin/common/actions"
)

type AttributeConfig struct {
	service.Service
}

// Insert 创建AttributeConfig对象
func (e *AttributeConfig) Insert(c *dto.AttributeConfigInsertReq) error {
	var err error
	var data models.AttributeConfig
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("AttributeConfigService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改AttributeConfig对象
func (e *AttributeConfig) Update(c *dto.AttributeConfigUpdateReq) error {
	var err error
	var data = models.AttributeConfig{}
	e.Orm.First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("AttributeConfigService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除AttributeConfig
func (e *AttributeConfig) Remove(d *dto.AttributeConfigDeleteReq, p *actions.DataPermission) error {
	var data models.AttributeConfig

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveAttributeConfig error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
