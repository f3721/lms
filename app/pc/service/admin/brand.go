package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/mozillazg/go-pinyin"
	"go-admin/app/pc/models"
	"go-admin/app/pc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
	"go-admin/common/global"
	"go-admin/common/utils"
	"gorm.io/gorm"
	"strings"
)

type Brand struct {
	service.Service
}

// GetPage 获取Brand列表
func (e *Brand) GetPage(c *dto.BrandGetPageReq, p *actions.DataPermission, list *[]models.Brand, count *int64) error {
	var err error
	var data models.Brand

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
			dto.BrandMakeCondition(c, p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("BrandService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Brand对象
func (e *Brand) Get(d *dto.BrandGetReq, p *actions.DataPermission, model *models.Brand) error {
	var data models.Brand

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetBrand error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建Brand对象
func (e *Brand) Insert(c *dto.BrandInsertReq) error {
	var err error
	var data models.Brand
	c.FirstLetter = strings.ToUpper(e.pinyinFirstLetter(c.BrandZh))
	c.Generate(&data)

	if c.BrandZh != "" {
		if c.BrandEn != "" {
			if resFlag, _ := data.CheckName(e.Orm, c.BrandZh, c.BrandEn, 0); !resFlag {
				return errors.New("此品牌中文+品牌英文已存在！")
			}
		}

		if c.Confirm == 0 {
			if resFlag, _ := data.CheckNameZh(e.Orm, c.BrandZh, 0); !resFlag {
				return errors.New("nameExists")
			}
		}
	}

	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("BrandService Insert error:%s \r\n", err)
		return err
	}
	// 生成日志
	dataLog, _ := json.Marshal(&c)
	brandLog := models.BrandLog{
		DataId:     data.Id,
		Type:       global.LogTypeCreate,
		Data:       string(dataLog),
		BeforeData: "",
		AfterData:  string(dataLog),
		ControlBy:  c.ControlBy,
	}
	_ = brandLog.CreateLog("brand", e.Orm)
	return nil
}

// Update 修改Brand对象
func (e *Brand) Update(c *dto.BrandUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Brand{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.FirstLetter = e.pinyinFirstLetter(c.BrandZh)
	oldData := data
	c.Generate(&data)

	if c.BrandZh != "" {
		if c.BrandEn != "" {
			if resFlag, _ := data.CheckName(e.Orm, c.BrandZh, c.BrandEn, data.Id); !resFlag {
				return errors.New("此品牌中文+品牌英文已存在！")
			}
		}
		if c.Confirm == 0 {
			if resFlag, _ := data.CheckNameZh(e.Orm, c.BrandZh, data.Id); !resFlag {
				return errors.New("nameExists")
			}
		}
	}

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("BrandService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}

	// 生成日志
	dataLog, _ := json.Marshal(&c)
	beforeDataStr, _ := json.Marshal(&oldData)
	afterDataStr, _ := json.Marshal(&data)
	brandLog := models.BrandLog{
		DataId:     data.Id,
		Type:       global.LogTypeUpdate,
		Data:       string(dataLog),
		BeforeData: string(beforeDataStr),
		AfterData:  string(afterDataStr),
		ControlBy:  c.ControlBy,
	}
	_ = brandLog.CreateLog("brand", e.Orm)

	return nil
}

// Remove 删除Brand
func (e *Brand) Remove(d *dto.BrandDeleteReq, p *actions.DataPermission) error {
	var err error
	var data models.Brand
	productService := Product{e.Service}

	for _, val := range d.Ids {
		var product models.Product
		_ = productService.GetProductByBrandId(val, &product)
		if product.Id > 0 {
			err = errors.New(fmt.Sprintf("品牌[ID:%d]有绑定产品【%s】,不可删除品牌!", val, product.SkuCode))
			break
		}
	}
	if err != nil {
		return err
	}
	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveBrand error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (e *Brand) FindBrandInfoByName(nameZh string, nameEn string, data *models.Brand) error {
	var model models.Brand
	db := e.Orm.Model(&model).Where("brand_zh = ?", nameZh).Where("status = ?", 1)
	if nameEn != "" {
		db.Where("brand_en", nameEn)
	}
	err := db.First(&data).Error
	return err
}

func (e *Brand) FindBrandListByName(c *dto.BrandGetPageReq, data *[]models.Brand) error {
	var model models.Brand
	err := e.Orm.Model(&model).Where("brand_zh = ?", c.BrandZh).Find(&data).Error
	return err
}

func (e *Brand) pinyinFirstLetter(name string) string {
	firstLetter := []rune(name)
	if utils.IsAlphanumeric(string(firstLetter[0])) {
		return string(firstLetter[0])
	} else if utils.ContainChinese(string(firstLetter[0])) {
		pinyinApp := pinyin.NewArgs()
		pinyinApp.Style = pinyin.FirstLetter
		letterArr := pinyin.Pinyin(string(firstLetter[0]), pinyinApp)
		return letterArr[0][0]
	}
	return ""
}
