package mall

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/pc/models"
	"go-admin/app/pc/service/mall/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type CategoryAttribute struct {
	service.Service
}

// GetPage 获取CategoryAttribute列表
func (e *CategoryAttribute) GetPage(c *dto.CategoryAttributeGetPageReq, p *actions.DataPermission, list *[]models.CategoryAttribute, count *int64) error {
	var err error
	var data models.CategoryAttribute

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
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
