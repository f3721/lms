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

type ProductCategory struct {
	service.Service
}

// GetPage 获取ProductCategory列表
func (e *ProductCategory) GetPage(c *dto.ProductCategoryGetPageReq, p *actions.DataPermission, list *[]models.ProductCategory, count *int64) error {
	var err error
	var data models.ProductCategory

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("ProductCategoryService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取ProductCategory对象
func (e *ProductCategory) Get(d *dto.ProductCategoryGetReq, p *actions.DataPermission, model *models.ProductCategory) error {
	var data models.ProductCategory

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetProductCategory error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建ProductCategory对象
func (e *ProductCategory) Insert(c *dto.ProductCategoryInsertReq) error {
	var err error
	var data models.ProductCategory
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("ProductCategoryService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// BatchInsert 批量插入
func (e *ProductCategory) BatchInsert(data *[]models.ProductCategory) error {
	err := e.Orm.Create(data).Error
	if err != nil {
		e.Log.Errorf("ProductCategoryService BatchInsert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改ProductCategory对象
func (e *ProductCategory) Update(c *dto.ProductCategoryUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.ProductCategory{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("ProductCategoryService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除ProductCategory
func (e *ProductCategory) Remove(skuCode string) error {
	var data models.ProductCategory
	db := e.Orm.Model(&data).Where("sku_code = ?", skuCode).Unscoped().Delete(&data)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

func (e *ProductCategory) GetCategoryBySkuMainCate(skuCodes []string, data *[]models.ProductCategory) error {
	var model models.ProductCategory
	err := e.Orm.Model(&model).Where("sku_code in ?", skuCodes).Where("main_cate_flag = 1").Find(&data).Error
	return err
}

func (e *ProductCategory) GetCategoryBySku(skuCodes []string, data *[]models.ProductCategory) error {
	var model models.ProductCategory
	err := e.Orm.Model(&model).Where("sku_code in ?", skuCodes).Find(&data).Error
	return err
}

// VerifyProductApproval 校验产线
func (e *ProductCategory) VerifyProductApproval(skuCode string, data *models.ProductCategory) error {
	var model models.ProductCategory
	err := e.Orm.Model(&model).Where("sku_code = ?", skuCode).Where("category_id != ''").Group("sku_code").First(&data).Error
	return err
}

func (e *ProductCategory) IsExistProductForCategory(skuCode string, categoryId int) bool {
	var model models.ProductCategory
	var count int64
	err := e.Orm.Model(&model).Where("sku_code = ?", skuCode).Where("category_id = ?", categoryId).Group("sku_code").Count(&count).Error
	if err != nil || count < 1 {
		return false
	}
	return true
}

func (e *ProductCategory) hasProduct(categoryId int) bool {
	var model models.ProductCategory
	var count int64
	err := e.Orm.Model(&model).Where("category_id = ?", categoryId).Count(&count).Error
	if err != nil || count < 1 {
		return false
	}
	return true
}
