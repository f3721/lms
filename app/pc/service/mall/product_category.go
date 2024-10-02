package mall

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-admin/common/middleware/mall_handler"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/pc/models"
	"go-admin/app/pc/service/mall/dto"
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

func (e *ProductCategory) GetCategoryBySkuMainCate(skuCodes []string, data *[]dto.CategoryInfo) error {
	var model models.ProductCategory
	ctx := e.Orm.Statement.Context.(*gin.Context)
	warehouseCode := mall_handler.GetUserConfig(ctx).SelectedWarehouseCode
	err := e.Orm.Model(&model).
		Joins("INNER JOIN category ON product_category.category_id = category.id").
		Joins("INNER JOIN goods ON product_category.sku_code = goods.sku_code").
		Where("product_category.sku_code in ?", skuCodes).
		Where("product_category.main_cate_flag = 1").
		Where("goods.warehouse_code = ?", warehouseCode).
		Select([]string{"product_category.sku_code", "category_id", "parent_id"}).
		Find(&data).Error
	return err
}

func (e *ProductCategory) GetCategoryBySkuMainCategoryId(skuCodes string) int {
	var model models.ProductCategory
	err := e.Orm.Model(&model).Where("sku_code = ?", skuCodes).Where("main_cate_flag = 1").First(&model).Error
	if err != nil {
		return 0
	}
	return model.CategoryId
}
