package mall

import (
	"errors"
	"github.com/samber/lo"
	"go-admin/common/utils"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/pc/models"
	"go-admin/app/pc/service/mall/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Category struct {
	service.Service
}

// GetPage 获取Category列表
func (e *Category) GetPage(c *dto.CategoryGetPageReq, p *actions.DataPermission, list *[]models.Category, count *int64) error {
	var err error
	var data models.Category

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Order("seq ASC,id ASC").
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("CategoryService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Category对象
func (e *Category) Get(d *dto.CategoryGetReq, p *actions.DataPermission, model *models.Category) error {
	var data models.Category

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCategory error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// GetList 获取Category对象
func (e *Category) GetList(d *dto.CategoryGetReq, p *actions.DataPermission, model *[]dto.CategoryGetPageResp) error {
	var data models.Category
	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Joins("LEFT JOIN category b on category.id = b.parent_id AND `b`.`deleted_at` IS NULL").
		Select("category.*,IF(COUNT(b.id), 1, 0) Haschild").
		Where("category.parent_id = ?", d.GetId()).
		Where("category.`status` = 1").
		Group("category.id").
		Order("category.seq ASC,category.id ASC").
		Find(model).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetCategory error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	return nil
}

func (e *Category) GetCategoryBySku(skuCodes []string, result *map[int]dto.CategoryList) error {
	productCategory := ProductCategory{e.Service}
	var productCategoryData []dto.CategoryInfo
	err := productCategory.GetCategoryBySkuMainCate(skuCodes, &productCategoryData)
	if err != nil {
		return err
	}
	var categoryIds []int
	utils.StructColumn(&categoryIds, productCategoryData, "CategoryId", "")
	categoryLists := make([]dto.CategoryList, 0)
	err = e.GetSkuCategoryList(categoryIds, &categoryLists)
	if err != nil {
		return err
	}
	list := *result
	for _, item := range categoryLists {
		categoryId := e.getCategoryId(item)
		list[categoryId] = item
	}
	result = &list
	return nil
}

func (e *Category) GetSkuCategoryList(categoryIds []int, data *[]dto.CategoryList) error {
	err := e.Orm.Raw(`
		SELECT
			c1.id Id1,
			c2.id Id2,
			c3.id Id3,
			c4.id Id4,
			c1.name_zh NameZh1,
			c2.name_zh NameZh2,
			c3.name_zh NameZh3,
			c4.name_zh NameZh4
		from category as c1
		LEFT JOIN category as c2 on c1.id = c2.parent_id and c2.status != 0 AND c2.deleted_at IS NULL
		LEFT JOIN category as c3 on c2.id = c3.parent_id and c3.status != 0 AND c3.deleted_at IS NULL
		LEFT JOIN category as c4 on c3.id = c4.parent_id and c4.status != 0 AND c4.deleted_at IS NULL
		WHERE
		c1.cate_level = 1
		AND c1.status != 0
		AND (c2.cate_level= 2 or c2.cate_level is NULL)
		AND (c3.cate_level = 3 or c3.cate_level is NULL)
		AND (c4.cate_level = 4 or c4.cate_level  is null )
		AND
		(
			c1.id in ? or c2.id in ? or c3.id in ? or c4.id in ?
		)
	`, categoryIds, categoryIds, categoryIds, categoryIds).Scan(&data).Error
	return err
}

func (e *Category) getCategoryId(data dto.CategoryList) int {
	if data.Id4 > 0 {
		return data.Id4
	} else if data.Id3 > 0 {
		return data.Id3
	} else if data.Id2 > 0 {
		return data.Id2
	} else {
		return data.Id1
	}
}

func (e *Category) GetCategoryListByParentId(parentId int, data *[]models.Category) error {
	var model models.Category
	err := e.Orm.Model(&model).Preload("MediaRelationship.MediaInstant").Where("parent_id = ?", parentId).Where("status = 1").Order("seq ASC,id ASC").Find(data).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("查看对象不存在或无权查看")
	}
	if err != nil {
		return err
	}
	return nil
}

func (e *Category) GetCategoryById(categoryId int, data *models.Category) error {
	var model models.Category
	err := e.Orm.Model(&model).Where("status = 1").First(data, categoryId).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("查看对象不存在或无权查看")
	}
	if err != nil {
		return err
	}
	return nil
}

func (e *Category) GetParentCategoryList(categoryId int, allSkuCategory []int) [][]dto.CategoryNav {
	categoryPath := e.GetCategoryPath(categoryId)
	var list [][]dto.CategoryNav
	for _, category := range categoryPath {
		if category.ParentId != 0 {
			categoryChild := make([]models.Category, 0)
			_ = e.GetCategoryListByParentId(category.ParentId, &categoryChild)
			categoryNav := make([]dto.CategoryNav, 0)
			for _, val := range categoryChild {
				selected := 0
				if category.Id == val.Id {
					selected = 1
				}
				if lo.Contains[int](allSkuCategory, val.Id) {
					categoryNav = append(categoryNav, dto.CategoryNav{
						CategoryId: val.Id,
						NameZh:     val.NameZh,
						Selected:   selected,
					})
				}
			}
			list = append(list, categoryNav)
		} else {
			list = append(list, []dto.CategoryNav{
				{
					CategoryId: category.Id,
					NameZh:     category.NameZh,
					Selected:   1,
				},
			})
		}
	}
	return list
}

func (e *Category) GetCategoryPath(lastId int) (result []models.Category) {
	var category4 models.Category
	if lastId > 0 {
		e.GetCategoryById(lastId, &category4)
	}
	var category3 models.Category
	if category4.ParentId > 0 {
		e.GetCategoryById(category4.ParentId, &category3)
	}
	var category2 models.Category
	if category3.ParentId > 0 {
		e.GetCategoryById(category3.ParentId, &category2)
	}
	var category1 models.Category
	if category2.ParentId > 0 {
		e.GetCategoryById(category2.ParentId, &category1)
	}

	if category4.Id != 0 {
		if category1.Id != 0 {
			result = append(result, category1)
		}
		if category2.Id != 0 {
			result = append(result, category2)
		}
		if category3.Id != 0 {
			result = append(result, category3)
		}
		result = append(result, category4)
	}
	return
}
