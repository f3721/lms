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

type AttributeDef struct {
	service.Service
}

// GetPage 获取AttributeDef列表
func (e *AttributeDef) GetPage(c *dto.AttributeDefGetPageReq, p *actions.DataPermission, list *[]models.AttributeDef, count *int64) error {
	var err error
	var data models.AttributeDef

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("AttributeDefService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取AttributeDef对象
func (e *AttributeDef) Get(d *dto.AttributeDefGetReq, p *actions.DataPermission, model *models.AttributeDef) error {
	var data models.AttributeDef

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetAttributeDef error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建AttributeDef对象
func (e *AttributeDef) Insert(c *dto.AttributeDefInsertReq) error {
	var err error
	var data models.AttributeDef
	c.Generate(&data)

	//数据插入校验
	if resFlag, _ := data.CheckNameZh(e.Orm, c.NameZh, 0); !resFlag {
		return errors.New("分类属性已存在")
	}
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("AttributeDefService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改AttributeDef对象
func (e *AttributeDef) Update(c *dto.AttributeDefUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.AttributeDef{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	if resFlag, _ := data.CheckNameZh(e.Orm, c.NameZh, c.Id); !resFlag {
		return errors.New("分类属性已存在")
	}

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("AttributeDefService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除AttributeDef
func (e *AttributeDef) Remove(d *dto.AttributeDefDeleteReq, p *actions.DataPermission) error {
	var data models.AttributeDef

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveAttributeDef error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
