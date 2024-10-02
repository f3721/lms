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

type CategoryAttribute struct {
	service.Service
}

// GetList 获取CategoryAttribute列表
func (e *CategoryAttribute) GetList(c *dto.CategoryAttributeGetPageReq, p *actions.DataPermission, list *[]models.CategoryAttribute) error {
	var err error
	var data models.CategoryAttribute

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			actions.Permission(data.TableName(), p),
		).
		Joins("Attribute").
		Find(list).Error
	if err != nil {
		e.Log.Errorf("CategoryAttributeService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取CategoryAttribute对象
func (e *CategoryAttribute) Get(d *dto.CategoryAttributeGetReq, p *actions.DataPermission, model *models.CategoryAttribute) error {
	var data models.CategoryAttribute

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Joins("Attribute").
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCategoryAttribute error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建CategoryAttribute对象
func (e *CategoryAttribute) Insert(c *dto.CategoryAttributeInsertReq) error {
	var err error

	resFlag, _ := e.GetCategoryAttribute(&dto.CategoryAttributeWhere{
		Id:          0,
		CategoryId:  c.CategoryId,
		AttributeId: c.AttributeId,
	})
	if resFlag {
		return errors.New("该属性已存在！")
	}

	var data models.CategoryAttribute
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("CategoryAttributeService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改CategoryAttribute对象
func (e *CategoryAttribute) Update(c *dto.CategoryAttributeUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.CategoryAttribute{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	resFlag, _ := e.GetCategoryAttribute(&dto.CategoryAttributeWhere{
		Id:          data.Id,
		CategoryId:  c.CategoryId,
		AttributeId: c.AttributeId,
	})
	if resFlag {
		return errors.New("该属性已存在！")
	}

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("CategoryAttributeService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除CategoryAttribute
func (e *CategoryAttribute) Remove(d *dto.CategoryAttributeDeleteReq, p *actions.DataPermission) error {
	var data models.CategoryAttribute

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Unscoped().Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveCategoryAttribute error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// GetAttrsByCategory 通过产品线获取产品附加属性
func (e *CategoryAttribute) GetAttrsByCategory(categoryIds []int, data *[]dto.AttrsKeyName) error {
	var model models.CategoryAttribute
	err := e.Orm.Model(&model).
		Joins("inner join attribute_def as a ON(category_attribute.attribute_id = a.id)").
		Where("category_attribute.category_id in ?", categoryIds).
		Select("a.name_zh as Name,category_attribute.attribute_id as KeyName").
		Group("category_attribute.attribute_id").
		Order("category_attribute.id ASC").
		Find(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *CategoryAttribute) GetCategoryAttribute(d *dto.CategoryAttributeWhere) (bool, error) {
	var result struct {
		Id int
	}
	var model models.CategoryAttribute
	db := e.Orm.Model(&model).
		Where("category_id = ?", d.CategoryId).
		Where("attribute_id = ?", d.AttributeId).
		Select("id")
	if d.Id != 0 {
		db.Where("id <> ?", d.Id)
	}
	if err := db.Scan(&result).Error; err != nil {
		return false, err
	}
	if result.Id != 0 {
		return true, nil
	}
	return false, nil
}
