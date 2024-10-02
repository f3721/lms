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

type ProductExtAttribute struct {
	service.Service
}

// GetPage 获取ProductExtAttribute列表
func (e *ProductExtAttribute) GetPage(c *dto.ProductExtAttributeGetPageReq, p *actions.DataPermission, list *[]models.ProductExtAttribute, count *int64) error {
	var err error
	var data models.ProductExtAttribute

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ProductExtAttributeService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取ProductExtAttribute对象
func (e *ProductExtAttribute) Get(d *dto.ProductExtAttributeGetReq, p *actions.DataPermission, model *models.ProductExtAttribute) error {
	var data models.ProductExtAttribute

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetProductExtAttribute error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

func (e *ProductExtAttribute) GetAttrs(c *dto.GetProductExtAttributeReq) ([]dto.AttrsKeyName, error) {
	result := make([]dto.AttrsKeyName, 0)
	categoryAttribute := CategoryAttribute{e.Service}
	err := categoryAttribute.GetAttrsByCategory([]int{c.CategoryId}, &result)
	if err != nil {
		return nil, err
	}
	if c.SkuCode != "" {
		productAttrs := make([]dto.ProductExtAttribute, 0)
		err := e.GetProductExtAttrs(c.SkuCode, &productAttrs)
		if err != nil {
			return nil, err
		}
		for k, m := range result {
			for _, attr := range productAttrs {
				if m.KeyName == attr.KeyName {
					result[k].Value = attr.Value
				}
			}
		}
	}

	return result, nil
}

// GetProductExtAttrs 获取产品附加属性
func (e *ProductExtAttribute) GetProductExtAttrs(skuCode string, data *[]dto.ProductExtAttribute) error {
	var model models.ProductExtAttribute
	err := e.Orm.Model(&model).
		Select("product_ext_attribute.attribute_id as KeyName,product_ext_attribute.value_zh as Value,ad.name_zh as Name").
		Joins("LEFT JOIN attribute_def as ad ON product_ext_attribute.attribute_id = ad.id").Where("product_ext_attribute.sku_code = ?", skuCode).Find(&data).Error
	if err != nil {
		return err
	}
	return nil
}
