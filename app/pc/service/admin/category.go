package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/samber/lo"
	"go-admin/common/global"
	cModel "go-admin/common/models"
	"go-admin/common/utils"
	"gorm.io/gorm"
	"strconv"
	"strings"

	"go-admin/app/pc/models"
	"go-admin/app/pc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Category struct {
	service.Service
}

var getPageReq *dto.CategoryGetPageReq

// GetPage 获取Category列表
func (e *Category) GetPage(c *dto.CategoryGetPageReq, p *actions.DataPermission, list *[]dto.CategoryGetPageResp, count *int) *Category {
	getPageReq = c
	var category = make([]dto.CategoryGetPageResp, 0)
	err := e.getPage(&category).Error
	if err != nil {
		_ = e.AddError(err)
		return e
	}
	var result = make([]dto.CategoryGetPageResp, 0)
	for i := 0; i < len(category); i++ {
		if category[i].ParentId != 0 {
			continue
		}
		categoryInfo := categoryCall(&category, category[i])
		result = append(result, categoryInfo)
	}
	var ids = make([]dto.Ids, 0)
	err = e.getList(c, &ids)
	if err != nil {
		_ = e.AddError(err)
		return e
	}
	var idArr []int
	_ = utils.StructColumn(&idArr, ids, "Id", "")
	categoryFilterList := categoryListFor(&result, idArr)
	lengths := len(categoryFilterList)
	*count = lengths
	// 分页
	start, end := utils.SlicePage(int64(c.GetPageIndex()), int64(c.GetPageSize()), int64(lengths))
	if len(categoryFilterList) > 0 {
		for _, cat := range categoryFilterList[start:end] {
			*list = append(*list, cat)
		}
	}

	return e
}

// getPage 分类分页列表
func (e *Category) getPage(list *[]dto.CategoryGetPageResp) *Category {
	var err error
	var data models.Category

	err = e.Orm.Model(&data).
		Scopes(
			cDto.OrderDest("seq", false),
		).
		Joins("left join category b on category.id = b.parent_id AND `b`.`deleted_at` IS NULL").
		Select("category.*,IF(COUNT(b.id), 1, 0) Haschild").
		//Where("category.`deleted_at` IS NULL").
		Group("category.id").
		Order("category.seq ASC,category.id ASC").
		Find(list).Error
	if err != nil {
		e.Log.Errorf("getCategoryPage error:%s", err)
		_ = e.AddError(err)
		return e
	}
	return e
}

func (e *Category) getList(c *dto.CategoryGetPageReq, ids *[]dto.Ids) error {
	var err error
	var data models.Category
	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			dto.MakeCategoryCondition(c),
		).
		Joins("left join category b on category.id = b.parent_id AND `b`.`deleted_at` IS NULL").
		Select("category.id").
		Group("category.id").
		Find(ids).Error
	if err != nil {
		return err
	}
	return err
}

// Get 获取Category对象
func (e *Category) Get(d *dto.CategoryGetReq, p *actions.DataPermission, model *dto.CategoryGetResp) error {
	var data models.Category
	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Preload("MediaRelationship.MediaInstant").
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
	if model.ParentId != 0 {
		var categoryLevle dto.CategoryLevel
		err = e.getCategoryById(model.ParentId, &categoryLevle)
		if err != nil {
			return err
		}
		model.Level1CatId = categoryLevle.Level1CatId
		model.Level2CatId = categoryLevle.Level2CatId
		model.Level3CatId = categoryLevle.Level3CatId
		model.Level4CatId = categoryLevle.Level4CatId
	}
	model.Addchild = true

	return nil
}

// GetList 获取Category对象
func (e *Category) GetList(d *dto.CategoryGetReq, p *actions.DataPermission, model *[]dto.CategoryGetPageResp) error {
	var data models.Category
	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
			cDto.MakeCondition(d.GetNeedSearch()),
			dto.MakeCategoryGetReqCondition(d),
		).
		Joins("LEFT JOIN category b on category.id = b.parent_id AND `b`.`deleted_at` IS NULL").
		Select("category.*,IF(COUNT(b.id), 1, 0) Haschild").
		//Where("category.`deleted_at` IS NULL").
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
	list := *model
	for i, resp := range list {
		score, _ := strconv.ParseFloat(resp.Tax, 64)
		if score != 0 {
			list[i].TaxTxt = strconv.FormatFloat(score*100, 'f', 0, 64) + "%"
		}
		list[i].StatusTxt = dto.StatusMap[resp.Status]
		//listSearchFieldHighlight(d, &list[i])
	}
	*model = list

	return nil
}

// Insert 创建Category对象
func (e *Category) Insert(c *dto.CategoryInsertReq) error {
	var err error
	var data models.Category

	c.InsertGenerate()
	c.Generate(&data)

	//数据插入校验
	if resFlag, _ := e.CheckNameZh(&data, 0); !resFlag {
		return errors.New("分类已存在")
	}

	tx := e.Orm.Debug().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.Create(&data).Error
	if err != nil {
		tx.Rollback()
		e.Log.Errorf("CategoryService Insert error:%s \r\n", err)
		return err
	}

	// 上传了图片
	if c.MediaInstance.MediaName != "" && c.MediaInstance.MediaDir != "" {
		mediaInstance := MediaInstance{service.Service{
			Orm: tx,
		}}
		var media models.MediaInstance
		// 保存图片实例
		err := mediaInstance.Insert(&dto.MediaInstanceInsertReq{
			MediaName: c.MediaInstance.MediaName,
			MediaDir:  c.MediaInstance.MediaDir,
		}, &media)
		if err != nil {
			tx.Rollback()
			e.Log.Errorf("CategoryService Insert error:%s \r\n", err)
			return err
		}
		// 保存关联关系
		err = e.createMediaRelationship(tx, &data, media.Id, c.MediaInstance)
		if err != nil {
			tx.Rollback()
			e.Log.Errorf("CategoryService Insert error:%s \r\n", err)
			return err
		}
	}
	// 保存分类关系
	err = e.createCategoryPath(&data)
	if err != nil {
		tx.Rollback()
		e.Log.Errorf("CategoryService Insert error:%s \r\n", err)
		return err
	}

	tx.Commit()
	// 生成日志
	dataLog, _ := json.Marshal(&c)
	categoryLog := models.CategoryLog{
		DataId:     data.Id,
		Type:       global.LogTypeCreate,
		Data:       string(dataLog),
		BeforeData: "",
		AfterData:  string(dataLog),
		ControlBy:  c.ControlBy,
	}
	_ = categoryLog.CreateLog("category", e.Orm)
	return nil
}

// Update 修改Category对象
func (e *Category) Update(c *dto.CategoryUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Category{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).
		Preload("MediaRelationship.MediaInstant").
		First(&data, c.GetId())
	oldData := data

	if data.Status == 1 && c.Status == 0 {
		productCategoryService := ProductCategory{e.Service}

		var category models.Category
		_ = e.GetCategoryByParentId(data.Id, &category)
		if category.Id != 0 {
			return errors.New("此分类下有子分类，不能禁用！")
		}

		if productCategoryService.hasProduct(data.Id) {
			return errors.New("此分类下有产品，不能禁用！")
		}
	}

	c.UpdateGenerate()
	c.Generate(&data)

	//数据插入校验
	if resFlag, _ := e.CheckNameZh(&data, data.Id); !resFlag {
		return errors.New("分类已存在")
	}
	tx := e.Orm.Debug().Begin()
	db := tx.Save(&data)
	if err = db.Error; err != nil {
		tx.Rollback()
		e.Log.Errorf("CategoryService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("无权更新该数据")
	}

	// 判断图片是否有更新
	newMedia := c.MediaInstance
	oldMedia := oldData.MediaRelationship.MediaInstant
	var media models.MediaInstance
	mediaInstance := MediaInstance{service.Service{
		Orm: tx,
	}}

	if newMedia.MediaDir != "" && oldMedia.MediaDir != "" && newMedia.MediaDir != oldMedia.MediaDir { // 更新
		mediaInstanceUpdate := models.MediaInstance{
			Model:     cModel.Model{Id: oldData.MediaRelationship.MediaInstantId},
			MediaDir:  newMedia.MediaDir,
			MediaName: newMedia.MediaName,
		}
		err = mediaInstance.Update(&mediaInstanceUpdate)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else if newMedia.MediaDir == "" && oldMedia.MediaDir != "" { // 删除操作
		err = e.deleteMediaRelationship(tx, &data)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else if newMedia.MediaDir != "" && oldMedia.MediaDir == "" { // 新增
		// 保存图片实例
		err = mediaInstance.Insert(&dto.MediaInstanceInsertReq{
			MediaName: newMedia.MediaName,
			MediaDir:  newMedia.MediaDir,
		}, &media)

		if err != nil {
			tx.Rollback()
			return err
		}
		// 保存关联关系
		err = e.createMediaRelationship(tx, &data, media.Id, newMedia)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = e.updateCategoryPath(&oldData, &data)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	// 生成日志
	dataLog, _ := json.Marshal(&c)
	beforeDataStr, _ := json.Marshal(&oldData)
	afterDataStr, _ := json.Marshal(&data)
	categoryLog := models.CategoryLog{
		DataId:     data.Id,
		Type:       global.LogTypeUpdate,
		Data:       string(dataLog),
		BeforeData: string(beforeDataStr),
		AfterData:  string(afterDataStr),
		ControlBy:  c.ControlBy,
	}
	_ = categoryLog.CreateLog("category", e.Orm)

	return nil
}

// Sort 批量排序
func (e *Category) Sort(c *dto.SortReq, p *actions.DataPermission) error {
	var err error
	var data models.Category
	tx := e.Orm.Begin()
	for _, sort := range c.Sort {
		update := map[string]interface{}{
			"seq":            sort.Seq,
			"update_by":      c.UpdateBy,
			"update_by_name": c.UpdateByName,
		}
		tx.Model(&data).Where("id = ?", sort.CategoryId).Updates(update)
	}
	err = tx.Commit().Error
	return err
}

// Remove 删除Category
func (e *Category) Remove(d *dto.CategoryDeleteReq, p *actions.DataPermission) error {
	var err error
	ids := d.GetId().([]int)

	if len(ids) <= 0 {
		return errors.New("请先选择数据！")
	}
	// 分类下有子分类，请先移除子分类！
	for _, val := range ids {
		var category models.Category
		_ = e.GetCategoryByParentId(val, &category)
		if category.Id != 0 {
			err = errors.New(fmt.Sprintf("分类[%s]下有子分类，请先移除子分类！", category.NameZh))
			break
		}
	}
	if err != nil {
		return err
	}
	//分类下有产品，请先移除产品！
	for _, v := range ids {
		productCategoryService := ProductCategory{e.Service}
		if productCategoryService.hasProduct(v) {
			err = errors.New(fmt.Sprintf("分类ID[%d]下有产品，请先移除产品！", v))
			break
		}
	}
	if err != nil {
		return err
	}
	var data models.Category
	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveCategory error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (e *Category) CheckNameZh(model *models.Category, id int) (bool, error) {
	var result struct {
		Id int
	}
	db := e.Orm.Model(&models.Category{}).Select("id")
	if id != 0 {
		db.Where("id <> ?", id)
	}
	if err := db.Where("name_zh = ?", model.NameZh).Where("parent_id = ?", model.ParentId).Scan(&result).Error; err != nil {
		return false, err
	}
	if result.Id != 0 {
		return false, nil
	}
	return true, nil
}

func (e *Category) createMediaRelationship(tx *gorm.DB, category *models.Category, mediaInstanceId int, mediaInstance dto.MediaInstanceInsertReq) error {
	var model models.MediaRelationship
	mediaRelationship := MediaRelationship{service.Service{
		Orm: tx,
	}}
	mediaRelationshipInsertReq := dto.MediaRelationshipInsertReq{
		MediaTypeId:    0, // 0分类 1商品档案
		BuszId:         strconv.Itoa(category.Id),
		MediaInstantId: mediaInstanceId,
		Watermark:      mediaInstance.WaterMark,
		Seq:            mediaInstance.Seq,
		ControlBy:      category.ControlBy,
	}
	err := mediaRelationship.Insert(&mediaRelationshipInsertReq, &model)
	return err
}

func (e *Category) deleteMediaRelationship(tx *gorm.DB, category *models.Category) error {
	mediaRelationship := MediaRelationship{service.Service{
		Orm: tx,
	}}
	err := mediaRelationship.Remove(0, category.Id)
	if err != nil {
		e.Log.Errorf("Service RemoveMediaRelationship error:%s \r\n", err)
		return err
	}
	return nil
}

func (e *Category) createCategoryPath(category *models.Category) error {
	var categoryPaths []models.CategoryPath
	err := e.Orm.Model(&models.CategoryPath{}).Where("category_id = ?", category.ParentId).Order("level ASC").Find(&categoryPaths).Error
	if err != nil {
		return err
	}
	var data []models.CategoryPath
	level := 0
	for _, categoryPath := range categoryPaths {
		data = append(data, models.CategoryPath{
			CategoryId: category.Id,
			PathId:     categoryPath.PathId,
			Level:      level,
		})
		level++
	}
	data = append(data, models.CategoryPath{
		CategoryId: category.Id,
		PathId:     category.Id,
		Level:      level,
	})
	err = e.Orm.Model(&models.CategoryPath{}).Save(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *Category) updateCategoryPath(beforeCategory *models.Category, category *models.Category) error {
	if beforeCategory.ParentId != category.ParentId {
		var model models.CategoryPath
		var data []models.CategoryPath
		categoryPaths, err := e.getCategoryPath(map[string]interface{}{
			"path_id": category.Id,
		})
		if err != nil {
			return err
		}
		if len(categoryPaths) > 0 {
			for _, categoryPath := range categoryPaths {
				deleteCategoryPath := models.CategoryPath{}
				e.Orm.Model(&model).Where("category_id = ?", categoryPath.CategoryId).Where("level < ?", categoryPath.Level).Unscoped().Delete(&deleteCategoryPath)
				var pathIds []int
				// Get the nodes new parents
				nodesCategoryPath, _ := e.getCategoryPath(map[string]interface{}{
					"category_id": category.ParentId,
				})
				for _, path := range nodesCategoryPath {
					pathIds = append(pathIds, path.PathId)
				}
				// Get whats left of the nodes current path
				currentPath, err := e.getCategoryPath(map[string]interface{}{
					"category_id": categoryPath.CategoryId,
				})
				if err != nil {
					return err
				}
				for _, path := range currentPath {
					pathIds = append(pathIds, path.PathId)
				}
				// Combine the paths with a new level
				level := 0
				for _, pathId := range pathIds {
					deleteCategoryPath := models.CategoryPath{}
					e.Orm.Model(&model).Where("category_id = ?", categoryPath.CategoryId).Where("level = ?", level).Where("path_id = ?", pathId).Unscoped().Delete(&deleteCategoryPath)
					data = append(data, models.CategoryPath{
						CategoryId: categoryPath.CategoryId,
						PathId:     pathId,
						Level:      level,
					})
					level++
				}
				err = e.Orm.Model(&models.CategoryPath{}).Save(&data).Error
				if err != nil {
					return err
				}
			}
		} else {
			var deleteCategoryPath models.CategoryPath
			e.Orm.Model(&model).Where("category_id = ?", category.Id).Unscoped().Delete(&deleteCategoryPath)
			categoryPaths, err := e.getCategoryPath(map[string]interface{}{
				"category_id": category.ParentId,
			})
			if err != nil {
				return err
			}
			level := 0
			for _, categoryPath := range categoryPaths {
				data = append(data, models.CategoryPath{
					CategoryId: category.Id,
					PathId:     categoryPath.PathId,
					Level:      level,
				})
				level++
			}
			data = append(data, models.CategoryPath{
				CategoryId: category.Id,
				PathId:     category.Id,
				Level:      level,
			})
			err = e.Orm.Model(&models.CategoryPath{}).Save(&data).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Category) getCategoryPath(configure map[string]interface{}) ([]models.CategoryPath, error) {
	var categoryPaths []models.CategoryPath
	err := e.Orm.Model(&models.CategoryPath{}).Where(configure).Order("level ASC").Find(&categoryPaths).Error
	if err != nil {
		return nil, err
	}
	return categoryPaths, nil
}

func (e *Category) GetCategoryByName(categoryName string, parentId int, category *models.Category) error {
	var model models.Category
	tx := e.Orm.Model(&model)
	if categoryName != "" {
		tx.Where("name_zh = ?", categoryName)
	}
	if parentId != -1 {
		tx.Where("parent_id = ?", parentId)
	}
	tx.Where("deleted_at IS NULL")
	err := tx.Order("id ASC").First(category).Error
	return err
}

func (e *Category) GetCategoryById(categoryId int, category *models.Category) error {
	var model models.Category
	err := e.Orm.Model(&model).Select("id,name_zh,parent_id").First(category, categoryId).Error
	return err
}

func (e *Category) GetCategoryPath(lastId int) (result []models.Category) {
	var category4 models.Category
	_ = e.GetCategoryById(lastId, &category4)
	var category3 models.Category
	_ = e.GetCategoryById(category4.ParentId, &category3)
	var category2 models.Category
	_ = e.GetCategoryById(category3.ParentId, &category2)
	var category1 models.Category
	_ = e.GetCategoryById(category2.ParentId, &category1)

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

func (e *Category) GetCategoryBySku(skuCodes []string) (map[string]int, map[int]dto.CategoryList) {
	productCategory := ProductCategory{e.Service}
	var productCategoryData []models.ProductCategory
	err := productCategory.GetCategoryBySkuMainCate(skuCodes, &productCategoryData)
	if err != nil {
		return nil, nil
	}
	skuMap := make(map[string]int)
	var categoryIds []int
	utils.StructColumn(&skuMap, productCategoryData, "CategoryId", "SkuCode")
	utils.StructColumn(&categoryIds, productCategoryData, "CategoryId", "")
	var categoryLists []dto.CategoryList
	err = e.GetSkuCategoryList(categoryIds, &categoryLists)
	if err != nil {
		return nil, nil
	}
	result := make(map[int]dto.CategoryList)
	for _, item := range categoryLists {
		categoryId := e.getCategoryId(&item)
		result[categoryId] = item
	}
	return skuMap, result
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
		AND c1.deleted_at IS NULL
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

func (e *Category) getCategoryId(data *dto.CategoryList) int {
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

func (e *Category) getCategoryById(parentId int, categoryLevel *dto.CategoryLevel) (err error) {
	var category models.Category
	err = e.GetCategory(parentId, &category)
	if err != nil {
		return err
	}
	if category.Id != 0 {
		currentId := category.Id
		currentParentId := category.ParentId
		currentLevel := category.CateLevel

		if currentLevel == 1 {
			categoryLevel.Level1CatId = currentId

		} else if currentLevel == 2 {
			categoryLevel.Level1CatId = currentParentId
			categoryLevel.Level2CatId = currentId
		} else if currentLevel == 3 {
			var currentCategory models.Category
			err = e.GetCategory(currentParentId, &currentCategory)
			if err != nil {
				return err
			}
			categoryLevel.Level1CatId = currentCategory.ParentId
			categoryLevel.Level2CatId = currentParentId
			categoryLevel.Level3CatId = currentId

		} else if currentLevel == 4 {
			var currentCategory models.Category
			err = e.GetCategory(currentParentId, &currentCategory)
			if err != nil {
				return err
			}
			if currentCategory.Id != 0 {
				categoryLevel.Level2CatId = currentCategory.ParentId
				_ = e.GetCategory(currentCategory.ParentId, &currentCategory)
				if currentCategory.Id != 0 {
					categoryLevel.Level1CatId = currentCategory.ParentId
				}
			}
			categoryLevel.Level3CatId = currentParentId
			categoryLevel.Level4CatId = currentId
		}
	}
	return nil
}

func (e *Category) GetCategory(id int, data *models.Category) error {
	var model models.Category
	err := e.Orm.Model(&model).First(&data, id).Error
	if err != nil {
		e.Log.Errorf("CategoryService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// 删除校验
func (e *Category) GetCategoryByParentId(categoryId int, data *models.Category) error {
	var model models.Category
	err := e.Orm.Model(&model).Where("parent_id = ?", categoryId).Take(data).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		return err
	}
	if err != nil {
		return err
	}
	return nil
}

// categoryCall 构建分类树
func categoryCall(categoryList *[]dto.CategoryGetPageResp, category dto.CategoryGetPageResp) dto.CategoryGetPageResp {
	list := *categoryList
	min := make([]dto.CategoryGetPageResp, 0)
	for j := 0; j < len(list); j++ {
		if category.Id != list[j].ParentId {
			continue
		}
		mi := dto.CategoryGetPageResp{}
		mi.Id = list[j].Id
		mi.CateLevel = list[j].CateLevel
		mi.Seq = list[j].Seq
		mi.NameZh = list[j].NameZh
		mi.NameEn = list[j].NameEn
		mi.ParentId = list[j].ParentId
		//mi.Description = list[j].Description
		mi.Status = list[j].Status
		mi.StatusTxt = dto.StatusMap[list[j].Status]
		//mi.KeyWords = list[j].KeyWords
		mi.Tax = list[j].Tax
		score, _ := strconv.ParseFloat(mi.Tax, 64)
		if score != 0 {
			mi.TaxTxt = strconv.FormatFloat(score*100, 'f', 0, 64) + "%"
		} else {
			mi.TaxTxt = ""
		}
		mi.CategoryTaxCode = list[j].CategoryTaxCode
		mi.Addchild = true
		mi.Haschild = list[j].Haschild
		//mi.CreatedAt = list[j].CreatedAt
		//mi.CreateByName = list[j].CreateByName
		//mi.UpdatedAt = list[j].UpdatedAt
		//mi.UpdateByName = list[j].UpdateByName
		mi.Children = []dto.CategoryGetPageResp{}
		if mi.Haschild {
			ms := categoryCall(categoryList, mi)
			min = append(min, ms)
		} else {
			min = append(min, mi)
		}
	}
	category.StatusTxt = dto.StatusMap[category.Status]
	score, _ := strconv.ParseFloat(category.Tax, 64)
	if score != 0 {
		category.TaxTxt = strconv.FormatFloat(score*100, 'f', 0, 64) + "%"
	} else {
		category.TaxTxt = ""
	}
	category.Children = min
	return category
}

func categoryListFor(categoryList *[]dto.CategoryGetPageResp, ids []int) []dto.CategoryGetPageResp {
	list := *categoryList
	category1 := make([]dto.CategoryGetPageResp, 0)
	for _, cat1 := range list {
		if len(cat1.Children) > 0 {
			cat1.Children = processCategories(ids, cat1.Children)
		}
		category1 = appendCategory(category1, ids, cat1)
	}
	return category1
}

func appendCategory(category []dto.CategoryGetPageResp, ids []int, c dto.CategoryGetPageResp) []dto.CategoryGetPageResp {
	if found := lo.IndexOf[int](ids, c.Id); found != -1 || len(c.Children) > 0 {
		if found != -1 {
			fieldHighlight(&c)
		}
		category = append(category, c)
	}
	return category
}

func processCategories(ids []int, children []dto.CategoryGetPageResp) []dto.CategoryGetPageResp {
	category := make([]dto.CategoryGetPageResp, 0)
	for _, child := range children {
		if len(child.Children) > 0 {
			child.Children = processCategories(ids, child.Children)
		}
		category = appendCategory(category, ids, child)
	}
	return category
}

func highlight(text, keyword string) string {
	highlighted := strings.Replace(text, keyword, "<span style=\"color: red\">"+keyword+"</span>", -1)
	return highlighted
}

func fieldHighlight(category *dto.CategoryGetPageResp) {
	if getPageReq.NameZh != "" {
		category.NameZh = highlight(category.NameZh, getPageReq.NameZh)
	}
	if getPageReq.Status > -1 {
		category.StatusTxt = highlight(category.StatusTxt, category.StatusTxt)
	}
	if getPageReq.Tax != "" {
		category.TaxTxt = highlight(category.TaxTxt, category.TaxTxt)
	}
	if getPageReq.CategoryTaxCode != "" {
		category.CategoryTaxCode = highlight(category.CategoryTaxCode, category.CategoryTaxCode)
	}
}

func listSearchFieldHighlight(req *dto.CategoryGetReq, category *dto.CategoryGetPageResp) {
	if req.NameZh != "" {
		category.NameZh = highlight(category.NameZh, req.NameZh)
	}
	if req.Status != "" {
		category.StatusTxt = highlight(category.StatusTxt, category.StatusTxt)
	}
	if req.Tax != "" {
		category.TaxTxt = highlight(category.TaxTxt, category.TaxTxt)
	}
	if req.CategoryTaxCode != "" {
		category.CategoryTaxCode = highlight(category.CategoryTaxCode, req.CategoryTaxCode)
	}
}
