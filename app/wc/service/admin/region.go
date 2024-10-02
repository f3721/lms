package admin

import (
	"errors"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/samber/lo"
	"go-admin/common/utils"
	"gorm.io/gorm"

	"go-admin/app/wc/models"
	"go-admin/app/wc/service/admin/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Region struct {
	service.Service
}

// GetPage 获取Region列表
func (e *Region) GetPage(c *dto.RegionGetPageReq, p *actions.DataPermission, list *[]models.Region, count *int64) error {
	var err error
	var data models.Region
	if c.Name == "" && c.ParentId == "" {
		c.ParentId = "0"
	}
	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSizeNegative(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("RegionService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Region对象
func (e *Region) Get(d *dto.RegionGetReq, p *actions.DataPermission, model *models.Region) error {
	var data models.Region

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetRegion error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建Region对象
func (e *Region) Insert(c *dto.RegionInsertReq) error {
	var err error
	var data models.Region
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("RegionService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改Region对象
func (e *Region) Update(c *dto.RegionUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Region{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("RegionService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除Region
func (e *Region) Remove(d *dto.RegionDeleteReq, p *actions.DataPermission) error {
	var data models.Region

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveRegion error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (e *Region) InnerGetByIds(c *dto.InnerRegionGetByIdsReq, list *[]models.Region) error {
	var err error
	var data models.Region
	err = e.Orm.Model(&data).
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				db.Where("id in ?", lo.Uniq(utils.SplitToInt(c.Ids)))
				return db
			},
		).
		Find(list).Error
	if err != nil {
		e.Log.Errorf("RegionService InnerGetByIds error:%s \r\n", err)
		return err
	}
	return nil
}
