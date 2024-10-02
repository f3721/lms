package admin

import (
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"
	"strconv"
	"strings"

	"go-admin/app/pc/models"
	"go-admin/app/pc/service/admin/dto"
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

// Insert 创建ProductExtAttribute对象
func (e *ProductExtAttribute) Insert(c *dto.ProductExtAttributeInsertReq) error {
	var err error
	var data models.ProductExtAttribute
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ProductExtAttributeService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改ProductExtAttribute对象
func (e *ProductExtAttribute) Update(c *dto.ProductExtAttributeUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.ProductExtAttribute{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("ProductExtAttributeService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除ProductExtAttribute
func (e *ProductExtAttribute) Remove(d *dto.ProductExtAttributeDeleteReq, p *actions.DataPermission) error {
	var data models.ProductExtAttribute

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveProductExtAttribute error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
func (e *ProductExtAttribute) Delete(skuCode string) error {
	var data models.ProductExtAttribute
	err := e.Orm.Model(&data).Where("sku_code", skuCode).Unscoped().Delete(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *ProductExtAttribute) Import(data *dto.ProductAttributeImportTemp) error {
	var err error
	// excel 附加导入产品属性
	if data.AttributeId != "" && data.AttributeValue != "" && data.CategoryId != 0 {
		// 通过产品线获取产品附加属性
		result := make([]dto.AttrsKeyName, 0)
		categoryAttribute := CategoryAttribute{e.Service}
		_ = categoryAttribute.GetAttrsByCategory([]int{data.CategoryId}, &result)
		if len(result) > 0 {
			var errMsg []string
			for _, attribute := range result {
				attributeId, _ := strconv.Atoi(data.AttributeId)
				if attribute.KeyName == attributeId {
					// 判断附加产品属性是否存在
					var productExtAttribute models.ProductExtAttribute
					_ = e.FindRow(data.SkuCode, attributeId, &productExtAttribute)
					// 更新
					if productExtAttribute.SkuCode != "" {
						err = e.importUpdate(&productExtAttribute, data)
						if err != nil {
							errMsg = append(errMsg, fmt.Sprintf("属性值【%s】更新失败！", productExtAttribute.ValueZh))
						}
					} else {
						// 新增
						err = e.importSave(data)
						if err != nil {
							errMsg = append(errMsg, fmt.Sprintf("属性值【%s】保存失败！", productExtAttribute.ValueZh))
						}
					}
				}
			}
			if len(errMsg) > 0 {
				return errors.New(strings.Join(errMsg, "、"))
			}
		}
	}
	return nil
}

func (e *ProductExtAttribute) FindRow(skuCode string, attributeId int, data *models.ProductExtAttribute) error {
	var model models.ProductExtAttribute
	err := e.Orm.Model(&model).Where("sku_code = ?", skuCode).Where("attribute_id = ?", attributeId).First(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *ProductExtAttribute) importUpdate(origin *models.ProductExtAttribute, c *dto.ProductAttributeImportTemp) error {
	origin.ValueZh = c.AttributeValue
	err := e.Orm.Save(&origin).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *ProductExtAttribute) importSave(c *dto.ProductAttributeImportTemp) error {
	var model models.ProductExtAttribute
	c.Generate(&model)
	err := e.Orm.Create(&model).Error
	if err != nil {
		e.Log.Errorf("ProductExtAttributeService importSave error:%s \r\n", err)
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
	if len(c.AttrList) > 0 {
		for k, m := range result {
			for _, attribute := range c.AttrList {
				if m.KeyName == attribute.KeyName {
					result[k].Value = attribute.Value
				}
			}
		}
	} else if c.SkuCode != "" {
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
